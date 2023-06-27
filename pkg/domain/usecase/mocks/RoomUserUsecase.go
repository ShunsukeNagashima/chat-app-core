// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/shunsukenagashima/chat-api/pkg/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// RoomUserUsecase is an autogenerated mock type for the RoomUserUsecase type
type RoomUserUsecase struct {
	mock.Mock
}

// AddUsersToRoom provides a mock function with given fields: ctx, roomId, userIds
func (_m *RoomUserUsecase) AddUsersToRoom(ctx context.Context, roomId string, userIds []string) error {
	ret := _m.Called(ctx, roomId, userIds)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []string) error); ok {
		r0 = rf(ctx, roomId, userIds)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllRoomsByUserID provides a mock function with given fields: ctx, userId
func (_m *RoomUserUsecase) GetAllRoomsByUserID(ctx context.Context, userId string) ([]*model.Room, error) {
	ret := _m.Called(ctx, userId)

	var r0 []*model.Room
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.Room, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Room); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Room)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveUserFromRoom provides a mock function with given fields: ctx, roomId, userId
func (_m *RoomUserUsecase) RemoveUserFromRoom(ctx context.Context, roomId string, userId string) error {
	ret := _m.Called(ctx, roomId, userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, roomId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRoomUserUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewRoomUserUsecase creates a new instance of RoomUserUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRoomUserUsecase(t mockConstructorTestingTNewRoomUserUsecase) *RoomUserUsecase {
	mock := &RoomUserUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
