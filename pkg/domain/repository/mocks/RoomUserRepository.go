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

// AddUserToRoom provides a mock function with given fields: ctx, roomID, userID
func (_m *RoomUserRepository) AddUserToRoom(ctx context.Context, roomID string, userID string) error {
	ret := _m.Called(ctx, roomID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, roomID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AddUsersToRoom provides a mock function with given fields: ctx, roomID, userIDs
func (_m *RoomUserRepository) AddUsersToRoom(ctx context.Context, roomID string, userIDs []string) error {
	ret := _m.Called(ctx, roomID, userIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []string) error); ok {
		r0 = rf(ctx, roomID, userIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllRoomsByUserID provides a mock function with given fields: ctx, userID
func (_m *RoomUserRepository) GetAllRoomsByUserID(ctx context.Context, userID string) ([]*model.Room, error) {
	ret := _m.Called(ctx, userID)

	var r0 []*model.Room
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.Room, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Room); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Room)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveUserFromRoom provides a mock function with given fields: ctx, roomID, userID
func (_m *RoomUserRepository) RemoveUserFromRoom(ctx context.Context, roomID string, userID string) error {
	ret := _m.Called(ctx, roomID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, roomID, userID)
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
