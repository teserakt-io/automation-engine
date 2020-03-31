// Copyright 2020 Teserakt AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"testing"

	"github.com/teserakt-io/automation-engine/internal/pb"
)

func TestValidator(t *testing.T) {
	validator := NewValidator()

	t.Run("ValidateTriggers properly returns error with bad trigger", func(t *testing.T) {
		badTriggerDataset := []struct {
			Trigger       Trigger
			ExpectedError error
		}{
			{Trigger: Trigger{}, ExpectedError: ErrUndefinedTriggerType},
			{Trigger: Trigger{TriggerType: pb.TriggerType_UNDEFINED_TRIGGER}, ExpectedError: ErrUndefinedTriggerType},
			{Trigger: Trigger{TriggerType: pb.TriggerType_TIME_INTERVAL}},
			{Trigger: Trigger{TriggerType: pb.TriggerType(-1)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_TIME_INTERVAL, Settings: []byte(`not_even_json`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_TIME_INTERVAL, Settings: []byte(`{"something":"bad"}`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_TIME_INTERVAL, Settings: []byte(`{"expr":"bad"}`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_TIME_INTERVAL, Settings: []byte(`{"expr":"* * *"}`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"bad":"json"}`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"maxOccurrence": 0}`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"maxOccurrence": -1}`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"maxOccurrence": -1}`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"eventType", "", maxOccurrence": 1}`)}},
			{Trigger: Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"eventType", "NOT_VALID", maxOccurrence": 1}`)}},
		}

		for _, testData := range badTriggerDataset {
			err := validator.ValidateTrigger(testData.Trigger)
			if err == nil {
				t.Errorf("Expected trigger %#v to produce a validation error, got nil", testData.Trigger)
			}

			if testData.ExpectedError != nil && err != testData.ExpectedError {
				t.Errorf("Expected error to be %v, got %v", testData.ExpectedError, err)
			}
		}
	})

	t.Run("ValidateTrigger does not returns error with valid triggers", func(t *testing.T) {
		validTriggers := []Trigger{
			Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"eventType": "CLIENT_UNSUBSCRIBED", "maxOccurrence": 1}`)},
			Trigger{TriggerType: pb.TriggerType_TIME_INTERVAL, Settings: []byte(`{"expr":"* * * * * *"}`)},
			Trigger{TriggerType: pb.TriggerType_TIME_INTERVAL, Settings: []byte(`{"expr":"* 0/3 * * * *"}`)},
		}

		for _, trigger := range validTriggers {
			err := validator.ValidateTrigger(trigger)
			if err != nil {
				t.Errorf("Expected no error when validating trigger %#v, got %v", trigger, err)
			}
		}
	})

	t.Run("ValidateTarget properly returns error with bad targets", func(t *testing.T) {
		badTargetDataset := []struct {
			Target        Target
			ExpectedError error
		}{
			{Target: Target{}, ExpectedError: ErrTargetExprRequired},
			{Target: Target{Expr: "bad(regexp"}},
		}

		for _, testData := range badTargetDataset {
			err := validator.ValidateTarget(testData.Target)
			if err == nil {
				t.Errorf("Expected target %#v to produce a validation error, got nil", testData.Target)
			}

			if testData.ExpectedError != nil && err != testData.ExpectedError {
				t.Errorf("Expected error to be %v, got %v", testData.ExpectedError, err)
			}
		}
	})

	t.Run("ValidateTarget does not returns error with valid targets", func(t *testing.T) {
		validTargets := []Target{
			Target{Expr: "client1"},
			Target{Expr: ".*"},
			Target{Expr: "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]$"},
		}

		for _, target := range validTargets {
			err := validator.ValidateTarget(target)
			if err != nil {
				t.Errorf("Expected no error when validating target %#v, got %v", target, err)
			}
		}
	})

	t.Run("ValidateRule properly returns error with bad rules", func(t *testing.T) {
		badRuleDataset := []struct {
			Rule          Rule
			ExpectedError error
		}{
			{Rule: Rule{}, ExpectedError: ErrUndefinedAction},
			{Rule: Rule{ActionType: pb.ActionType_UNDEFINED_ACTION}, ExpectedError: ErrUndefinedAction},
			{Rule: Rule{ActionType: pb.ActionType(-1)}, ExpectedError: ErrUnknownActionType},
			{Rule: Rule{ActionType: pb.ActionType_KEY_ROTATION, Triggers: []Trigger{Trigger{}}}},
			{Rule: Rule{ActionType: pb.ActionType_KEY_ROTATION, Targets: []Target{Target{}}}},
		}

		for _, testData := range badRuleDataset {
			err := validator.ValidateRule(testData.Rule)
			if err == nil {
				t.Errorf("Expected rule %#v to produce a validation error, got nil", testData.Rule)
			}

			if testData.ExpectedError != nil && err != testData.ExpectedError {
				t.Errorf("Expected error to be %v, got %v", testData.ExpectedError, err)
			}
		}
	})

	t.Run("ValidateRule does not returns error with valid rules", func(t *testing.T) {
		validRules := []Rule{
			Rule{ActionType: pb.ActionType_KEY_ROTATION},
			Rule{ActionType: pb.ActionType_KEY_ROTATION, Triggers: []Trigger{Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"eventType": "CLIENT_SUBSCRIBED", "maxOccurrence": 1}`)}}},
			Rule{ActionType: pb.ActionType_KEY_ROTATION, Targets: []Target{Target{Expr: "abc"}}},
			Rule{
				ActionType: pb.ActionType_KEY_ROTATION,
				Triggers:   []Trigger{Trigger{TriggerType: pb.TriggerType_EVENT, Settings: []byte(`{"eventType": "CLIENT_SUBSCRIBED", "maxOccurrence": 1}`)}},
				Targets:    []Target{Target{Expr: "abc"}},
			},
		}

		for _, Rule := range validRules {
			err := validator.ValidateRule(Rule)
			if err != nil {
				t.Errorf("Expected no error when validating rule %#v, got %v", Rule, err)
			}
		}
	})
}
