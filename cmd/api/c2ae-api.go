package main

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"os/signal"

	"github.com/go-kit/kit/log"

	"gitlab.com/teserakt/c2ae/internal/api"
	"gitlab.com/teserakt/c2ae/internal/config"
	"gitlab.com/teserakt/c2ae/internal/engine"
	"gitlab.com/teserakt/c2ae/internal/engine/actions"
	"gitlab.com/teserakt/c2ae/internal/engine/watchers"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/monitoring"
	"gitlab.com/teserakt/c2ae/internal/pb"
	"gitlab.com/teserakt/c2ae/internal/services"
	slibcfg "gitlab.com/teserakt/serverlib/config"
	slibpath "gitlab.com/teserakt/serverlib/path"
)

// Provided by build script
var gitCommit string
var gitTag string
var buildDate string

func main() {
	exitCode := 0
	defer os.Exit(exitCode)

	printVersion()

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

	ruleService := services.NewRuleService(db)

	globalErrorChan := make(chan error)

	c2ClientFactory, err := pb.NewC2PbClientFactory(appConfig.C2Endpoint, appConfig.C2Certificate)
	if err != nil {
		logger.Log("msg", "cannot create C2 client factory", "error", err)
		exitCode = 1
		return
	}

	c2Requester := services.NewC2Requester(c2ClientFactory)
	// TODO: we might want to test the connection here and fail if C2 isn't running, or certs are bad...
	// Maybe add a Ping() to the C2 server allowing to test the connection / C2 health without
	// sending a real commands. Otherwise we just establish C2 connections only when sending an actual
	// command, ie: when a rule trigger.
	c2client := services.NewC2(c2Requester)

	triggerWatcherFactory := watchers.NewTriggerWatcherFactory(log.With(logger, "type", "triggerWatcher"))
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
		appConfig.Addr,
		ruleService,
		converter,
		log.With(logger, "type", "apiServer"),
	)

	globalCtx, globalCancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer func() {
		signal.Stop(sigChan)
		globalCancel()
	}()

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

	engineCtx, engineCancel := context.WithCancel(globalCtx)
	if err := automationEngine.Start(engineCtx); err != nil {
		logger.Log("msg", "error when starting automation engine", "error", err)
		exitCode = 1
		return
	}
	logger.Log("msg", "started automation engine")

	go func() {
		err := server.ListenAndServe(globalCtx)
		logger.Log("msg", "failed to listen and serve api", "error", err)
		globalCancel()
		return
	}()

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
