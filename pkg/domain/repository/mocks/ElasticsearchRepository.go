// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// ElasticsearchRepository is an autogenerated mock type for the ElasticsearchRepository type
type ElasticsearchRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, index, id, body
func (_m *ElasticsearchRepository) Create(ctx context.Context, index string, id string, body io.Reader) error {
	ret := _m.Called(ctx, index, id, body)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, io.Reader) error); ok {
		r0 = rf(ctx, index, id, body)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, index, id
func (_m *ElasticsearchRepository) Delete(ctx context.Context, index string, id string) error {
	ret := _m.Called(ctx, index, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, index, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, index, id, body
func (_m *ElasticsearchRepository) Update(ctx context.Context, index string, id string, body io.Reader) error {
	ret := _m.Called(ctx, index, id, body)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, io.Reader) error); ok {
		r0 = rf(ctx, index, id, body)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewElasticsearchRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewElasticsearchRepository creates a new instance of ElasticsearchRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewElasticsearchRepository(t mockConstructorTestingTNewElasticsearchRepository) *ElasticsearchRepository {
	mock := &ElasticsearchRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}