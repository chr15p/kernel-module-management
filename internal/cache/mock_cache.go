// Code generated by MockGen. DO NOT EDIT.
// Source: cache.go

// Package cache is a generated GoMock package.
package cache

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache[T comparable] struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder[T]
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder[T comparable] struct {
	mock *MockCache[T]
}

// NewMockCache creates a new mock instance.
func NewMockCache[T comparable](ctrl *gomock.Controller) *MockCache[T] {
	mock := &MockCache[T]{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache[T]) EXPECT() *MockCacheMockRecorder[T] {
	return m.recorder
}

// DeleteExpired mocks base method.
func (m *MockCache[T]) DeleteExpired() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteExpired")
}

// DeleteExpired indicates an expected call of DeleteExpired.
func (mr *MockCacheMockRecorder[T]) DeleteExpired() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpired", reflect.TypeOf((*MockCache[T])(nil).DeleteExpired))
}

// Get mocks base method.
func (m *MockCache[T]) Get(key T) (interface{}, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder[T]) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache[T])(nil).Get), key)
}

// Set mocks base method.
func (m *MockCache[T]) Set(key T, value interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", key, value)
}

// Set indicates an expected call of Set.
func (mr *MockCacheMockRecorder[T]) Set(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCache[T])(nil).Set), key, value)
}

// StartCollecting mocks base method.
func (m *MockCache[T]) StartCollecting(ctx context.Context, interval time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartCollecting", ctx, interval)
}

// StartCollecting indicates an expected call of StartCollecting.
func (mr *MockCacheMockRecorder[T]) StartCollecting(ctx, interval interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartCollecting", reflect.TypeOf((*MockCache[T])(nil).StartCollecting), ctx, interval)
}

// WaitForTermination mocks base method.
func (m *MockCache[T]) WaitForTermination() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WaitForTermination")
}

// WaitForTermination indicates an expected call of WaitForTermination.
func (mr *MockCacheMockRecorder[T]) WaitForTermination() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForTermination", reflect.TypeOf((*MockCache[T])(nil).WaitForTermination))
}