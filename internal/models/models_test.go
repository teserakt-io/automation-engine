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

package models

import (
	"reflect"
	"testing"
)

func TestFilterNonExistingTriggers(t *testing.T) {
	testDataSet := []struct {
		Old      []Trigger
		New      []Trigger
		Expected []Trigger
	}{
		{
			Old:      []Trigger{Trigger{ID: 1}, Trigger{ID: 2}, Trigger{ID: 3}},
			New:      []Trigger{Trigger{ID: 1}, Trigger{ID: 3}, Trigger{ID: 4}},
			Expected: []Trigger{Trigger{ID: 2}},
		},
		{
			Old:      []Trigger{},
			New:      []Trigger{Trigger{ID: 1}, Trigger{ID: 2}},
			Expected: []Trigger{},
		},
		{
			Old:      []Trigger{Trigger{ID: 1}, Trigger{ID: 2}},
			New:      []Trigger{},
			Expected: []Trigger{Trigger{ID: 1}, Trigger{ID: 2}},
		},
	}

	for _, testData := range testDataSet {
		filtered := FilterNonExistingTriggers(testData.Old, testData.New)
		if reflect.DeepEqual(filtered, testData.Expected) == false {
			t.Errorf("Expected filtered Triggers to be %#v, got %#v", testData.Expected, filtered)
		}
	}
}

func TestFilterNonExistingTargets(t *testing.T) {
	testDataSet := []struct {
		Old      []Target
		New      []Target
		Expected []Target
	}{
		{
			Old:      []Target{Target{ID: 1}, Target{ID: 2}, Target{ID: 3}},
			New:      []Target{Target{ID: 1}, Target{ID: 3}, Target{ID: 4}},
			Expected: []Target{Target{ID: 2}},
		},
		{
			Old:      []Target{},
			New:      []Target{Target{ID: 1}, Target{ID: 2}},
			Expected: []Target{},
		},
		{
			Old:      []Target{Target{ID: 1}, Target{ID: 2}},
			New:      []Target{},
			Expected: []Target{Target{ID: 1}, Target{ID: 2}},
		},
	}

	for _, testData := range testDataSet {
		filtered := FilterNonExistingTargets(testData.Old, testData.New)
		if reflect.DeepEqual(filtered, testData.Expected) == false {
			t.Errorf("Expected filtered Targets to be %#v, got %#v", testData.Expected, filtered)
		}
	}
}
