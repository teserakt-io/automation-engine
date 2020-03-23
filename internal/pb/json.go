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
