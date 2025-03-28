// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/handlers/apishorten.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/rycln/shorturl/internal/app/storage"
)

// MockapiStorager is a mock of apiStorager interface.
type MockapiStorager struct {
	ctrl     *gomock.Controller
	recorder *MockapiStoragerMockRecorder
}

// MockapiStoragerMockRecorder is the mock recorder for MockapiStorager.
type MockapiStoragerMockRecorder struct {
	mock *MockapiStorager
}

// NewMockapiStorager creates a new mock instance.
func NewMockapiStorager(ctrl *gomock.Controller) *MockapiStorager {
	mock := &MockapiStorager{ctrl: ctrl}
	mock.recorder = &MockapiStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockapiStorager) EXPECT() *MockapiStoragerMockRecorder {
	return m.recorder
}

// AddURL mocks base method.
func (m *MockapiStorager) AddURL(arg0 context.Context, arg1 storage.ShortenedURL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddURL", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddURL indicates an expected call of AddURL.
func (mr *MockapiStoragerMockRecorder) AddURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddURL", reflect.TypeOf((*MockapiStorager)(nil).AddURL), arg0, arg1)
}

// GetShortURL mocks base method.
func (m *MockapiStorager) GetShortURL(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURL", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURL indicates an expected call of GetShortURL.
func (mr *MockapiStoragerMockRecorder) GetShortURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURL", reflect.TypeOf((*MockapiStorager)(nil).GetShortURL), arg0, arg1)
}

// MockapiConfiger is a mock of apiConfiger interface.
type MockapiConfiger struct {
	ctrl     *gomock.Controller
	recorder *MockapiConfigerMockRecorder
}

// MockapiConfigerMockRecorder is the mock recorder for MockapiConfiger.
type MockapiConfigerMockRecorder struct {
	mock *MockapiConfiger
}

// NewMockapiConfiger creates a new mock instance.
func NewMockapiConfiger(ctrl *gomock.Controller) *MockapiConfiger {
	mock := &MockapiConfiger{ctrl: ctrl}
	mock.recorder = &MockapiConfigerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockapiConfiger) EXPECT() *MockapiConfigerMockRecorder {
	return m.recorder
}

// GetBaseAddr mocks base method.
func (m *MockapiConfiger) GetBaseAddr() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBaseAddr")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetBaseAddr indicates an expected call of GetBaseAddr.
func (mr *MockapiConfigerMockRecorder) GetBaseAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBaseAddr", reflect.TypeOf((*MockapiConfiger)(nil).GetBaseAddr))
}

// GetKey mocks base method.
func (m *MockapiConfiger) GetKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetKey indicates an expected call of GetKey.
func (mr *MockapiConfigerMockRecorder) GetKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockapiConfiger)(nil).GetKey))
}
