package pb

import (
	"encoding/json"
	fmt "fmt"
)

type jsonTrigger struct {
	ID       int32           `json:"id,omitempty"`
	Type     TriggerType     `json:"type,omitempty"`
	Settings TriggerSettings `json:"settings,omitempty"`
}

// MarshalJSON is a custom json marshalling for Trigger type
// allowing to decode binary Settings and State to a proper json representation
func (t *Trigger) MarshalJSON() ([]byte, error) {
	var triggerSettings TriggerSettings

	switch t.Type {
	case TriggerType_TIME_INTERVAL:
		triggerSettings = &TriggerSettingsTimeInterval{}
	case TriggerType_CLIENT_UNSUBSCRIBED, TriggerType_CLIENT_SUBSCRIBED:
		triggerSettings = &TriggerSettingsEvent{}
	default:
		return nil, fmt.Errorf("json marshalling trigger type %s is not supported", t.Type)
	}

	triggerSettings.Decode(t.Settings)

	jsonTrigger := jsonTrigger{
		ID:       t.Id,
		Type:     t.Type,
		Settings: triggerSettings,
	}

	return json.Marshal(jsonTrigger)
}
