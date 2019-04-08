package services

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"

	"gitlab.com/teserakt/c2se/internal/models"
	"gitlab.com/teserakt/c2se/internal/pb"
)

func getTestDB(t *testing.T) (models.Database, func()) {
	f, err := ioutil.TempFile(os.TempDir(), "ruleServiceTestDb-")
	if err != nil {
		t.Fatalf("Cannot create temporary file: %s", err)
	}

	db, err := models.NewDB(models.DBConfig{
		Dialect:   models.DBDialectSQLite,
		CnxString: f.Name(),
		LogMode:   true,
		Models:    models.All,
	})

	if err != nil {
		t.Fatalf("Cannot open database: %s", err)
	}

	return db, func() {
		db.Close()
		f.Close()
		os.Remove(f.Name())
	}
}

func TestRuleService(t *testing.T) {

	t.Run("All returns all rules", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewRuleService(db)

		rules, err := srv.All()
		if err != nil {
			t.Errorf("Expected error to be nil, got %s", err)
		}

		if len(rules) != 0 {
			t.Errorf("Expected 0 rules, got %d", len(rules))
		}

		rule1 := models.Rule{
			ActionType:  pb.ActionType_KEY_ROTATION,
			Description: "rule1",
			Targets: []models.Target{
				models.Target{ID: 1},
			},
			Triggers: []models.Trigger{
				models.Trigger{ID: 1},
			},
		}
		rule2 := models.Rule{
			ActionType:  pb.ActionType_KEY_ROTATION,
			Description: "rule2",
			Targets: []models.Target{
				models.Target{ID: 2},
			},
			Triggers: []models.Trigger{
				models.Trigger{ID: 2},
			},
		}

		err = srv.Save(&rule1)
		if err != nil {
			t.Errorf("Expected nil error, got %s", err)
		}

		err = srv.Save(&rule2)
		if err != nil {
			t.Errorf("Expected nil error, got %s", err)
		}

		rules, err = srv.All()
		if len(rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(rules))
		}

		if reflect.DeepEqual(rules, []models.Rule{rule1, rule2}) == false {
			t.Errorf("Expected first rule to be %#v, got %#v", []models.Rule{rule1, rule2}, rules)
		}
	})

	t.Run("Save create the entity if it doesn't exists and update it if it does", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewRuleService(db)

		rule := models.Rule{
			ActionType:  pb.ActionType_KEY_ROTATION,
			Description: "rule",
			Targets: []models.Target{
				models.Target{ID: 1},
			},
			Triggers: []models.Trigger{
				models.Trigger{ID: 1},
			},
		}

		err := srv.Save(&rule)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if rule.ID != 1 {
			t.Errorf("Expected rule ID to have been set to 1, got %d", rule.ID)
		}

		if rule.Targets[0].RuleID != rule.ID {
			t.Errorf("Expected first target ID to have been set to rule ID %d, got %d", rule.ID, rule.Targets[0].RuleID)
		}
		if rule.Triggers[0].RuleID != rule.ID {
			t.Errorf("Expected first trigger ID to have been set to rule ID %d, got %d", rule.ID, rule.Triggers[0].RuleID)
		}

		expectedDescription := "New description"
		rule.Description = expectedDescription
		err = srv.Save(&rule)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if rule.Description != expectedDescription {
			t.Errorf("Expected rule description to have been updated to %s, got %s", expectedDescription, rule.Description)
		}
		if rule.ID != 1 {
			t.Errorf("Expected rule to have been updated and keep id %d, got new id %d", 1, rule.ID)
		}
	})

	t.Run("ByID returns the proper rule", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewRuleService(db)

		rule1 := models.Rule{
			ActionType:  pb.ActionType_KEY_ROTATION,
			Description: "rule1",
			Targets: []models.Target{
				models.Target{ID: 1},
			},
			Triggers: []models.Trigger{
				models.Trigger{ID: 1},
			},
		}
		rule2 := models.Rule{
			ActionType:  pb.ActionType_KEY_ROTATION,
			Description: "rule2",
			Targets: []models.Target{
				models.Target{ID: 2},
			},
			Triggers: []models.Trigger{
				models.Trigger{ID: 2},
			},
		}

		if err := srv.Save(&rule1); err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}
		if err := srv.Save(&rule2); err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		rule, err := srv.ByID(1)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if reflect.DeepEqual(rule, rule1) == false {
			t.Errorf("Expected rule 1 to be %#v, got %#v", rule1, rule)
		}

		rule, err = srv.ByID(2)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if reflect.DeepEqual(rule, rule2) == false {
			t.Errorf("Expected rule 2 to be %#v, got %#v", rule2, rule)
		}

		_, err = srv.ByID(3)
		if err == nil {
			t.Errorf("Expected err to be %s, got %s", gorm.ErrRecordNotFound, err)
		}

	})
}
