// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/shunsukenagashima/chat-api/pkg/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// RoomUsecase is an autogenerated mock type for the RoomUsecase type
type RoomUsecase struct {
	mock.Mock
}

// CreateRoom provides a mock function with given fields: ctx, room, ownerID
func (_m *RoomUsecase) CreateRoom(ctx context.Context, room *model.Room, ownerID string) error {
	ret := _m.Called(ctx, room, ownerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Room, string) error); ok {
		r0 = rf(ctx, room, ownerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteRoom provides a mock function with given fields: ctx, roomID
func (_m *RoomUsecase) DeleteRoom(ctx context.Context, roomID string) error {
	ret := _m.Called(ctx, roomID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, roomID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllPublicRoom provides a mock function with given fields: ctx
func (_m *RoomUsecase) GetAllPublicRoom(ctx context.Context) ([]*model.Room, error) {
	ret := _m.Called(ctx)

	var r0 []*model.Room
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.Room, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Room); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Room)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRoomByID provides a mock function with given fields: ctx, roomID
func (_m *RoomUsecase) GetRoomByID(ctx context.Context, roomID string) (*model.Room, error) {
	ret := _m.Called(ctx, roomID)

	var r0 *model.Room
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Room, error)); ok {
		return rf(ctx, roomID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Room); ok {
		r0 = rf(ctx, roomID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Room)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, roomID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRoom provides a mock function with given fields: ctx, room
func (_m *RoomUsecase) UpdateRoom(ctx context.Context, room *model.Room) error {
	ret := _m.Called(ctx, room)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Room) error); ok {
		r0 = rf(ctx, room)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRoomUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewRoomUsecase creates a new instance of RoomUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRoomUsecase(t mockConstructorTestingTNewRoomUsecase) *RoomUsecase {
	mock := &RoomUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
