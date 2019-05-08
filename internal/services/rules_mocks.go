// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/teserakt/c2se/internal/services (interfaces: RuleService)

// Package services is a generated GoMock package.
package services

import (
	gomock "github.com/golang/mock/gomock"
	models "gitlab.com/teserakt/c2se/internal/models"
	reflect "reflect"
)

// MockRuleService is a mock of RuleService interface
type MockRuleService struct {
	ctrl     *gomock.Controller
	recorder *MockRuleServiceMockRecorder
}

// MockRuleServiceMockRecorder is the mock recorder for MockRuleService
type MockRuleServiceMockRecorder struct {
	mock *MockRuleService
}

// NewMockRuleService creates a new mock instance
func NewMockRuleService(ctrl *gomock.Controller) *MockRuleService {
	mock := &MockRuleService{ctrl: ctrl}
	mock.recorder = &MockRuleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRuleService) EXPECT() *MockRuleServiceMockRecorder {
	return m.recorder
}

// All mocks base method
func (m *MockRuleService) All() ([]models.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All")
	ret0, _ := ret[0].([]models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All
func (mr *MockRuleServiceMockRecorder) All() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockRuleService)(nil).All))
}

// ByID mocks base method
func (m *MockRuleService) ByID(arg0 int) (models.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ByID", arg0)
	ret0, _ := ret[0].(models.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ByID indicates an expected call of ByID
func (mr *MockRuleServiceMockRecorder) ByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ByID", reflect.TypeOf((*MockRuleService)(nil).ByID), arg0)
}

// Delete mocks base method
func (m *MockRuleService) Delete(arg0 models.Rule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRuleServiceMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRuleService)(nil).Delete), arg0)
}

// Save mocks base method
func (m *MockRuleService) Save(arg0 *models.Rule) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockRuleServiceMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRuleService)(nil).Save), arg0)
}

// TargetByID mocks base method
func (m *MockRuleService) TargetByID(arg0 int) (models.Target, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TargetByID", arg0)
	ret0, _ := ret[0].(models.Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TargetByID indicates an expected call of TargetByID
func (mr *MockRuleServiceMockRecorder) TargetByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TargetByID", reflect.TypeOf((*MockRuleService)(nil).TargetByID), arg0)
}

// TriggerByID mocks base method
func (m *MockRuleService) TriggerByID(arg0 int) (models.Trigger, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TriggerByID", arg0)
	ret0, _ := ret[0].(models.Trigger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TriggerByID indicates an expected call of TriggerByID
func (mr *MockRuleServiceMockRecorder) TriggerByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TriggerByID", reflect.TypeOf((*MockRuleService)(nil).TriggerByID), arg0)
}