// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/HardDie/mmr_boost_server/internal/dto"
	entity "github.com/HardDie/mmr_boost_server/internal/entity"

	mock "github.com/stretchr/testify/mock"
)

// IPostgresApplication is an autogenerated mock type for the IPostgresApplication type
type IPostgresApplication struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, req
func (_m *IPostgresApplication) Create(ctx context.Context, req *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error) {
	ret := _m.Called(ctx, req)

	var r0 *entity.ApplicationPublic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationCreateRequest) (*entity.ApplicationPublic, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationCreateRequest) *entity.ApplicationPublic); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ApplicationPublic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.ApplicationCreateRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Item provides a mock function with given fields: ctx, req
func (_m *IPostgresApplication) Item(ctx context.Context, req *dto.ApplicationItemRequest) (*entity.ApplicationPublic, error) {
	ret := _m.Called(ctx, req)

	var r0 *entity.ApplicationPublic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationItemRequest) (*entity.ApplicationPublic, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationItemRequest) *entity.ApplicationPublic); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ApplicationPublic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.ApplicationItemRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, req
func (_m *IPostgresApplication) List(ctx context.Context, req *dto.ApplicationListRequest) ([]*entity.ApplicationPublic, error) {
	ret := _m.Called(ctx, req)

	var r0 []*entity.ApplicationPublic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationListRequest) ([]*entity.ApplicationPublic, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationListRequest) []*entity.ApplicationPublic); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.ApplicationPublic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.ApplicationListRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrivateItem provides a mock function with given fields: ctx, req
func (_m *IPostgresApplication) PrivateItem(ctx context.Context, req *dto.ApplicationItemRequest) (*entity.ApplicationPrivate, error) {
	ret := _m.Called(ctx, req)

	var r0 *entity.ApplicationPrivate
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationItemRequest) (*entity.ApplicationPrivate, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationItemRequest) *entity.ApplicationPrivate); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ApplicationPrivate)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.ApplicationItemRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateItem provides a mock function with given fields: ctx, req
func (_m *IPostgresApplication) UpdateItem(ctx context.Context, req *dto.ApplicationUpdateRequest) (*entity.ApplicationPublic, error) {
	ret := _m.Called(ctx, req)

	var r0 *entity.ApplicationPublic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationUpdateRequest) (*entity.ApplicationPublic, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationUpdateRequest) *entity.ApplicationPublic); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ApplicationPublic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.ApplicationUpdateRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePrivate provides a mock function with given fields: ctx, req
func (_m *IPostgresApplication) UpdatePrivate(ctx context.Context, req *dto.ApplicationUpdatePrivateRequest) (*entity.ApplicationPrivate, error) {
	ret := _m.Called(ctx, req)

	var r0 *entity.ApplicationPrivate
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationUpdatePrivateRequest) (*entity.ApplicationPrivate, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationUpdatePrivateRequest) *entity.ApplicationPrivate); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ApplicationPrivate)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.ApplicationUpdatePrivateRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: ctx, req
func (_m *IPostgresApplication) UpdateStatus(ctx context.Context, req *dto.ApplicationUpdateStatusRequest) (*entity.ApplicationPublic, error) {
	ret := _m.Called(ctx, req)

	var r0 *entity.ApplicationPublic
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationUpdateStatusRequest) (*entity.ApplicationPublic, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.ApplicationUpdateStatusRequest) *entity.ApplicationPublic); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.ApplicationPublic)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.ApplicationUpdateStatusRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIPostgresApplication interface {
	mock.TestingT
	Cleanup(func())
}

// NewIPostgresApplication creates a new instance of IPostgresApplication. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIPostgresApplication(t mockConstructorTestingTNewIPostgresApplication) *IPostgresApplication {
	mock := &IPostgresApplication{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
