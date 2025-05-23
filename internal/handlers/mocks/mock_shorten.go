// Code generated by MockGen. DO NOT EDIT.
// Source: shorten.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/rycln/shorturl/internal/models"
)

// MockshortenServicer is a mock of shortenServicer interface.
type MockshortenServicer struct {
	ctrl     *gomock.Controller
	recorder *MockshortenServicerMockRecorder
}

// MockshortenServicerMockRecorder is the mock recorder for MockshortenServicer.
type MockshortenServicerMockRecorder struct {
	mock *MockshortenServicer
}

// NewMockshortenServicer creates a new mock instance.
func NewMockshortenServicer(ctrl *gomock.Controller) *MockshortenServicer {
	mock := &MockshortenServicer{ctrl: ctrl}
	mock.recorder = &MockshortenServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockshortenServicer) EXPECT() *MockshortenServicerMockRecorder {
	return m.recorder
}

// ShortenURL mocks base method.
func (m *MockshortenServicer) ShortenURL(arg0 context.Context, arg1 models.UserID, arg2 models.OrigURL) (*models.URLPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShortenURL", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.URLPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShortenURL indicates an expected call of ShortenURL.
func (mr *MockshortenServicerMockRecorder) ShortenURL(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShortenURL", reflect.TypeOf((*MockshortenServicer)(nil).ShortenURL), arg0, arg1, arg2)
}

// MockshortenAuthServicer is a mock of shortenAuthServicer interface.
type MockshortenAuthServicer struct {
	ctrl     *gomock.Controller
	recorder *MockshortenAuthServicerMockRecorder
}

// MockshortenAuthServicerMockRecorder is the mock recorder for MockshortenAuthServicer.
type MockshortenAuthServicerMockRecorder struct {
	mock *MockshortenAuthServicer
}

// NewMockshortenAuthServicer creates a new mock instance.
func NewMockshortenAuthServicer(ctrl *gomock.Controller) *MockshortenAuthServicer {
	mock := &MockshortenAuthServicer{ctrl: ctrl}
	mock.recorder = &MockshortenAuthServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockshortenAuthServicer) EXPECT() *MockshortenAuthServicerMockRecorder {
	return m.recorder
}

// GetUserIDFromCtx mocks base method.
func (m *MockshortenAuthServicer) GetUserIDFromCtx(arg0 context.Context) (models.UserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIDFromCtx", arg0)
	ret0, _ := ret[0].(models.UserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIDFromCtx indicates an expected call of GetUserIDFromCtx.
func (mr *MockshortenAuthServicerMockRecorder) GetUserIDFromCtx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIDFromCtx", reflect.TypeOf((*MockshortenAuthServicer)(nil).GetUserIDFromCtx), arg0)
}

// MockerrShortenConflict is a mock of errShortenConflict interface.
type MockerrShortenConflict struct {
	ctrl     *gomock.Controller
	recorder *MockerrShortenConflictMockRecorder
}

// MockerrShortenConflictMockRecorder is the mock recorder for MockerrShortenConflict.
type MockerrShortenConflictMockRecorder struct {
	mock *MockerrShortenConflict
}

// NewMockerrShortenConflict creates a new mock instance.
func NewMockerrShortenConflict(ctrl *gomock.Controller) *MockerrShortenConflict {
	mock := &MockerrShortenConflict{ctrl: ctrl}
	mock.recorder = &MockerrShortenConflictMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockerrShortenConflict) EXPECT() *MockerrShortenConflictMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockerrShortenConflict) Error() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error")
	ret0, _ := ret[0].(string)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *MockerrShortenConflictMockRecorder) Error() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockerrShortenConflict)(nil).Error))
}

// IsErrConflict mocks base method.
func (m *MockerrShortenConflict) IsErrConflict() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsErrConflict")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsErrConflict indicates an expected call of IsErrConflict.
func (mr *MockerrShortenConflictMockRecorder) IsErrConflict() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsErrConflict", reflect.TypeOf((*MockerrShortenConflict)(nil).IsErrConflict))
}
