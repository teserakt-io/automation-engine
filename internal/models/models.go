package models

import (
	"time"

	"github.com/teserakt-io/automation-engine/internal/pb"
)

// Rule holds database information of a rule.
type Rule struct {
	ID           int `gorm:"primary_key:true"`
	Description  string
	ActionType   pb.ActionType
	LastExecuted time.Time
	Triggers     []Trigger
	Targets      []Target
}

// Target holds database informations for a rule target
type Target struct {
	ID     int `gorm:"primary_key"`
	RuleID int `gorm:"type:int REFERENCES rules(id) ON DELETE CASCADE; index;"`
	Type   pb.TargetType
	Expr   string
}

// Trigger holds database informations for a rule trigger
type Trigger struct {
	ID          int `gorm:"primary_key"`
	RuleID      int `gorm:"type:int REFERENCES rules(id) ON DELETE CASCADE; index;"`
	TriggerType pb.TriggerType
	Settings    []byte
}

// TriggerState holds data to be persisted by a trigger watcher
type TriggerState struct {
	ID        int `gorm:"primary_key"`
	TriggerID int `gorm:"type:int REFERENCES triggers(id) ON DELETE CASCADE; unique_index; NOT NULL;"`
	Counter   int
}

// FilterNonExistingTriggers will returns a slice of Triggers
// from `old` which does not exists in `new`
func FilterNonExistingTriggers(old []Trigger, new []Trigger) []Trigger {
	filtered := []Trigger{}
	for _, oldTrigger := range old {
		if !containsTrigger(oldTrigger, new) {
			filtered = append(filtered, oldTrigger)
		}
	}

	return filtered
}

func containsTrigger(needle Trigger, haystack []Trigger) bool {
	for _, trigger := range haystack {
		if trigger.ID == needle.ID {
			return true
		}
	}

	return false
}

// FilterNonExistingTargets will returns a slice of Targets
// from `old` which does not exists in `new`
func FilterNonExistingTargets(old []Target, new []Target) []Target {
	filtered := []Target{}
	for _, oldTarget := range old {
		if !containsTarget(oldTarget, new) {
			filtered = append(filtered, oldTarget)
		}
	}

	return filtered
}

func containsTarget(needle Target, haystack []Target) bool {
	for _, target := range haystack {
		if target.ID == needle.ID {
			return true
		}
	}

	return false
}
