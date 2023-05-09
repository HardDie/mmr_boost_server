// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/HardDie/mmr_boost_server/internal/entity"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// IPostgresEmailValidation is an autogenerated mock type for the IPostgresEmailValidation type
type IPostgresEmailValidation struct {
	mock.Mock
}

// CreateOrUpdate provides a mock function with given fields: ctx, userID, code, expiredAt
func (_m *IPostgresEmailValidation) CreateOrUpdate(ctx context.Context, userID int32, code string, expiredAt time.Time) (*entity.EmailValidation, error) {
	ret := _m.Called(ctx, userID, code, expiredAt)

	var r0 *entity.EmailValidation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, string, time.Time) (*entity.EmailValidation, error)); ok {
		return rf(ctx, userID, code, expiredAt)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32, string, time.Time) *entity.EmailValidation); ok {
		r0 = rf(ctx, userID, code, expiredAt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.EmailValidation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32, string, time.Time) error); ok {
		r1 = rf(ctx, userID, code, expiredAt)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByID provides a mock function with given fields: ctx, id
func (_m *IPostgresEmailValidation) DeleteByID(ctx context.Context, id int32) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int32) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByCode provides a mock function with given fields: ctx, code
func (_m *IPostgresEmailValidation) GetByCode(ctx context.Context, code string) (*entity.EmailValidation, error) {
	ret := _m.Called(ctx, code)

	var r0 *entity.EmailValidation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.EmailValidation, error)); ok {
		return rf(ctx, code)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.EmailValidation); ok {
		r0 = rf(ctx, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.EmailValidation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIPostgresEmailValidation interface {
	mock.TestingT
	Cleanup(func())
}

// NewIPostgresEmailValidation creates a new instance of IPostgresEmailValidation. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIPostgresEmailValidation(t mockConstructorTestingTNewIPostgresEmailValidation) *IPostgresEmailValidation {
	mock := &IPostgresEmailValidation{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
