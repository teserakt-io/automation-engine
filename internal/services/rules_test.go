package services

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/jinzhu/gorm"

	"gitlab.com/teserakt/c2ae/internal/config"
	"gitlab.com/teserakt/c2ae/internal/models"
	"gitlab.com/teserakt/c2ae/internal/pb"
	slibcfg "gitlab.com/teserakt/serverlib/config"
)

func TestRuleService(t *testing.T) {
	testRuleServiceDatabase(t, sqliteTestDB)

	if os.Getenv("C2AETEST_POSTGRES") == "" {
		t.Skip("C2AETEST_POSTGRES environment is not set")

		return
	}
	testRuleServiceDatabase(t, postgresTestDB)

}

func sqliteTestDB(t *testing.T) (models.Database, func()) {

	f, err := ioutil.TempFile(os.TempDir(), "ruleServiceTestDb-")
	if err != nil {
		t.Fatalf("Cannot create temporary file: %s", err)
	}

	logger := log.New(os.Stdout, "", 0)

	dbConfig := config.DBCfg{
		Type:    slibcfg.DBTypeSQLite,
		File:    f.Name(),
		Logging: false,
	}

	db, err := models.NewDB(dbConfig, logger)

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

func postgresTestDB(t *testing.T) (models.Database, func()) {
	logger := log.New(os.Stdout, "", 0)

	dbConfig := config.DBCfg{
		Type:             slibcfg.DBTypePostgres,
		Username:         "c2ae_test",
		Password:         "teserakte4",
		Schema:           "c2ae_test_unit",
		SecureConnection: slibcfg.DBSecureConnectionInsecure,
		Passphrase:       "unittest-passphrase",
		Host:             "127.0.0.1",
		Database:         "e4",
		Logging:          false,
	}

	db, err := models.NewDB(dbConfig, logger)
	if err != nil {
		t.Fatalf("Cannot open database: %s", err)
	}
	db.Connection().Exec("CREATE SCHEMA c2ae_test_unit AUTHORIZATION c2ae_test;")

	if err := db.Migrate(); err != nil {
		t.Fatalf("Expected no error when migrating database, got %v", err)
	}

	return db, func() {
		db.Connection().Exec("DROP SCHEMA c2ae_test_unit CASCADE;")
		db.Close()
	}
}

func createRules(t *testing.T, srv RuleService, validator *models.MockValidator) (rule1 models.Rule, rule2 models.Rule) {
	rule1 = models.Rule{
		ActionType:  pb.ActionType_KEY_ROTATION,
		Description: "rule1",
		Targets: []models.Target{
			models.Target{
				ID:   1,
				Type: pb.TargetType_CLIENT,
				Expr: "target1Expr",
			},
		},
		Triggers: []models.Trigger{
			models.Trigger{
				ID:          1,
				TriggerType: pb.TriggerType_TIME_INTERVAL,
				Settings:    []byte("settings1"),
			},
		},
	}

	rule2 = models.Rule{
		ActionType:  pb.ActionType_KEY_ROTATION,
		Description: "rule2",
		Targets: []models.Target{
			models.Target{
				ID:   2,
				Type: pb.TargetType_TOPIC,
				Expr: "target2Expr",
			},
		},
		Triggers: []models.Trigger{
			models.Trigger{
				ID:          2,
				TriggerType: pb.TriggerType_EVENT,
				Settings:    []byte("settings2"),
			},
		},
	}

	ctx := context.Background()

	validator.EXPECT().ValidateRule(rule1).Times(1)
	err := srv.Save(ctx, &rule1)
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}

	validator.EXPECT().ValidateRule(rule2).Times(1)
	err = srv.Save(ctx, &rule2)
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}

	return rule1, rule2
}

