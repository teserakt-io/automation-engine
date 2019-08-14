// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/teserakt/c2ae/internal/services (interfaces: C2,C2EventStreamClient)

// Package services is a generated GoMock package.
package services

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	pb "gitlab.com/teserakt/c2/pkg/pb"
	metadata "google.golang.org/grpc/metadata"
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
func (m *MockC2) NewClientKey(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewClientKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewClientKey indicates an expected call of NewClientKey
func (mr *MockC2MockRecorder) NewClientKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewClientKey", reflect.TypeOf((*MockC2)(nil).NewClientKey), arg0, arg1)
}

// NewTopicKey mocks base method
func (m *MockC2) NewTopicKey(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTopicKey", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewTopicKey indicates an expected call of NewTopicKey
func (mr *MockC2MockRecorder) NewTopicKey(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTopicKey", reflect.TypeOf((*MockC2)(nil).NewTopicKey), arg0, arg1)
}

// SubscribeToEventStream mocks base method
func (m *MockC2) SubscribeToEventStream(arg0 context.Context) (C2EventStreamClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToEventStream", arg0)
	ret0, _ := ret[0].(C2EventStreamClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeToEventStream indicates an expected call of SubscribeToEventStream
func (mr *MockC2MockRecorder) SubscribeToEventStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToEventStream", reflect.TypeOf((*MockC2)(nil).SubscribeToEventStream), arg0)
}

// MockC2EventStreamClient is a mock of C2EventStreamClient interface
type MockC2EventStreamClient struct {
	ctrl     *gomock.Controller
	recorder *MockC2EventStreamClientMockRecorder
}

// MockC2EventStreamClientMockRecorder is the mock recorder for MockC2EventStreamClient
type MockC2EventStreamClientMockRecorder struct {
	mock *MockC2EventStreamClient
}

// NewMockC2EventStreamClient creates a new mock instance
func NewMockC2EventStreamClient(ctrl *gomock.Controller) *MockC2EventStreamClient {
	mock := &MockC2EventStreamClient{ctrl: ctrl}
	mock.recorder = &MockC2EventStreamClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockC2EventStreamClient) EXPECT() *MockC2EventStreamClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method
func (m *MockC2EventStreamClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockC2EventStreamClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockC2EventStreamClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockC2EventStreamClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockC2EventStreamClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockC2EventStreamClient)(nil).Context))
}

// Header mocks base method
func (m *MockC2EventStreamClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockC2EventStreamClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockC2EventStreamClient)(nil).Header))
}

// Recv mocks base method
func (m *MockC2EventStreamClient) Recv() (*pb.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*pb.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockC2EventStreamClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockC2EventStreamClient)(nil).Recv))
}

// RecvMsg mocks base method
func (m *MockC2EventStreamClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockC2EventStreamClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockC2EventStreamClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method
func (m *MockC2EventStreamClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockC2EventStreamClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockC2EventStreamClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method
func (m *MockC2EventStreamClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockC2EventStreamClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockC2EventStreamClient)(nil).Trailer))
}
