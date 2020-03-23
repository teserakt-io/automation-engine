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

package pb

import (
	"encoding/json"
	"errors"
	fmt "fmt"

	"github.com/gorhill/cronexpr"
	c2pb "github.com/teserakt-io/c2/pkg/pb"
)

// EventType  is a custom type for all available Trigger Events
type EventType string

var (
	// EventTypeClientSubscribed describes the event emitted when a client subscribe
	EventTypeClientSubscribed EventType = EventType(c2pb.EventType_CLIENT_SUBSCRIBED.String())
	// EventTypeClientUnsubscribed describes the event emitted when a client unsubscribe
	EventTypeClientUnsubscribed EventType = EventType(c2pb.EventType_CLIENT_UNSUBSCRIBED.String())
)

// TriggerSettings defines a generic trigger settings structure
type TriggerSettings interface {
	Validate() error
	Encode() ([]byte, error)
	Decode([]byte) error
}

// TriggerSettingsTimeInterval holds settings for pb.TriggerType_TIME_INTERVAL trigger types
type TriggerSettingsTimeInterval struct {
	Expr string `json:"expr,omitempty"`
}

// TriggerSettingsEvent holds settings for event driven trigger types
type TriggerSettingsEvent struct {
	EventType     EventType `json:"eventType,omitempty"`
	MaxOccurrence int       `json:"maxOccurrence,omitempty"`
}

var _ TriggerSettings = &TriggerSettingsTimeInterval{}
var _ TriggerSettings = &TriggerSettingsEvent{}

// Decode will attempt to turn []byte settings into matching struct given the trigger type
func Decode(t TriggerType, settings []byte) (TriggerSettings, error) {
	var triggerSettings TriggerSettings

	switch t {
	case TriggerType_TIME_INTERVAL:
		triggerSettings = &TriggerSettingsTimeInterval{}
	case TriggerType_EVENT:
		triggerSettings = &TriggerSettingsEvent{}
	default:
		return nil, fmt.Errorf("trigger type %s is not supported", t)
	}

	triggerSettings.Decode(settings)

	return triggerSettings, nil
}

// Validate implements TriggerSettings and returns an error when the settings are invalid
func (t *TriggerSettingsTimeInterval) Validate() error {
	if len(t.Expr) == 0 {
		return errors.New("expr field is required and must be a valid cron expression")
	}

	_, err := cronexpr.Parse(t.Expr)
	if err != nil {
		return fmt.Errorf("failed to parse cron expression from Expr field: %s", err)
	}

	return nil
}

// Encode json encode settings to []byte
func (t *TriggerSettingsTimeInterval) Encode() ([]byte, error) {
	return jsonEncode(t)
}

// Decode json decode bytes to settings
func (t *TriggerSettingsTimeInterval) Decode(b []byte) error {
	return jsonDecode(t, b)
}

// Validate implements TriggerSettings and returns an error when the settings are invalid
func (t *TriggerSettingsEvent) Validate() error {
	if len(t.EventType) <= 0 {
		return errors.New("EventType is required")
	}

	switch t.EventType {
	case EventTypeClientSubscribed, EventTypeClientUnsubscribed:
	default:
		return fmt.Errorf("EventType must be one of %v", []EventType{EventTypeClientSubscribed, EventTypeClientUnsubscribed})
	}

	if t.MaxOccurrence <= 0 {
		return errors.New("MaxOccurrence must be greater than 0")
	}

	return nil
}

// Encode json encode settings to []byte
func (t *TriggerSettingsEvent) Encode() ([]byte, error) {
	return jsonEncode(t)
}

// Decode json decode bytes to settings
func (t *TriggerSettingsEvent) Decode(b []byte) error {
	return jsonDecode(t, b)
}

func jsonEncode(t interface{}) ([]byte, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func jsonDecode(t interface{}, b []byte) error {
	if err := json.Unmarshal(b, t); err != nil {
		return err
	}

	return nil
}
