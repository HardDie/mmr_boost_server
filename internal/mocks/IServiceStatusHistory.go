// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/HardDie/mmr_boost_server/internal/dto"
	entity "github.com/HardDie/mmr_boost_server/internal/entity"

	mock "github.com/stretchr/testify/mock"
)

// IServiceStatusHistory is an autogenerated mock type for the IServiceStatusHistory type
type IServiceStatusHistory struct {
	mock.Mock
}

// StatusHistory provides a mock function with given fields: ctx, req
func (_m *IServiceStatusHistory) StatusHistory(ctx context.Context, req *dto.StatusHistoryListRequest) ([]*entity.StatusHistory, error) {
	ret := _m.Called(ctx, req)

	var r0 []*entity.StatusHistory
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *dto.StatusHistoryListRequest) ([]*entity.StatusHistory, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *dto.StatusHistoryListRequest) []*entity.StatusHistory); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.StatusHistory)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *dto.StatusHistoryListRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIServiceStatusHistory interface {
	mock.TestingT
	Cleanup(func())
}

// NewIServiceStatusHistory creates a new instance of IServiceStatusHistory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIServiceStatusHistory(t mockConstructorTestingTNewIServiceStatusHistory) *IServiceStatusHistory {
	mock := &IServiceStatusHistory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
