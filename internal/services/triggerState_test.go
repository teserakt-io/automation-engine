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

package services

import (
	"context"
	"os"
	reflect "reflect"
	"testing"

	"github.com/teserakt-io/automation-engine/internal/models"
)

func TestTriggerStateService(t *testing.T) {
	testTriggerStateServiceDatabase(t, sqliteTestDB)

	if os.Getenv("C2AETEST_POSTGRES") == "" {
		t.Skip("C2AETEST_POSTGRES environment is not set")

		return
	}
	testTriggerStateServiceDatabase(t, postgresTestDB)
}

func createTrigger(t *testing.T, db models.Database) models.Trigger {
	rule := &models.Rule{}
	if res := db.Connection().Save(rule); res.Error != nil {
		t.Fatalf("Failed to save dummy rule: %v", res.Error)
	}

	trigger := &models.Trigger{
		RuleID: rule.ID,
	}
	if res := db.Connection().Save(trigger); res.Error != nil {
		t.Fatalf("Failed to save dummy trigger: %v", res.Error)
	}

	return *trigger
}

func testTriggerStateServiceDatabase(t *testing.T, getTestDB func(t *testing.T) (models.Database, func())) {
	ctx := context.Background()

	t.Run("Save properly save triggerState", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewTriggerStateService(db)

		t1 := &models.TriggerState{
			Counter: 5,
		}

		if err := srv.Save(ctx, t1); err == nil {
			t.Error("Expected save with missing trigger foreign key to fail")
		}

		trigger := createTrigger(t, db)

		t2 := &models.TriggerState{
			TriggerID: trigger.ID,
			Counter:   5,
		}

		if err := srv.Save(ctx, t2); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if t2.ID == 0 {
			t.Errorf("Expected triggerState to get affected an ID, got %d", t1.ID)
		}
	})

	t.Run("ByTriggerID returns triggerState", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewTriggerStateService(db)

		trigger := createTrigger(t, db)

		expectedState := models.TriggerState{
			TriggerID: trigger.ID,
			Counter:   5,
		}

		if err := srv.Save(ctx, &expectedState); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		state, err := srv.ByTriggerID(ctx, expectedState.ID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if reflect.DeepEqual(state, expectedState) == false {
			t.Errorf("Expected state to be %#v, got %#v", expectedState, state)
		}
	})

	t.Run("ByTriggerID returns a default state when it doesn't exist and no errors", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewTriggerStateService(db)

		state, err := srv.ByTriggerID(ctx, 1)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if reflect.DeepEqual(state, models.TriggerState{TriggerID: 1}) == false {
			t.Errorf("Expected state to be %#v, got %#v", models.TriggerState{}, state)
		}
	})

	t.Run("Deleting a trigger cascade delete the trigger state", func(t *testing.T) {
		db, closeFunc := getTestDB(t)
		defer closeFunc()

		srv := NewTriggerStateService(db)

		trigger1 := createTrigger(t, db)
		trigger2 := createTrigger(t, db)

		trigger1State := models.TriggerState{
			TriggerID: trigger1.ID,
			Counter:   5,
		}
		trigger2State := models.TriggerState{
			TriggerID: trigger2.ID,
			Counter:   6,
		}

		for _, s := range []models.TriggerState{trigger1State, trigger2State} {
			if err := srv.Save(ctx, &s); err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		}

		db.Connection().Delete(trigger1)

		state, err := srv.ByTriggerID(ctx, trigger1.ID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if state.Counter != 0 {
			t.Errorf("Expected trigger1 state to have been deleted, got %#v", state)
		}

		state, err = srv.ByTriggerID(ctx, trigger2.ID)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if state.Counter != 6 {
			t.Errorf("Expected trigger2 state to not be modified, got %#v", state)
		}
	})
}
