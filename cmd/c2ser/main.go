package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/jinzhu/gorm"
)

var gitCommit string
var buildDate string

var db *gorm.DB

// RuleClient is rule for client key rotation
type RuleClient struct {
	ID         int    `gorm:"primary_key:true"`
	E4IDAlias  string `gorm:"unique;NOT NULL"`
	KeyPeriod  int    `gorm:"NOT NULL"`
	LastUpdate uint64 `gorm:"NOT NULL"`
}

// RuleTopic is rule for topic key rotation
type RuleTopic struct {
	ID         int    `gorm:"primary_key:true"`
	Topic      string `gorm:"unique;NOT NULL"`
	KeyPeriod  int    `gorm:"NOT NULL"`
	LastUpdate uint64 `gorm:"NOT NULL"`
}

// parseRules creates rules objects from the file, aborting upon an invalid line
func parseRules(rulesFilePath string) (rc []RuleClient, rt []RuleTopic, err error) {

	log.Printf("\nparsing %v...", rulesFilePath)

	file, err := os.Open(rulesFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip blank lines
		if len(line) == 0 {
			continue
		}

		splitLine := strings.Split(line, " ")
		if len(splitLine) != 3 {
			return nil, nil, errors.New("a line does not have 3 space-separated values")
		}

		keyPeriod, err := strconv.Atoi(splitLine[2])
		if err != nil || keyPeriod < 0 || keyPeriod > 9000 {
			return nil, nil, errors.New("key period must be a positive integer less than 9000")
		}

		if splitLine[0] == "C" {
			e4IDAlias := splitLine[1]
			r := RuleClient{E4IDAlias: e4IDAlias, KeyPeriod: keyPeriod}
			rc = append(rc, r)
		} else if splitLine[0] == "T" {
			r := RuleTopic{Topic: splitLine[1], KeyPeriod: keyPeriod}
			rt = append(rt, r)
		} else {
			return nil, nil, errors.New("first value must be C or T")
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	log.Printf("found %v client rules and %v topic rules", len(rc), len(rt))

	return rc, rt, nil
}

func updateRulesClient(rules []RuleClient) {

	for _, rule := range rules {

		// if frequency is 0 then remove any entry for this client
		if rule.KeyPeriod == 0 {
			var oldrule RuleClient
			res := db.Where(&RuleClient{E4IDAlias: rule.E4IDAlias}).First(&oldrule)
			if !gorm.IsRecordNotFoundError(res.Error) {
				if res = db.Delete(&oldrule); res.Error != nil {
					log.Printf("ERROR\tfailed to delete rule for client %v: %v", oldrule.E4IDAlias, res.Error)
				} else {
					log.Printf("OK\trule deleted for client %v", oldrule.E4IDAlias)
				}
			} else {
				log.Printf("WARN\trule with Period=0 entered but no rule to delete for client %v", rule.E4IDAlias)
			}
			continue
		}

		// otherwise insert rule, if doesn't already exists
		var somerule RuleClient
		res := db.Where(&RuleClient{E4IDAlias: rule.E4IDAlias}).First(&somerule)
		if gorm.IsRecordNotFoundError(res.Error) {

			if res := db.Create(&rule); res.Error != nil {
				log.Printf("ERROR\tfailed to add rule for client %v: %v", rule.E4IDAlias, res.Error)
			} else {
				log.Printf("OK\tkey period set to %v for client %v", rule.KeyPeriod, rule.E4IDAlias)
			}
		} else {
			log.Printf("WARN\ta rule already exists for client %v, to remove it set Period=0", rule.E4IDAlias)
		}
	}
}

func updateRulesTopic(rules []RuleTopic) error {

	return nil
}

func main() {

	log.SetFlags(0)
	var err error

	// take a list of files as arguments
	fileNames := os.Args[1:]
	if len(fileNames) == 0 {
		log.Fatal("please provide e4 script files as arguments")
	}

	// open db
	db, err = gorm.Open("sqlite3", "/tmp/e4se.sqlite")
	if err != nil {
		log.Fatalf("FATAL: database opening failed: %v", err)
	}
	defer db.Close()

	// initialize db
	if result := db.AutoMigrate(&RuleClient{}); result.Error != nil {
		log.Fatalf("FATAL: database initialization failed: %v", result.Error)
	}
	if result := db.AutoMigrate(&RuleTopic{}); result.Error != nil {
		log.Fatalf("FATAL: database initialization failed: %v", result.Error)
	}
	db.LogMode(false)

	fmt.Printf("E4: C2 script reader - version %s-%s\n", buildDate, gitCommit[:4])
	fmt.Println("Copyright (c) Teserakt AG, 2018")

	for _, fileName := range fileNames {
		rulesClient, rulesTopic, err := parseRules(fileName)
		if err != nil {
			log.Printf("ERROR\tinvalid script file: %v", err.Error())
			continue
		}
		updateRulesClient(rulesClient)
		updateRulesTopic(rulesTopic)
	}
}
