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
		Models: []interface{}{
			models.Trigger{},
			models.Target{},
			models.Rule{},
		},
	}

	db, err := models.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("FATAL: database creation failed: %v", err)
	}
	defer db.Close()

	ruleService := services.NewRuleService(db)

	server := api.NewServer(
		":8080",
		ruleService,
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