func testRuleServiceDatabase(t *testing.T, getTestDB func(t *testing.T) (models.Database, func())) {
	ctx := context.Background()

	t.Run("All returns all rules", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)

		rules, err := srv.All(ctx)
		if err != nil {
			t.Errorf("Expected error to be nil, got %s", err)
		}

		if len(rules) != 0 {
			t.Errorf("Expected 0 rules, got %d", len(rules))
		}

		rule1, rule2 := createRules(t, srv, validator)

		rules, err = srv.All(ctx)
		if len(rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(rules))
		}

		// Ignore lastExecuted time to simplify next assertion
		// which fail on postgres (too slow)
		for i := range rules {
			rules[i].LastExecuted = time.Time{}
		}

		if reflect.DeepEqual(rules, []models.Rule{rule1, rule2}) == false {
			t.Errorf("Expected rules to be %#v, got %#v", []models.Rule{rule1, rule2}, rules)
		}
	})

	t.Run("Save create the entity if it doesn't exists and update it if it does", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)

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

		validator.EXPECT().ValidateRule(rule).Times(1).Return(nil)

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

		validator.EXPECT().ValidateRule(rule).Times(1).Return(nil)
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

	t.Run("Save with invalid rules returns a validation error", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)

		rule := models.Rule{}

		validationError := errors.New("validation error")
		validator.EXPECT().ValidateRule(rule).Times(1).Return(validationError)

		err := srv.Save(ctx, &rule)
		if err == nil {
			t.Error("Expected an error, got nil")
		}

		if !strings.Contains(err.Error(), validationError.Error()) {
			t.Errorf("Expected err to contains %v, got %v", validationError, err)
		}
	})

	t.Run("ByID returns the proper rule", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)

		rule1, rule2 := createRules(t, srv, validator)

		rule, err := srv.ByID(ctx, rule1.ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		// Ignore lastExecuted time to simplify next assertion
		// which fail on postgres (too slow)
		rule.LastExecuted = time.Time{}

		if reflect.DeepEqual(rule, rule1) == false {
			t.Errorf("Expected rule 1 to be %#v, got %#v", rule1, rule)
		}

		rule, err = srv.ByID(ctx, rule2.ID)
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		// Ignore lastExecuted time to simplify next assertion
		// which fail on postgres (too slow)
		rule.LastExecuted = time.Time{}
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

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)

		rule1, rule2 := createRules(t, srv, validator)

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

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)

		rule1, rule2 := createRules(t, srv, validator)

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

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)

		rule1, rule2 := createRules(t, srv, validator)

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

	t.Run("DeleteTargets properly delete given targets", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)
		rule1, _ := createRules(t, srv, validator)

		originalTargets := make([]models.Target, len(rule1.Targets))
		copy(originalTargets, rule1.Targets)

		targets := []models.Target{
			models.Target{ID: 1000},
			models.Target{ID: 1001},
		}

		rule1.Targets = append(rule1.Targets, targets...)

		validator.EXPECT().ValidateRule(rule1).Times(1)
		if err := srv.Save(ctx, &rule1); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		err := srv.DeleteTargets(ctx, targets...)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		r, err := srv.ByID(ctx, rule1.ID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if reflect.DeepEqual(originalTargets, r.Targets) == false {
			t.Errorf("Expected Targets to be %#v, got %#v", originalTargets, r.Targets)
		}
	})

	t.Run("DeleteTriggers properly delete given triggers", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		validator := models.NewMockValidator(mockCtrl)

		srv := NewRuleService(db, validator)
		rule1, _ := createRules(t, srv, validator)

		originalTriggers := make([]models.Trigger, len(rule1.Triggers))
		copy(originalTriggers, rule1.Triggers)

		triggers := []models.Trigger{
			models.Trigger{ID: 1000},
			models.Trigger{ID: 1001},
		}

		rule1.Triggers = append(rule1.Triggers, triggers...)

		validator.EXPECT().ValidateRule(rule1).Times(1)
		if err := srv.Save(ctx, &rule1); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		err := srv.DeleteTriggers(ctx, triggers...)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		r, err := srv.ByID(ctx, rule1.ID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if reflect.DeepEqual(originalTriggers, r.Triggers) == false {
			t.Errorf("Expected Triggers to be %#v, got %#v", originalTriggers, r.Triggers)
		}
	})
}
