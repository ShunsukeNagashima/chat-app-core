// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/shunsukenagashima/chat-api/pkg/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// RoomUserRepository is an autogenerated mock type for the RoomUserRepository type
type RoomUserRepository struct {
	mock.Mock
}

// AddUsersToRoom provides a mock function with given fields: ctx, roomId, userIDs
func (_m *RoomUserRepository) AddUsersToRoom(ctx context.Context, roomId string, userIDs []string) error {
	ret := _m.Called(ctx, roomId, userIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []string) error); ok {
		r0 = rf(ctx, roomId, userIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllRoomsByUserID provides a mock function with given fields: ctx, userId
func (_m *RoomUserRepository) GetAllRoomsByUserID(ctx context.Context, userId string) ([]*model.RoomUser, error) {
	ret := _m.Called(ctx, userId)

	var r0 []*model.RoomUser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.RoomUser, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.RoomUser); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.RoomUser)
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
func (_m *RoomUserRepository) RemoveUserFromRoom(ctx context.Context, roomId string, userId string) error {
	ret := _m.Called(ctx, roomId, userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, roomId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRoomUserRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRoomUserRepository creates a new instance of RoomUserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRoomUserRepository(t mockConstructorTestingTNewRoomUserRepository) *RoomUserRepository {
	mock := &RoomUserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
