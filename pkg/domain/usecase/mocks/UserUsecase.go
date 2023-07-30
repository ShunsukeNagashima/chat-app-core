// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/shunsukenagashima/chat-api/pkg/domain/model"
	mock "github.com/stretchr/testify/mock"
)

// UserUsecase is an autogenerated mock type for the UserUsecase type
type UserUsecase struct {
	mock.Mock
}

// BatchGetUsers provides a mock function with given fields: ctx, userIds
func (_m *UserUsecase) BatchGetUsers(ctx context.Context, userIds []string) ([]*model.User, error) {
	ret := _m.Called(ctx, userIds)

	var r0 []*model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) ([]*model.User, error)); ok {
		return rf(ctx, userIds)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) []*model.User); ok {
		r0 = rf(ctx, userIds)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, userIds)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, user, idToken
func (_m *UserUsecase) CreateUser(ctx context.Context, user *model.User, idToken string) error {
	ret := _m.Called(ctx, user, idToken)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User, string) error); ok {
		r0 = rf(ctx, user, idToken)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserByID provides a mock function with given fields: ctx, userId
func (_m *UserUsecase) GetUserByID(ctx context.Context, userId string) (*model.User, error) {
	ret := _m.Called(ctx, userId)

	var r0 *model.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.User, error)); ok {
		return rf(ctx, userId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.User); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchUsers provides a mock function with given fields: ctx, query, nextKey, size
func (_m *UserUsecase) SearchUsers(ctx context.Context, query string, nextKey string, size int) ([]*model.User, string, error) {
	ret := _m.Called(ctx, query, nextKey, size)

	var r0 []*model.User
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) ([]*model.User, string, error)); ok {
		return rf(ctx, query, nextKey, size)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) []*model.User); ok {
		r0 = rf(ctx, query, nextKey, size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, int) string); ok {
		r1 = rf(ctx, query, nextKey, size)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string, int) error); ok {
		r2 = rf(ctx, query, nextKey, size)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewUserUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserUsecase creates a new instance of UserUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserUsecase(t mockConstructorTestingTNewUserUsecase) *UserUsecase {
	mock := &UserUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
