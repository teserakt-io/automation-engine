package models

import (
	"time"

	"gitlab.com/teserakt/c2ae/internal/pb"
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
	ID     int   `gorm:"primary_key"`
	RuleID int   `gorm:"type:int REFERENCES rules(id) ON DELETE CASCADE"`
	Rule   *Rule `gorm:"-"`
	Type   pb.TargetType
	Expr   string
}

// Trigger holds database informations for a rule trigger
type Trigger struct {
	ID          int   `gorm:"primary_key"`
	RuleID      int   `gorm:"type:int REFERENCES rules(id) ON DELETE CASCADE"`
	Rule        *Rule `gorm:"-"`
	TriggerType pb.TriggerType
	Settings    []byte
	State       []byte
}

// AfterFind update triggers and targets with current rule
func (r *Rule) AfterFind() error {
	for i := range r.Triggers {
		r.Triggers[i].Rule = r
	}

	for i := range r.Targets {
		r.Targets[i].Rule = r
	}

	return nil
}
