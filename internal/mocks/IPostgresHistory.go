// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IPostgresHistory is an autogenerated mock type for the IPostgresHistory type
type IPostgresHistory struct {
	mock.Mock
}

// NewEvent provides a mock function with given fields: ctx, userID, message
func (_m *IPostgresHistory) NewEvent(ctx context.Context, userID int32, message string) error {
	ret := _m.Called(ctx, userID, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, string) error); ok {
		r0 = rf(ctx, userID, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEventWithAffected provides a mock function with given fields: ctx, userID, affectedUserID, message
func (_m *IPostgresHistory) NewEventWithAffected(ctx context.Context, userID int32, affectedUserID int32, message string) error {
	ret := _m.Called(ctx, userID, affectedUserID, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, int32, string) error); ok {
		r0 = rf(ctx, userID, affectedUserID, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewIPostgresHistory interface {
	mock.TestingT
	Cleanup(func())
}

// NewIPostgresHistory creates a new instance of IPostgresHistory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIPostgresHistory(t mockConstructorTestingTNewIPostgresHistory) *IPostgresHistory {
	mock := &IPostgresHistory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}