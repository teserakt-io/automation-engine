package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.com/teserakt/c2se/internal/events"

	"gitlab.com/teserakt/c2se/internal/engine/watchers"

	"gitlab.com/teserakt/c2se/internal/api"
	"gitlab.com/teserakt/c2se/internal/config"
	"gitlab.com/teserakt/c2se/internal/engine"
	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/services"
)

// Provided by build script
var gitCommit string
var gitTag string
var buildDate string

func main() {

	defer os.Exit(1)

	printVersion()

	var appConfig config.API
	flag.StringVar(&appConfig.DBFilepath, "db", "", "path to the sqlite db file")
	flag.StringVar(&appConfig.Addr, "addr", "127.0.0.1:5556", "interface:port to listen for incoming connections")
	flag.Parse()

	if err := appConfig.Validate(); err != nil {
		fmt.Printf("\ninvalid settings: %s\n\n", err)
		flag.Usage()

		return
	}

	dbConfig := models.DBConfig{
		Dialect:   models.DBDialectSQLite,
		CnxString: appConfig.DBFilepath,
		LogMode:   false,
	}

	db, err := models.NewDB(dbConfig)
	if err != nil {
		log.Printf("FATAL: database creation failed: %s", err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		log.Printf("FATAL: database migration failed: %s", err)

		return
	}
	converter := models.NewConverter()

	ruleService := services.NewRuleService(db)

	globalErrorChan := make(chan error)

	triggerWatcherFactory := watchers.NewTriggerWatcherFactory()
	ruleWatcherFactory := watchers.NewRuleWatcherFactory(
		ruleService,
		triggerWatcherFactory,
		make(chan events.TriggerEvent),
		globalErrorChan,
	)

	scriptEngine := engine.NewScriptEngine(ruleService, ruleWatcherFactory)

	server := api.NewServer(
		appConfig.Addr,
		ruleService,
		converter,
	)

	err = scriptEngine.Start()
	if err != nil {
		log.Printf("Error when starting script engine: %s", err)

		return
	}

	go server.ListenAndServe(globalErrorChan)

	for {
		select {
		case <-server.RulesModifiedChan():
			log.Println("Rules modified, restarting script engine!")
			scriptEngine.Stop()
			err = scriptEngine.Start()
			if err != nil {
				log.Printf("ERROR: failed to restart script engine: %s", err)
			}
		case err := <-globalErrorChan:
			log.Printf("ERROR: %s", err)
		}
	}
}

func printVersion() {
	if len(gitTag) == 0 {
		fmt.Printf("E4: C2 script reader api - version %s-%s\n", buildDate, gitCommit)
	} else {
		fmt.Printf("E4: C2 script reader api - version %s (%s-%s)\n", gitTag, buildDate, gitCommit)
	}
	fmt.Println("Copyright (c) Teserakt AG, 2018-2019")
}
