package models

import (
	"time"

	"gitlab.com/teserakt/c2se/internal/pb"
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

// Trigger holds database informations for a rule trigger
type Trigger struct {
	ID          int `gorm:"primary_key"`
	RuleID      int `gorm:"type:int REFERENCES rules(id) ON DELETE CASCADE"`
	TriggerType pb.TriggerType
	Settings    []byte
	State       []byte
}

// Target holds database informations for a rule target
type Target struct {
	ID     int `gorm:"primary_key"`
	RuleID int `gorm:"type:int REFERENCES rules(id) ON DELETE CASCADE"`
	Type   pb.TargetType
	Expr   string
}
