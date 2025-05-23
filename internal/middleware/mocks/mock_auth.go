// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/rycln/shorturl/internal/models"
)

// MockauthServicer is a mock of authServicer interface.
type MockauthServicer struct {
	ctrl     *gomock.Controller
	recorder *MockauthServicerMockRecorder
}

// MockauthServicerMockRecorder is the mock recorder for MockauthServicer.
type MockauthServicerMockRecorder struct {
	mock *MockauthServicer
}

// NewMockauthServicer creates a new mock instance.
func NewMockauthServicer(ctrl *gomock.Controller) *MockauthServicer {
	mock := &MockauthServicer{ctrl: ctrl}
	mock.recorder = &MockauthServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockauthServicer) EXPECT() *MockauthServicerMockRecorder {
	return m.recorder
}

// NewJWTString mocks base method.
func (m *MockauthServicer) NewJWTString(arg0 models.UserID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewJWTString", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewJWTString indicates an expected call of NewJWTString.
func (mr *MockauthServicerMockRecorder) NewJWTString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewJWTString", reflect.TypeOf((*MockauthServicer)(nil).NewJWTString), arg0)
}

// ParseIDFromAuthHeader mocks base method.
func (m *MockauthServicer) ParseIDFromAuthHeader(arg0 string) (models.UserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseIDFromAuthHeader", arg0)
	ret0, _ := ret[0].(models.UserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseIDFromAuthHeader indicates an expected call of ParseIDFromAuthHeader.
func (mr *MockauthServicerMockRecorder) ParseIDFromAuthHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseIDFromAuthHeader", reflect.TypeOf((*MockauthServicer)(nil).ParseIDFromAuthHeader), arg0)
}
