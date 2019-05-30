// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/teserakt/c2ae/internal/services (interfaces: C2,C2Requester)

// Package services is a generated GoMock package.
package services

import (
	gomock "github.com/golang/mock/gomock"
	e4common "gitlab.com/teserakt/e4common"
	reflect "reflect"
)

// MockC2 is a mock of C2 interface
type MockC2 struct {
	ctrl     *gomock.Controller
	recorder *MockC2MockRecorder
}

// MockC2MockRecorder is the mock recorder for MockC2
type MockC2MockRecorder struct {
	mock *MockC2
}

// NewMockC2 creates a new mock instance
func NewMockC2(ctrl *gomock.Controller) *MockC2 {
	mock := &MockC2{ctrl: ctrl}
	mock.recorder = &MockC2MockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockC2) EXPECT() *MockC2MockRecorder {
	return m.recorder
}

// NewClientKey mocks base method
func (m *MockC2) NewClientKey(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClientKey", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewClientKey indicates an expected call of NewClientKey
func (mr *MockC2MockRecorder) NewClientKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClientKey", reflect.TypeOf((*MockC2)(nil).NewClientKey), arg0)
}

// NewTopicKey mocks base method
func (m *MockC2) NewTopicKey(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTopicKey", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewTopicKey indicates an expected call of NewTopicKey
func (mr *MockC2MockRecorder) NewTopicKey(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTopicKey", reflect.TypeOf((*MockC2)(nil).NewTopicKey), arg0)
}

// MockC2Requester is a mock of C2Requester interface
type MockC2Requester struct {
	ctrl     *gomock.Controller
	recorder *MockC2RequesterMockRecorder
}

// MockC2RequesterMockRecorder is the mock recorder for MockC2Requester
type MockC2RequesterMockRecorder struct {
	mock *MockC2Requester
}

// NewMockC2Requester creates a new mock instance
func NewMockC2Requester(ctrl *gomock.Controller) *MockC2Requester {
	mock := &MockC2Requester{ctrl: ctrl}
	mock.recorder = &MockC2RequesterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockC2Requester) EXPECT() *MockC2RequesterMockRecorder {
	return m.recorder
}

// C2Request mocks base method
func (m *MockC2Requester) C2Request(arg0 e4common.C2Request) (e4common.C2Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "C2Request", arg0)
	ret0, _ := ret[0].(e4common.C2Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// C2Request indicates an expected call of C2Request
func (mr *MockC2RequesterMockRecorder) C2Request(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "C2Request", reflect.TypeOf((*MockC2Requester)(nil).C2Request), arg0)
}
