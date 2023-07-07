// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/shunsukenagashima/chat-api/pkg/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// MessageRepository is an autogenerated mock type for the MessageRepository type
type MessageRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, message
func (_m *MessageRepository) Create(ctx context.Context, message *model.Message) error {
	ret := _m.Called(ctx, message)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Message) error); ok {
		r0 = rf(ctx, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, messageId
func (_m *MessageRepository) Delete(ctx context.Context, messageId string) error {
	ret := _m.Called(ctx, messageId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, messageId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllMessagesByRoomID provides a mock function with given fields: ctx, roomId
func (_m *MessageRepository) GetAllMessagesByRoomID(ctx context.Context, roomId string) ([]*model.Message, error) {
	ret := _m.Called(ctx, roomId)

	var r0 []*model.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*model.Message, error)); ok {
		return rf(ctx, roomId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Message); ok {
		r0 = rf(ctx, roomId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, roomId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, messageId, newContent
func (_m *MessageRepository) Update(ctx context.Context, messageId string, newContent string) error {
	ret := _m.Called(ctx, messageId, newContent)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, messageId, newContent)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMessageRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageRepository creates a new instance of MessageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageRepository(t mockConstructorTestingTNewMessageRepository) *MessageRepository {
	mock := &MessageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
