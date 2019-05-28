package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-kit/kit/log"

	"gitlab.com/teserakt/c2ae/internal/api"
	"gitlab.com/teserakt/c2ae/internal/config"
	"gitlab.com/teserakt/c2ae/internal/engine"
	"gitlab.com/teserakt/c2ae/internal/engine/actions"
	"gitlab.com/teserakt/c2ae/internal/engine/watchers"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
	"gitlab.com/teserakt/c2ae/internal/services"
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

	var appConfig config.API
	flag.StringVar(&appConfig.DBFilepath, "db", "", "path to the sqlite db file")
	flag.StringVar(&appConfig.Addr, "addr", "localhost:5556", "interface:port to listen for incoming connections")
	flag.StringVar(&appConfig.C2Endpoint, "c2", "localhost:5555", "tcp://interface:port to the c2 backend")
	flag.StringVar(&appConfig.C2Certificate, "c2cert", "", "path to the c2 backend certificate")
	flag.Parse()

	if err := appConfig.Validate(); err != nil {
		logger.Log("msg", "invalid settings", "error", err)
		flag.Usage()
		exitCode = 1
		return
	}
	logger.Log("msg", "successfully loaded configuration")

	dbConfig := models.DBConfig{
		Dialect:   models.DBDialectSQLite,
		CnxString: appConfig.DBFilepath,
		LogMode:   false,
	}

	db, err := models.NewDB(dbConfig)
	if err != nil {
		logger.Log("msg", "database creation failed", "error", err)
		exitCode = 1
		return
	}
	defer db.Close()

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

	go server.ListenAndServe(globalCtx, globalErrorChan)

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
