package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"

	"gitlab.com/teserakt/c2ae/internal/config"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
	slibcfg "gitlab.com/teserakt/serverlib/config"
)

func getTestDB(t *testing.T) (models.Database, func()) {
	f, err := ioutil.TempFile(os.TempDir(), "ruleServiceTestDb-")
	if err != nil {
		t.Fatalf("Cannot create temporary file: %s", err)
	}

	logger := log.New(os.Stdout, "", 0)

	db, err := models.NewDB(config.DBCfg{
		Type:    slibcfg.DBTypeSQLite,
		File:    f.Name(),
		Logging: false,
	}, logger)

	if err != nil {
		t.Fatalf("Cannot open database: %s", err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatalf("Cannot migrate database: %s", err)
	}

	return db, func() {
		db.Close()
		f.Close()
		os.Remove(f.Name())
	}
}

func createRules(t *testing.T, srv RuleService) (rule1 models.Rule, rule2 models.Rule) {
	rule1 = models.Rule{
		ActionType:  pb.ActionType_KEY_ROTATION,
		Description: "rule1",
		Targets: []models.Target{
			models.Target{ID: 1},
		},
		Triggers: []models.Trigger{
			models.Trigger{ID: 1},
		},
	}

	rule1.Targets[0].Rule = &rule1
	rule1.Triggers[0].Rule = &rule1

	rule2 = models.Rule{
		ActionType:  pb.ActionType_KEY_ROTATION,
		Description: "rule2",
		Targets: []models.Target{
			models.Target{ID: 2},
		},
		Triggers: []models.Trigger{
			models.Trigger{ID: 2},
		},
	}

	rule2.Targets[0].Rule = &rule2
	rule2.Triggers[0].Rule = &rule2

	ctx := context.Background()

	err := srv.Save(ctx, &rule1)
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}

	err = srv.Save(ctx, &rule2)
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}

	return rule1, rule2
}

func TestRuleService(t *testing.T) {
	ctx := context.Background()

	t.Run("All returns all rules", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewRuleService(db)

		rules, err := srv.All(ctx)
		if err != nil {
			t.Errorf("Expected error to be nil, got %s", err)
		}

		if len(rules) != 0 {
			t.Errorf("Expected 0 rules, got %d", len(rules))
		}

		rule1, rule2 := createRules(t, srv)

		rules, err = srv.All(ctx)
		if len(rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(rules))
		}

		if reflect.DeepEqual(rules, []models.Rule{rule1, rule2}) == false {
			t.Errorf("Expected rules to be %#v, got %#v", []models.Rule{rule1, rule2}, rules)
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

		err := srv.Save(ctx, &rule)
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
		err = srv.Save(ctx, &rule)
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

		rule1, rule2 := createRules(t, srv)

		rule, err := srv.ByID(ctx, rule1.ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if reflect.DeepEqual(rule, rule1) == false {
			t.Errorf("Expected rule 1 to be %#v, got %#v", rule1, rule)
		}

		rule, err = srv.ByID(ctx, rule2.ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if reflect.DeepEqual(rule, rule2) == false {
			t.Errorf("Expected rule 2 to be %#v, got %#v", rule2, rule)
		}

		_, err = srv.ByID(ctx, 3)
		if err == nil {
			t.Errorf("Expected err to be %s, got %s", gorm.ErrRecordNotFound, err)
		}
	})

	t.Run("Delete removes the rule and dependancies from database", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewRuleService(db)

		rule1, rule2 := createRules(t, srv)

		if err := srv.Delete(ctx, rule1); err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		_, err := srv.ByID(ctx, rule1.ID)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("Expected err to be %s, got %s", gorm.ErrRecordNotFound, err)
		}

		_, err = srv.TriggerByID(ctx, rule1.Triggers[0].ID)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("Expected err to be %s, got %s", gorm.ErrRecordNotFound, err)
		}

		_, err = srv.TargetByID(ctx, rule1.Targets[0].ID)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("Expected err to be %s, got %s", gorm.ErrRecordNotFound, err)
		}

		_, err = srv.ByID(ctx, rule2.ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}
		_, err = srv.TriggerByID(ctx, rule2.Triggers[0].ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		_, err = srv.TargetByID(ctx, rule2.Targets[0].ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}
	})

	t.Run("TriggerByID retrieve proper trigger", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewRuleService(db)

		rule1, rule2 := createRules(t, srv)

		trigger, err := srv.TriggerByID(ctx, rule1.Triggers[0].ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if reflect.DeepEqual(trigger, rule1.Triggers[0]) == false {
			t.Errorf("Expected trigger to be %#v, got %#v", rule1.Triggers[0], trigger)
		}

		trigger, err = srv.TriggerByID(ctx, rule2.Triggers[0].ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if reflect.DeepEqual(trigger, rule2.Triggers[0]) == false {
			t.Errorf("Expected trigger to be %#v, got %#v", rule2.Triggers[0], trigger)
		}

		_, err = srv.TriggerByID(ctx, 3)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("Expected err to be %s, got %s", gorm.ErrRecordNotFound, err)
		}
	})

	t.Run("TargetByID retrieve proper target", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewRuleService(db)

		rule1, rule2 := createRules(t, srv)

		target, err := srv.TargetByID(ctx, rule1.Targets[0].ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if reflect.DeepEqual(target, rule1.Targets[0]) == false {
			t.Errorf("Expected trigger to be %#v, got %#v", rule1.Targets[0], target)
		}

		target, err = srv.TargetByID(ctx, 2)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if reflect.DeepEqual(target, rule2.Targets[0]) == false {
			t.Errorf("Expected trigger to be %#v, got %#v", rule2.Targets[0], target)
		}

		_, err = srv.TargetByID(ctx, 3)
		if err != gorm.ErrRecordNotFound {
			t.Errorf("Expected err to be %s, got %s", gorm.ErrRecordNotFound, err)
		}
	})
}
