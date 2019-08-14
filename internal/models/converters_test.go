package models

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"gitlab.com/teserakt/c2ae/internal/pb"
)

func TestConverter(t *testing.T) {

	converter := NewConverter()

	trigger1 := Trigger{
		ID:          1,
		Settings:    []byte("settings1"),
		TriggerType: pb.TriggerType_EVENT,
	}
	trigger2 := Trigger{
		ID:          2,
		Settings:    []byte("settings2"),
		TriggerType: pb.TriggerType_TIME_INTERVAL,
	}
	trigger3 := Trigger{
		ID:          3,
		Settings:    []byte("settings3"),
		TriggerType: pb.TriggerType_EVENT,
	}

	target1 := Target{
		ID:   1,
		Expr: "expr1",
		Type: pb.TargetType_ANY,
	}
	target2 := Target{
		ID:   2,
		Expr: "expr2",
		Type: pb.TargetType_TOPIC,
	}
	target3 := Target{
		ID:   3,
		Expr: "expr3",
		Type: pb.TargetType_CLIENT,
	}

	rule1 := Rule{
		ID:           1,
		Description:  "description1",
		ActionType:   pb.ActionType_KEY_ROTATION,
		LastExecuted: time.Now(),
		Targets:      []Target{target1, target2},
		Triggers:     []Trigger{trigger1, trigger2},
	}
	rule2 := Rule{
		ID:           2,
		Description:  "description2",
		ActionType:   pb.ActionType_KEY_ROTATION,
		LastExecuted: time.Now(),
		Targets:      []Target{target1, target2, target3},
		Triggers:     []Trigger{trigger1, trigger2, trigger3},
	}

	t.Run("RulesToPb and PbToRules properly converts []models.Rule to []*pb.Rule and back", func(t *testing.T) {
		rules := []Rule{rule1, rule2}
		pbRules, err := converter.RulesToPb(rules)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		if len(rules) != len(pbRules) {
			t.Errorf("Expected %d converted rules, got %d", len(rules), len(pbRules))
		}

		for i, rule := range rules {
			assertSameRule(t, rule, pbRules[i])
		}

		origRules, err := converter.PbToRules(pbRules)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		if len(rules) != len(origRules) {
			t.Errorf("Expected %d rules once converted back, got %d", len(rules), len(origRules))
		}

		for i, rule := range rules {
			if rule.ID != origRules[i].ID {
				t.Errorf("Expected rule id to be %d, got %d", rule.ID, origRules[i].ID)
			}
			if rule.Description != origRules[i].Description {
				t.Errorf("Expected rule description to be %s, got %s", rule.Description, origRules[i].Description)
			}
			if rule.ActionType != origRules[i].ActionType {
				t.Errorf("Expected rule action type to be %v, got %v", rule.ActionType, origRules[i].ActionType)
			}
			if rule.LastExecuted.UnixNano() != origRules[i].LastExecuted.UnixNano() {
				t.Errorf("Expected last executed to be %#v, got %#v", rule.LastExecuted, origRules[i].LastExecuted)
			}
			if reflect.DeepEqual(rule.Targets, origRules[i].Targets) == false {
				t.Errorf("Expected targets to be %#v, got %#v", rule.Targets, origRules[i].Targets)
			}
			if reflect.DeepEqual(rule.Triggers, origRules[i].Triggers) == false {
				t.Errorf("Expected triggers to be %#v, got %#v", rule.Triggers, origRules[i].Triggers)
			}
		}
	})
}

func assertSameRule(t *testing.T, rule Rule, pbRule *pb.Rule) {
	if rule.ID != int(pbRule.Id) {
		t.Errorf("Expected rule ID to be %d, got %d", rule.ID, pbRule.Id)
	}

	if rule.Description != pbRule.Description {
		t.Errorf("Expected rule description to be %s, got %s", rule.Description, pbRule.Description)
	}
	time, err := ptypes.Timestamp(pbRule.LastExecuted)
	if err != nil {
		t.Errorf("Converted rule have an invalid timestamp: %s", err)
	}
	if rule.LastExecuted.UnixNano() != time.UnixNano() {
		t.Errorf("Expected rule last executed to be %v, got %v", rule.LastExecuted, pbRule.LastExecuted)
	}

	if len(rule.Triggers) != len(pbRule.Triggers) {
		t.Errorf("Expected %d rule triggers, got %d", len(rule.Triggers), len(pbRule.Triggers))
	}
	for i, rTrig := range rule.Triggers {
		assertSameTrigger(t, rTrig, pbRule.Triggers[i])
	}

	if len(rule.Targets) != len(pbRule.Targets) {
		t.Errorf("Expected %d rule targets, got %d", len(rule.Targets), len(pbRule.Targets))
	}
	for i, rTarg := range rule.Targets {
		assertSameTarget(t, rTarg, pbRule.Targets[i])
	}
}

func assertSameTrigger(t *testing.T, trigger Trigger, pbTrigger *pb.Trigger) {
	if trigger.ID != int(pbTrigger.Id) {
		t.Errorf("Expected trigger id to be %d, got %d", trigger.ID, pbTrigger.Id)
	}

	if trigger.TriggerType != pbTrigger.Type {
		t.Errorf("Expected trigger type to be %v, got %v", trigger.TriggerType, pbTrigger.Type)
	}

	if bytes.Equal(trigger.Settings, pbTrigger.Settings) == false {
		t.Errorf("Expected trigger settings to be %v, got %v", trigger.Settings, pbTrigger.Settings)
	}
}

func assertSameTarget(t *testing.T, target Target, pbTarget *pb.Target) {
	if target.ID != int(pbTarget.Id) {
		t.Errorf("Expected target id to be %d, got %d", target.ID, pbTarget.Id)
	}

	if target.Expr != pbTarget.Expr {
		t.Errorf("Expected target expr to be %s, got %s", target.Expr, pbTarget.Expr)
	}

	if target.Type != pbTarget.Type {
		t.Errorf("Expected target type to be %v, got %v", target.Type, pbTarget.Type)
	}
}
