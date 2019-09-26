package pb

import (
	"strings"
	"testing"
)

func TestJson(t *testing.T) {
	t.Run("Trigger MarshalJson properly marshall TriggerSettingsTimeInterval", func(t *testing.T) {
		expectedExpr := "expectedExpr"
		triggerSettings := &TriggerSettingsTimeInterval{Expr: expectedExpr}
		encodedSettings, err := triggerSettings.Encode()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		trigger := &Trigger{
			Type:     TriggerType_TIME_INTERVAL,
			Settings: encodedSettings,
		}

		json, err := trigger.MarshalJSON()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if strings.Contains(string(json), expectedExpr) == false {
			t.Errorf("Expected json settings to contains '%s', but got '%s'", expectedExpr, string(json))
		}
	})

	t.Run("Trigger MarshalJson properly marshall TriggerSettingsEvent", func(t *testing.T) {
		triggerSettings := &TriggerSettingsEvent{MaxOccurrence: 5, EventType: EventTypeClientSubscribed}
		encodedSettings, err := triggerSettings.Encode()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		trigger := &Trigger{
			Type:     TriggerType_EVENT,
			Settings: encodedSettings,
		}

		json, err := trigger.MarshalJSON()
		if err != nil {
			t.Errorf("Expected err to be nil, got %s", err)
		}

		if strings.Contains(string(json), string(EventTypeClientSubscribed)) == false {
			t.Errorf("Expected json settings to contains '%s', but got '%s'", EventTypeClientSubscribed, string(json))
		}
	})

	t.Run("Trigger MarshalJson returns an error on unsupported trigger type", func(t *testing.T) {
		trigger := &Trigger{}

		_, err := trigger.MarshalJSON()
		if err == nil {
			t.Error("Expected err to be not nil")
		}
	})
}
