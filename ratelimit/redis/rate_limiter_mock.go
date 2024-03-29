// Code generated by MockGen. DO NOT EDIT.
// Source: ratelimiter.go

// Package redis is a generated GoMock package.
package redis

import (
	context "context"
	reflect "reflect"

	redis_rate "github.com/go-redis/redis_rate/v9"
	gomock "github.com/golang/mock/gomock"
)

// MockactionAllower is a mock of actionAllower interface.
type MockactionAllower struct {
	ctrl     *gomock.Controller
	recorder *MockactionAllowerMockRecorder
}

// MockactionAllowerMockRecorder is the mock recorder for MockactionAllower.
type MockactionAllowerMockRecorder struct {
	mock *MockactionAllower
}

// NewMockactionAllower creates a new mock instance.
func NewMockactionAllower(ctrl *gomock.Controller) *MockactionAllower {
	mock := &MockactionAllower{ctrl: ctrl}
	mock.recorder = &MockactionAllowerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockactionAllower) EXPECT() *MockactionAllowerMockRecorder {
	return m.recorder
}

// Allow mocks base method.
func (m *MockactionAllower) Allow(ctx context.Context, key string, limit redis_rate.Limit) (*redis_rate.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Allow", ctx, key, limit)
	ret0, _ := ret[0].(*redis_rate.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Allow indicates an expected call of Allow.
func (mr *MockactionAllowerMockRecorder) Allow(ctx, key, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Allow", reflect.TypeOf((*MockactionAllower)(nil).Allow), ctx, key, limit)
}
