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
// Source: github.com/teserakt-io/automation-engine/internal/services (interfaces: TriggerStateService)

// Package services is a generated GoMock package.
package services

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "github.com/teserakt-io/automation-engine/internal/models"
	reflect "reflect"
)

// MockTriggerStateService is a mock of TriggerStateService interface
type MockTriggerStateService struct {
	ctrl     *gomock.Controller
	recorder *MockTriggerStateServiceMockRecorder
}

// MockTriggerStateServiceMockRecorder is the mock recorder for MockTriggerStateService
type MockTriggerStateServiceMockRecorder struct {
	mock *MockTriggerStateService
}

// NewMockTriggerStateService creates a new mock instance
func NewMockTriggerStateService(ctrl *gomock.Controller) *MockTriggerStateService {
	mock := &MockTriggerStateService{ctrl: ctrl}
	mock.recorder = &MockTriggerStateServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTriggerStateService) EXPECT() *MockTriggerStateServiceMockRecorder {
	return m.recorder
}

// ByTriggerID mocks base method
func (m *MockTriggerStateService) ByTriggerID(arg0 context.Context, arg1 int) (models.TriggerState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByTriggerID", arg0, arg1)
	ret0, _ := ret[0].(models.TriggerState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByTriggerID indicates an expected call of ByTriggerID
func (mr *MockTriggerStateServiceMockRecorder) ByTriggerID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByTriggerID", reflect.TypeOf((*MockTriggerStateService)(nil).ByTriggerID), arg0, arg1)
}

// Save mocks base method
func (m *MockTriggerStateService) Save(arg0 context.Context, arg1 *models.TriggerState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockTriggerStateServiceMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTriggerStateService)(nil).Save), arg0, arg1)
}
