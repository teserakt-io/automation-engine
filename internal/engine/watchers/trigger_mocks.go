// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/teserakt/c2ae/internal/engine/watchers (interfaces: TriggerWatcher)

// Package watchers is a generated GoMock package.
package watchers

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockTriggerWatcher is a mock of TriggerWatcher interface
type MockTriggerWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockTriggerWatcherMockRecorder
}

// MockTriggerWatcherMockRecorder is the mock recorder for MockTriggerWatcher
type MockTriggerWatcherMockRecorder struct {
	mock *MockTriggerWatcher
}

// NewMockTriggerWatcher creates a new mock instance
func NewMockTriggerWatcher(ctrl *gomock.Controller) *MockTriggerWatcher {
	mock := &MockTriggerWatcher{ctrl: ctrl}
	mock.recorder = &MockTriggerWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTriggerWatcher) EXPECT() *MockTriggerWatcherMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockTriggerWatcher) Start(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", arg0)
}

// Start indicates an expected call of Start
func (mr *MockTriggerWatcherMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockTriggerWatcher)(nil).Start), arg0)
}

// UpdateLastExecuted mocks base method
func (m *MockTriggerWatcher) UpdateLastExecuted(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLastExecuted", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLastExecuted indicates an expected call of UpdateLastExecuted
func (mr *MockTriggerWatcherMockRecorder) UpdateLastExecuted(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLastExecuted", reflect.TypeOf((*MockTriggerWatcher)(nil).UpdateLastExecuted), arg0)
}