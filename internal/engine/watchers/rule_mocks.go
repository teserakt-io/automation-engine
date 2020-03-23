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
//

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/teserakt-io/automation-engine/internal/engine/watchers (interfaces: RuleWatcher)

// Package watchers is a generated GoMock package.
package watchers

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRuleWatcher is a mock of RuleWatcher interface
type MockRuleWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockRuleWatcherMockRecorder
}

// MockRuleWatcherMockRecorder is the mock recorder for MockRuleWatcher
type MockRuleWatcherMockRecorder struct {
	mock *MockRuleWatcher
}

// NewMockRuleWatcher creates a new mock instance
func NewMockRuleWatcher(ctrl *gomock.Controller) *MockRuleWatcher {
	mock := &MockRuleWatcher{ctrl: ctrl}
	mock.recorder = &MockRuleWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRuleWatcher) EXPECT() *MockRuleWatcherMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockRuleWatcher) Start(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", arg0)
}

// Start indicates an expected call of Start
func (mr *MockRuleWatcherMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockRuleWatcher)(nil).Start), arg0)
}
