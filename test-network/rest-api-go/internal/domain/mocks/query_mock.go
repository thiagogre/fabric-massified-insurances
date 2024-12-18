// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/query.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQueryInterface is a mock of QueryInterface interface.
type MockQueryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockQueryInterfaceMockRecorder
}

// MockQueryInterfaceMockRecorder is the mock recorder for MockQueryInterface.
type MockQueryInterfaceMockRecorder struct {
	mock *MockQueryInterface
}

// NewMockQueryInterface creates a new mock instance.
func NewMockQueryInterface(ctrl *gomock.Controller) *MockQueryInterface {
	mock := &MockQueryInterface{ctrl: ctrl}
	mock.recorder = &MockQueryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueryInterface) EXPECT() *MockQueryInterfaceMockRecorder {
	return m.recorder
}

// ExecuteQuery mocks base method.
func (m *MockQueryInterface) ExecuteQuery(channelID, chainCodeName, function string, args []string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteQuery", channelID, chainCodeName, function, args)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteQuery indicates an expected call of ExecuteQuery.
func (mr *MockQueryInterfaceMockRecorder) ExecuteQuery(channelID, chainCodeName, function, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteQuery", reflect.TypeOf((*MockQueryInterface)(nil).ExecuteQuery), channelID, chainCodeName, function, args)
}
