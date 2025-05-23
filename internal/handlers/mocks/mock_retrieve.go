// Code generated by MockGen. DO NOT EDIT.
// Source: retrieve.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/rycln/shorturl/internal/models"
)

// MockretrieveServicer is a mock of retrieveServicer interface.
type MockretrieveServicer struct {
	ctrl     *gomock.Controller
	recorder *MockretrieveServicerMockRecorder
}

// MockretrieveServicerMockRecorder is the mock recorder for MockretrieveServicer.
type MockretrieveServicerMockRecorder struct {
	mock *MockretrieveServicer
}

// NewMockretrieveServicer creates a new mock instance.
func NewMockretrieveServicer(ctrl *gomock.Controller) *MockretrieveServicer {
	mock := &MockretrieveServicer{ctrl: ctrl}
	mock.recorder = &MockretrieveServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockretrieveServicer) EXPECT() *MockretrieveServicerMockRecorder {
	return m.recorder
}

// GetOrigURLByShort mocks base method.
func (m *MockretrieveServicer) GetOrigURLByShort(arg0 context.Context, arg1 models.ShortURL) (models.OrigURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrigURLByShort", arg0, arg1)
	ret0, _ := ret[0].(models.OrigURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrigURLByShort indicates an expected call of GetOrigURLByShort.
func (mr *MockretrieveServicerMockRecorder) GetOrigURLByShort(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrigURLByShort", reflect.TypeOf((*MockretrieveServicer)(nil).GetOrigURLByShort), arg0, arg1)
}

// GetShortURLFromCtx mocks base method.
func (m *MockretrieveServicer) GetShortURLFromCtx(arg0 context.Context) (models.ShortURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURLFromCtx", arg0)
	ret0, _ := ret[0].(models.ShortURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURLFromCtx indicates an expected call of GetShortURLFromCtx.
func (mr *MockretrieveServicerMockRecorder) GetShortURLFromCtx(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURLFromCtx", reflect.TypeOf((*MockretrieveServicer)(nil).GetShortURLFromCtx), arg0)
}

// MockerrRetrieveDeletedURL is a mock of errRetrieveDeletedURL interface.
type MockerrRetrieveDeletedURL struct {
	ctrl     *gomock.Controller
	recorder *MockerrRetrieveDeletedURLMockRecorder
}

// MockerrRetrieveDeletedURLMockRecorder is the mock recorder for MockerrRetrieveDeletedURL.
type MockerrRetrieveDeletedURLMockRecorder struct {
	mock *MockerrRetrieveDeletedURL
}

// NewMockerrRetrieveDeletedURL creates a new mock instance.
func NewMockerrRetrieveDeletedURL(ctrl *gomock.Controller) *MockerrRetrieveDeletedURL {
	mock := &MockerrRetrieveDeletedURL{ctrl: ctrl}
	mock.recorder = &MockerrRetrieveDeletedURLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockerrRetrieveDeletedURL) EXPECT() *MockerrRetrieveDeletedURLMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockerrRetrieveDeletedURL) Error() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error")
	ret0, _ := ret[0].(string)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *MockerrRetrieveDeletedURLMockRecorder) Error() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockerrRetrieveDeletedURL)(nil).Error))
}

// IsErrDeletedURL mocks base method.
func (m *MockerrRetrieveDeletedURL) IsErrDeletedURL() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsErrDeletedURL")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsErrDeletedURL indicates an expected call of IsErrDeletedURL.
func (mr *MockerrRetrieveDeletedURLMockRecorder) IsErrDeletedURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsErrDeletedURL", reflect.TypeOf((*MockerrRetrieveDeletedURL)(nil).IsErrDeletedURL))
}
