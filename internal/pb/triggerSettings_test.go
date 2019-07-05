package pb

import (
	"reflect"
	"testing"
)

func TestTriggerSettings(t *testing.T) {
	t.Run("TriggerSettings encode / decode properly works", func(t *testing.T) {

		testData := map[TriggerSettings]TriggerSettings{
			&TriggerSettingsTimeInterval{}: &TriggerSettingsTimeInterval{Expr: "something"},
			&TriggerSettingsEvent{}:        &TriggerSettingsEvent{EventType: "something", MaxOccurence: 5},
		}

		for settings, expectedSettings := range testData {
			encoded, err := expectedSettings.Encode()
			if err != nil {
				t.Errorf("Expected err to be nil, got %s", err)
			}

			settings.Decode(encoded)

			if reflect.DeepEqual(settings, expectedSettings) == false {
				t.Errorf("Expected settings to be %#v, got %#v", expectedSettings, settings)
			}
		}
	})
}

func TestTriggerSettingsTimeInterval(t *testing.T) {
	t.Run("Validate properly checks settings", func(t *testing.T) {
		testData := map[*TriggerSettingsTimeInterval]bool{
			&TriggerSettingsTimeInterval{Expr: ""}:                     false,
			&TriggerSettingsTimeInterval{Expr: "*****"}:                false,
			&TriggerSettingsTimeInterval{Expr: "* * * * *"}:            true,
			&TriggerSettingsTimeInterval{Expr: "0/5 * * * *"}:          true,
			&TriggerSettingsTimeInterval{Expr: "0 0 12 ? * WED,SAT *"}: true,
			&TriggerSettingsTimeInterval{Expr: "0 0 2 ? 1 MON#1 *"}:    true,
		}

		for settings, valid := range testData {
			err := settings.Validate()

			if valid && err != nil {
				t.Errorf("Expected err to be nil, got %s with settings: %#v", err, settings)
			} else if !valid && err == nil {
				t.Errorf("Expected err to be not nil with settings: %#v", settings)
			}
		}
	})
}

// TODO: enable back this test when we'll have some real usage of TriggerSettingsEvent

// func TestTriggerSettingsEvent(t *testing.T) {
// 	t.Run("Validate properly checks settings", func(t *testing.T) {
// 		testData := map[*TriggerSettingsEvent]bool{
// 			&TriggerSettingsEvent{EventType: EventType(""), MaxOccurence: 0}:                 false,
// 			&TriggerSettingsEvent{EventType: EventType("something"), MaxOccurence: 0}:        false,
// 			&TriggerSettingsEvent{EventType: EventType("something"), MaxOccurence: -1}:       false,
// 			&TriggerSettingsEvent{EventType: EventType("something"), MaxOccurence: 1}:        true,
// 			&TriggerSettingsEvent{EventType: EventTypeClientSubscribed, MaxOccurence: 5}:     true,
// 			&TriggerSettingsEvent{EventType: EventTypeClientUnsubscribed, MaxOccurence: 100}: true,
// 		}

// 		for settings, valid := range testData {
// 			err := settings.Validate()

// 			if valid && err != nil {
// 				t.Errorf("Expected err to be nil, got %s with settings: %#v", err, settings)
// 			} else if !valid && err == nil {
// 				t.Errorf("Expected err to be not nil with settings: %#v", settings)
// 			}
// 		}
// 	})
// }
