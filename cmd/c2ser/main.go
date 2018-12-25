package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var gitCommit string
var buildDate string

var db *gorm.DB

// RuleClient is rule for client key rotation
type RuleClient struct {
	ID         int    `gorm:"primary_key:true"`
	Frequency  uint32 `gorm:"NOT NULL"`
	LastUpdate uint64 `gorm:"NOT NULL"`
}

// RuleTopic is rule for topic key rotation
type RuleTopic struct {
	ID         int    `gorm:"primary_key:true"`
	Topic      string `gorm:"NOT NULL"`
	Frequency  uint32 `gorm:"NOT NULL"`
	LastUpdate uint64 `gorm:"NOT NULL"`
}

// parseRules creates rules objects from the file, aborting upon an invalid line
func parseRules(rulesFilePath string) ([]RuleClient, []RuleTopic, error) {

	return nil, nil, nil
}

func updateRulesClient(rules []RuleClient) error {

	for _, rule := range rules {

		// if frequency
		if rule.Frequency == 0 {
			var oldrule RuleClient
			result := db.Where(&RuleClient{})
		}
	}

	return nil
}

func updateRulesTopic(rules []RuleTopic) error {

	return nil
}

func main() {

	log.SetFlags(0)

	// take a list of files as arguments
	fileNames := os.Args[1:]
	if len(fileNames) == 0 {
		log.Fatal("please provide e4 script files as arguments")
	}

	// open db
	db, err := gorm.Open("sqlite3", "/tmp/e4se.sqlite")
	if err != nil {
		log.Fatalf("database opening failed: %v", err)
	}
	defer db.Close()

	// initialize db
	if result := db.AutoMigrate(&RuleClient{}); result.Error != nil {
		log.Fatalf("database initialization failed: %v", result.Error)
	}
	if result := db.AutoMigrate(&RuleTopic{}); result.Error != nil {
		log.Fatalf("database initialization failed: %v", result.Error)
	}

	// for each file:
	// 	- check that it exists
	//  - check that it's valid
	//  - add rules to the db

	log.SetFlags(0)

}
