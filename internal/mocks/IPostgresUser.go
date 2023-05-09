// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/HardDie/mmr_boost_server/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// IPostgresUser is an autogenerated mock type for the IPostgresUser type
type IPostgresUser struct {
	mock.Mock
}

// ActivateRecord provides a mock function with given fields: ctx, userID
func (_m *IPostgresUser) ActivateRecord(ctx context.Context, userID int32) (*entity.User, error) {
	ret := _m.Called(ctx, userID)

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32) (*entity.User, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32) *entity.User); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, email, name
func (_m *IPostgresUser) Create(ctx context.Context, email string, name string) (*entity.User, error) {
	ret := _m.Called(ctx, email, name)

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*entity.User, error)); ok {
		return rf(ctx, email, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *entity.User); ok {
		r0 = rf(ctx, email, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, email, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *IPostgresUser) GetByID(ctx context.Context, id int32) (*entity.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32) (*entity.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32) *entity.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: ctx, name
func (_m *IPostgresUser) GetByName(ctx context.Context, name string) (*entity.User, error) {
	ret := _m.Called(ctx, name)

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*entity.User, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *entity.User); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSteamID provides a mock function with given fields: ctx, userID, steamID
func (_m *IPostgresUser) UpdateSteamID(ctx context.Context, userID int32, steamID string) (*entity.User, error) {
	ret := _m.Called(ctx, userID, steamID)

	var r0 *entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int32, string) (*entity.User, error)); ok {
		return rf(ctx, userID, steamID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int32, string) *entity.User); ok {
		r0 = rf(ctx, userID, steamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int32, string) error); ok {
		r1 = rf(ctx, userID, steamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIPostgresUser interface {
	mock.TestingT
	Cleanup(func())
}

// NewIPostgresUser creates a new instance of IPostgresUser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIPostgresUser(t mockConstructorTestingTNewIPostgresUser) *IPostgresUser {
	mock := &IPostgresUser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
