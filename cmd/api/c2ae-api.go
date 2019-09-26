package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/olivere/elastic"
	slibcfg "github.com/teserakt-io/serverlib/config"
	sliblog "github.com/teserakt-io/serverlib/log"
	slibpath "github.com/teserakt-io/serverlib/path"

	"github.com/teserakt-io/automation-engine/internal/api"
	"github.com/teserakt-io/automation-engine/internal/config"
	"github.com/teserakt-io/automation-engine/internal/engine"
	"github.com/teserakt-io/automation-engine/internal/engine/actions"
	"github.com/teserakt-io/automation-engine/internal/engine/watchers"
	"github.com/teserakt-io/automation-engine/internal/events"
	"github.com/teserakt-io/automation-engine/internal/models"
	"github.com/teserakt-io/automation-engine/internal/monitoring"
	"github.com/teserakt-io/automation-engine/internal/pb"
	"github.com/teserakt-io/automation-engine/internal/services"
)

// Provided by build script
var gitCommit string
var gitTag string
var buildDate string

func main() {
	exitCode := 0
	defer os.Exit(exitCode)

	printVersion()

	globalCtx, globalCancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer func() {
		signal.Stop(sigChan)
		globalCancel()
	}()

	// init logger
	logFileName := fmt.Sprintf("/var/log/e4_c2ae.log")
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		fmt.Printf("[ERROR] logs: unable to open file '%v' to write logs: %v\n", logFileName, err)
		fmt.Print("[WARN] logs: falling back to standard output only\n")
		logFile = os.Stdout
	}

	defer logFile.Close()

	logger := log.NewJSONLogger(logFile)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	defer logger.Log("msg", "goodbye")

	// compatibility for packages that do not understand go-kit logger:
	stdLogger := stdlog.New(log.NewStdlibAdapter(logger), "", 0)

	configResolver, err := slibpath.NewAppPathResolver(os.Args[0])
	if err != nil {
		logger.Log("msg", "failed to create configuration resolver", "error", err)
		exitCode = 1
		return
	}

	configLoader := slibcfg.NewViperLoader("config", configResolver)

	appConfig := config.NewAPI()
	if err := configLoader.Load(appConfig.ViperCfgFields()); err != nil {
		logger.Log("msg", "failed to load configuration", "error", err)
		exitCode = 1
		return
	}

	if err := appConfig.Validate(); err != nil {
		logger.Log("msg", "failed to validate configuration", "error", err)
		exitCode = 1
		return
	}

	logger.Log("msg", "successfully loaded configuration")

	// extend logger to forward log to ES
	if appConfig.ES.LoggingEnabled {
		esClient, err := elastic.NewClient(
			elastic.SetURL(appConfig.ES.URLs...),
			elastic.SetSniff(false),
		)

		if err != nil {
			logger.Log("msg", "failed to create elastic client", "error", err)
			exitCode = 1
			return
		}

		esLogger, err := sliblog.WithElasticSearch(logger, esClient, appConfig.ES.LoggingIndex)
		logger = log.With(esLogger)
		if err != nil {
			logger.Log("msg", "failed to initialize elastic log forwarding", "error", err)
			exitCode = 1
			return
		}

		logger.Log("msg", "elasticsearch log forwarding enabled")
	}

	db, err := models.NewDB(appConfig.DB, stdLogger)
	if err != nil {
		logger.Log("msg", "database creation failed", "error", err)
		exitCode = 1
		return
	}
	defer db.Close()
	logger.Log(append(appConfig.DB.Log(), "msg", "successfully connected to database")...)

	if err := db.Migrate(); err != nil {
		logger.Log("msg", "database migration failed", "error", err)
		exitCode = 1
		return
	}
	converter := models.NewConverter()
	validator := models.NewValidator()

	ruleService := services.NewRuleService(db, validator)
	triggerStateService := services.NewTriggerStateService(db)

	globalErrorChan := make(chan error)

	c2ClientFactory, err := pb.NewC2PbClientFactory(appConfig.C2Endpoint, appConfig.C2Certificate)
	if err != nil {
		logger.Log("msg", "cannot create C2 client factory", "error", err)
		exitCode = 1
		return
	}

	// TODO: we might want to test the connection here and fail if C2 isn't running, or certs are bad...
	// Maybe add a Ping() to the C2 server allowing to test the connection / C2 health without
	// sending a real commands. Otherwise we just establish C2 connections only when sending an actual
	// command, ie: when a rule trigger.
	c2client := services.NewC2(c2ClientFactory)

	eventStreamer := events.NewStreamer(c2client, log.With(logger, "type", "eventStreamer"))

	triggerWatcherFactory := watchers.NewTriggerWatcherFactory(
		events.NewStreamListenerFactory(eventStreamer),
		triggerStateService,
		validator,
		log.With(logger, "type", "triggerWatcher"),
	)
	actionFactory := actions.NewActionFactory(c2client, globalErrorChan, log.With(logger, "type", "ruleAction"))

	ruleWatcherFactory := watchers.NewRuleWatcherFactory(
		ruleService,
		triggerWatcherFactory,
		actionFactory,
		globalErrorChan,
		log.With(logger, "type", "ruleWatcher"),
	)

	automationEngine := engine.NewAutomationEngine(
		ruleService,
		ruleWatcherFactory,
		log.With(logger, "type", "automationEngine"),
	)

	server := api.NewServer(
		appConfig.Server,
		ruleService,
		converter,
		log.With(logger, "type", "apiServer"),
	)

	if err := monitoring.Setup(appConfig.OpencensusAddress, appConfig.OpencensusSampleAll); err != nil {
		logger.Log("msg", "failed to setup monitoring", "error", err)
		exitCode = 1
		return
	}

	go func() {
		select {
		case <-sigChan:
			logger.Log("msg", "shutdown requested, cancelling context")
			globalCancel()
		case <-globalCtx.Done():
		}
	}()

	// Start automation engine. Will start a background routine for every rules
	// and every rule's triggers until engineCtx or globalCtx get cancelled.
	engineCtx, engineCancel := context.WithCancel(globalCtx)
	if err := automationEngine.Start(engineCtx); err != nil {
		logger.Log("msg", "error when starting automation engine", "error", err)
		exitCode = 1
		return
	}
	logger.Log("msg", "started automation engine")

	// Start C2AE api server
	go func() {
		err := server.ListenAndServe(globalCtx)
		logger.Log("msg", "failed to listen and serve api", "error", err)
		globalCancel()
	}()

	// Start event stream from C2 server.
	// In case the C2 is not available, or crash after some time
	// the event streamer will try to reconnect every seconds until it succeed
	// or the context get canceled.
	go func() {
		for {
			err := eventStreamer.StartStream(globalCtx)
			logger.Log("msg", "event streamer stopped", "error", err)
			select {
			case <-globalCtx.Done():
				return
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Listen for changes in the database and stop / restart the automation engine,
	// creating a fresh engineCtx.
	// Cancelling the globalCtx stop it from restarting indefinitely.
	for {
		select {
		case <-server.RulesModifiedChan():
			logger.Log("msg", "rules modified, restarting automation engine!")

			engineCancel()

			engineCtx, engineCancel = context.WithCancel(globalCtx)
			if err := automationEngine.Start(engineCtx); err != nil {
				logger.Log("msg", "failed to restart automation engine", "error", err)
			}
			logger.Log("msg", "restarted automation engine")
		case err := <-globalErrorChan:
			logger.Log("msg", "a goroutine emitted an error", "error", err)

		case <-globalCtx.Done():
			engineCancel()
			return
		}
	}
}

func printVersion() {
	if len(gitTag) == 0 {
		fmt.Printf("E4: C2 automation engine api - version %s-%s\n", buildDate, gitCommit)
	} else {
		fmt.Printf("E4: C2 automation engine api - version %s (%s-%s)\n", gitTag, buildDate, gitCommit)
	}
	fmt.Println("Copyright (c) Teserakt AG, 2018-2019")
}
