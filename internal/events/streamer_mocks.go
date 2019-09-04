// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/teserakt/c2ae/internal/events (interfaces: Streamer)

// Package events is a generated GoMock package.
package events

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStreamer is a mock of Streamer interface
type MockStreamer struct {
	ctrl     *gomock.Controller
	recorder *MockStreamerMockRecorder
}

// MockStreamerMockRecorder is the mock recorder for MockStreamer
type MockStreamerMockRecorder struct {
	mock *MockStreamer
}

// NewMockStreamer creates a new mock instance
func NewMockStreamer(ctrl *gomock.Controller) *MockStreamer {
	mock := &MockStreamer{ctrl: ctrl}
	mock.recorder = &MockStreamerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStreamer) EXPECT() *MockStreamerMockRecorder {
	return m.recorder
}

// AddListener mocks base method
func (m *MockStreamer) AddListener(arg0 StreamListener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddListener", arg0)
}

// AddListener indicates an expected call of AddListener
func (mr *MockStreamerMockRecorder) AddListener(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddListener", reflect.TypeOf((*MockStreamer)(nil).AddListener), arg0)
}

// Listeners mocks base method
func (m *MockStreamer) Listeners() []StreamListener {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Listeners")
	ret0, _ := ret[0].([]StreamListener)
	return ret0
}

// Listeners indicates an expected call of Listeners
func (mr *MockStreamerMockRecorder) Listeners() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Listeners", reflect.TypeOf((*MockStreamer)(nil).Listeners))
}

// RemoveListener mocks base method
func (m *MockStreamer) RemoveListener(arg0 StreamListener) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveListener", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveListener indicates an expected call of RemoveListener
func (mr *MockStreamerMockRecorder) RemoveListener(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveListener", reflect.TypeOf((*MockStreamer)(nil).RemoveListener), arg0)
}

// StartStream mocks base method
func (m *MockStreamer) StartStream(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartStream", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartStream indicates an expected call of StartStream
func (mr *MockStreamerMockRecorder) StartStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartStream", reflect.TypeOf((*MockStreamer)(nil).StartStream), arg0)
}