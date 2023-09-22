// Code generated by MockGen. DO NOT EDIT.
// Source: mirror.go
//
// Generated by this command:
//
//	mockgen -source=mirror.go -package=worker -destination=mock_mirror.go
//
// Package worker is a generated GoMock package.
package worker

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMirrorResolver is a mock of MirrorResolver interface.
type MockMirrorResolver struct {
	ctrl     *gomock.Controller
	recorder *MockMirrorResolverMockRecorder
}

// MockMirrorResolverMockRecorder is the mock recorder for MockMirrorResolver.
type MockMirrorResolverMockRecorder struct {
	mock *MockMirrorResolver
}

// NewMockMirrorResolver creates a new mock instance.
func NewMockMirrorResolver(ctrl *gomock.Controller) *MockMirrorResolver {
	mock := &MockMirrorResolver{ctrl: ctrl}
	mock.recorder = &MockMirrorResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMirrorResolver) EXPECT() *MockMirrorResolverMockRecorder {
	return m.recorder
}

// GetAllReferences mocks base method.
func (m *MockMirrorResolver) GetAllReferences(imageName string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllReferences", imageName)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllReferences indicates an expected call of GetAllReferences.
func (mr *MockMirrorResolverMockRecorder) GetAllReferences(imageName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllReferences", reflect.TypeOf((*MockMirrorResolver)(nil).GetAllReferences), imageName)
}
