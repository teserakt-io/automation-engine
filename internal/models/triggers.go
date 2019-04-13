package models

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	"gitlab.com/teserakt/c2se/internal/pb"

	"github.com/gorhill/cronexpr"
)

// Trigger holds database informations for a rule trigger
type Trigger struct {
	ID          int `gorm:"primary_key"`
	RuleID      int `gorm:"type:int REFERENCES rules(id) ON DELETE CASCADE"`
	TriggerType pb.TriggerType
	Settings    []byte
	State       []byte
}

// EventType  is a custom type for all available Trigger Events
type EventType string

var (
	// EventTypeClientSubscribed describes the event emitted when a client subscribe
	EventTypeClientSubscribed EventType = "client_subscribed"
	// EventTypeClientUnsubscribed describes the event emitted when a client unsubscribe
	EventTypeClientUnsubscribed EventType = "client_unsubscribed"
)

// TriggerSettings defines a generic trigger settings structure
type TriggerSettings interface {
	Validate() error
	Encode() ([]byte, error)
	Decode([]byte) error
}

// TriggerSettingsTimeInterval holds settings for pb.TriggerType_TIME_INTERVAL trigger types
type TriggerSettingsTimeInterval struct {
	Expr string
}

// TriggerSettingsEvent holds settings for event driven trigger types
type TriggerSettingsEvent struct {
	EventType    EventType
	MaxOccurence int
}

var _ TriggerSettings = TriggerSettingsTimeInterval{}
var _ TriggerSettings = TriggerSettingsEvent{}

// Validate implements TriggerSettings and returns an error when the settings are invalid
func (t TriggerSettingsTimeInterval) Validate() error {
	if len(t.Expr) == 0 {
		return errors.New("Expr field is required and must be a valid cron expression")
	}

	_, err := cronexpr.Parse(t.Expr)
	if err != nil {
		return fmt.Errorf("failed to parse cron expression from Expr field: %s", err)
	}

	return nil
}

// Encode gob encode settings to []byte
func (t TriggerSettingsTimeInterval) Encode() ([]byte, error) {
	return gobEncode(&t)
}

// Decode gob decode bytes to settings
func (t TriggerSettingsTimeInterval) Decode(b []byte) error {
	return gobDecode(&t, b)
}

func gobEncode(t interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buf)

	if err := encoder.Encode(t); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func gobDecode(t interface{}, b []byte) error {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)

	if err := decoder.Decode(&t); err != nil {
		return err
	}

	return nil
}

// Validate implements TriggerSettings and returns an error when the settings are invalid
func (t TriggerSettingsEvent) Validate() error {
	if len(t.EventType) == 0 {
		return errors.New("EventType is required")
	}

	return nil
}

// Encode gob encode settings to []byte
func (t TriggerSettingsEvent) Encode() ([]byte, error) {
	return gobEncode(&t)
}

// Decode gob decode bytes to settings
func (t TriggerSettingsEvent) Decode(b []byte) error {
	return gobDecode(&t, b)
}
