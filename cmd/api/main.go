package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gitlab.com/teserakt/c2se/internal/api"
	"gitlab.com/teserakt/c2se/internal/config"
	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/services"
)

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
		LogMode:   true,
		Models:    models.All,
	}

	db, err := models.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("FATAL: database creation failed: %v", err)
	}
	defer db.Close()

	ruleService := services.NewRuleService(db)
	converter := models.NewConverter()

	server := api.NewServer(
		appConfig.Addr,
		ruleService,
		converter,
	)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("FATAL: server failed: %s", err)
	}
}

func printVersion() {
	if len(gitTag) == 0 {
		fmt.Printf("E4: C2 script reader - version %s-%s\n", buildDate, gitCommit)
	} else {
		fmt.Printf("E4: C2 script reader - version %s (%s-%s)\n", gitTag, buildDate, gitCommit)
	}
	fmt.Println("Copyright (c) Teserakt AG, 2018-2019")
}
