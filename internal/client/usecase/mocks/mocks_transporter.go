// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	proto "github.com/h2p2f/dedicated-vault/proto"
	mock "github.com/stretchr/testify/mock"
)

// Transporter is an autogenerated mock type for the Transporter type
type Transporter struct {
	mock.Mock
}

// ChangePassword provides a mock function with given fields: ctx, user, newPassword
func (_m *Transporter) ChangePassword(ctx context.Context, user *proto.User, newPassword string) (string, error) {
	ret := _m.Called(ctx, user, newPassword)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.User, string) (string, error)); ok {
		return rf(ctx, user, newPassword)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.User, string) string); ok {
		r0 = rf(ctx, user, newPassword)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.User, string) error); ok {
		r1 = rf(ctx, user, newPassword)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChangeSecret provides a mock function with given fields: ctx, data
func (_m *Transporter) ChangeSecret(ctx context.Context, data *proto.SecretData) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.SecretData) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSecret provides a mock function with given fields: ctx, uuid
func (_m *Transporter) DeleteSecret(ctx context.Context, uuid string) error {
	ret := _m.Called(ctx, uuid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, uuid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListSecrets provides a mock function with given fields: ctx
func (_m *Transporter) ListSecrets(ctx context.Context) ([]*proto.SecretData, error) {
	ret := _m.Called(ctx)

	var r0 []*proto.SecretData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*proto.SecretData, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*proto.SecretData); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*proto.SecretData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: ctx, user
func (_m *Transporter) Login(ctx context.Context, user *proto.User) (string, error) {
	ret := _m.Called(ctx, user)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.User) (string, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.User) string); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, user
func (_m *Transporter) Register(ctx context.Context, user *proto.User) (string, error) {
	ret := _m.Called(ctx, user)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.User) (string, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.User) string); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveSecret provides a mock function with given fields: ctx, data
func (_m *Transporter) SaveSecret(ctx context.Context, data *proto.SecretData) error {
	ret := _m.Called(ctx, data)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.SecretData) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransporter creates a new instance of Transporter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransporter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Transporter {
	mock := &Transporter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
