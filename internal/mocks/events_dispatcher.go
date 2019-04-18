// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/teserakt/c2se/internal/events (interfaces: Dispatcher)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	events "gitlab.com/teserakt/c2se/internal/events"
)

// MockDispatcher is a mock of Dispatcher interface
type MockDispatcher struct {
	ctrl     *gomock.Controller
	recorder *MockDispatcherMockRecorder
}

// MockDispatcherMockRecorder is the mock recorder for MockDispatcher
type MockDispatcherMockRecorder struct {
	mock *MockDispatcher
}

// NewMockDispatcher creates a new mock instance
func NewMockDispatcher(ctrl *gomock.Controller) *MockDispatcher {
	mock := &MockDispatcher{ctrl: ctrl}
	mock.recorder = &MockDispatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDispatcher) EXPECT() *MockDispatcherMockRecorder {
	return m.recorder
}

// ClearListeners mocks base method
func (m *MockDispatcher) ClearListeners() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearListeners")
}

// ClearListeners indicates an expected call of ClearListeners
func (mr *MockDispatcherMockRecorder) ClearListeners() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearListeners", reflect.TypeOf((*MockDispatcher)(nil).ClearListeners))
}

// Dispatch mocks base method
func (m *MockDispatcher) Dispatch(arg0 events.Type, arg1, arg2 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispatch", arg0, arg1, arg2)
}

// Dispatch indicates an expected call of Dispatch
func (mr *MockDispatcherMockRecorder) Dispatch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockDispatcher)(nil).Dispatch), arg0, arg1, arg2)
}

// Register mocks base method
func (m *MockDispatcher) Register(arg0 events.Type, arg1 events.Listener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", arg0, arg1)
}

// Register indicates an expected call of Register
func (mr *MockDispatcherMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockDispatcher)(nil).Register), arg0, arg1)
}

// Start mocks base method
func (m *MockDispatcher) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start
func (mr *MockDispatcherMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockDispatcher)(nil).Start))
}
