package main

import (
	"bufio"
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

	log.Printf("parsing file %v...", rulesFilePath)

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
			log.Fatalf("FATAL ERROR: this line does not have 3 space-separated values: %v", line)
		}

		keyPeriod, err := strconv.Atoi(splitLine[2])
		if err != nil || keyPeriod < 0 || keyPeriod > 9000 {
			log.Fatalf("FATAL ERROR: third value (key period) must be a positive integer less than 9000, here it is %v", splitLine[2])
		}

		if splitLine[0] == "C" {
			e4IDAlias := splitLine[1]
			r := RuleClient{E4IDAlias: e4IDAlias, KeyPeriod: keyPeriod}
			rc = append(rc, r)
		} else if splitLine[0] == "T" {
			r := RuleTopic{Topic: splitLine[1], KeyPeriod: keyPeriod}
			rt = append(rt, r)
		} else {
			log.Fatalf("FATAL ERROR: first value must be C or T, here it is %v", splitLine[0])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	log.Printf("file %v parsed: %v client rules and %v topic rules found", rulesFilePath, len(rc), len(rt))

	return rc, rt, nil
}

func updateRulesClient(rules []RuleClient) {

	log.Printf("processing client rules...")

	for _, rule := range rules {

		// if frequency is 0 then remove any entry for this client
		if rule.KeyPeriod == 0 {
			var oldrule RuleClient
			res := db.Where(&RuleClient{E4IDAlias: rule.E4IDAlias}).First(&oldrule)
			if !gorm.IsRecordNotFoundError(res.Error) {
				if res = db.Delete(&oldrule); res.Error != nil {
					log.Printf("ERROR: failed to delete rule for client %v: %v", oldrule.E4IDAlias, res.Error)
				} else {
					log.Printf("OK: rule deleted for client %v", oldrule.E4IDAlias)
				}
			} else {
				log.Printf("WARNING: found 0-period rule, but found no rule to delete for client %v", rule.E4IDAlias)
			}
			continue
		}

		// otherwise insert rule, if doesn't already exists
		if db.NewRecord(rule) {

			if res := db.Create(&rule); res.Error != nil {
				log.Printf("ERROR: failed to add rule for client %v: %v", rule.E4IDAlias, res.Error)
			} else {
				log.Printf("OK: added key period %v for client %v", rule.KeyPeriod, rule.E4IDAlias)
			}
		} else {
			log.Printf("WARNING: rule already exists for client %v, to remove it set Period=0", rule.E4IDAlias)
		}
	}
	log.Printf("client rules processing done!")
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
		log.Fatalf("FATAL ERROR: database opening failed: %v", err)
	}
	defer db.Close()

	// initialize db
	if result := db.AutoMigrate(&RuleClient{}); result.Error != nil {
		log.Fatalf("FATAL ERROR: database initialization failed: %v", result.Error)
	}
	if result := db.AutoMigrate(&RuleTopic{}); result.Error != nil {
		log.Fatalf("FATAL ERROR: database initialization failed: %v", result.Error)
	}
	db.LogMode(false)

	for _, fileName := range fileNames {
		rulesClient, rulesTopic, err := parseRules(fileName)
		if err != nil {
			log.Printf("ERROR: %v", err.Error())
			continue
		}
		updateRulesClient(rulesClient)
		updateRulesTopic(rulesTopic)
	}
}
