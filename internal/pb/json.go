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
	triggerSettings, err := Decode(t.Type, t.Settings)
	if err != nil {
		return nil, fmt.Errorf("json marshalling failed: %v", err)
	}

	jsonTrigger := jsonTrigger{
		ID:       t.Id,
		Type:     t.Type,
		Settings: triggerSettings,
	}

	return json.Marshal(jsonTrigger)
}
