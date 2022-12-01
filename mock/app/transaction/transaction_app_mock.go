// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/danilotadeu/pismo/app/transaction (interfaces: App)

// Package mockAppTransaction is a generated GoMock package.
package mockAppTransaction

import (
	context "context"
	reflect "reflect"

	transaction "github.com/danilotadeu/pismo/model/transaction"
	gomock "github.com/golang/mock/gomock"
)

// MockApp is a mock of App interface.
type MockApp struct {
	ctrl     *gomock.Controller
	recorder *MockAppMockRecorder
}

// MockAppMockRecorder is the mock recorder for MockApp.
type MockAppMockRecorder struct {
	mock *MockApp
}

// NewMockApp creates a new mock instance.
func NewMockApp(ctrl *gomock.Controller) *MockApp {
	mock := &MockApp{ctrl: ctrl}
	mock.recorder = &MockAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApp) EXPECT() *MockAppMockRecorder {
	return m.recorder
}

// CreateTransaction mocks base method.
func (m *MockApp) CreateTransaction(arg0 context.Context, arg1 transaction.TransactionRequest) (*int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", arg0, arg1)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockAppMockRecorder) CreateTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockApp)(nil).CreateTransaction), arg0, arg1)
}
