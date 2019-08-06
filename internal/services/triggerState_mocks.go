// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/teserakt/c2ae/internal/services (interfaces: TriggerStateService)

// Package services is a generated GoMock package.
package services

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "gitlab.com/teserakt/c2ae/internal/models"
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
