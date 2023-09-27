// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/h2p2f/dedicated-vault/internal/server/models"
)

// DataHandler is an autogenerated mock type for the DataHandler type
type DataHandler struct {
	mock.Mock
}

// ChangeData provides a mock function with given fields: ctx, user, data
func (_m *DataHandler) ChangeData(ctx context.Context, user models.User, data models.VaultData) (int64, error) {
	ret := _m.Called(ctx, user, data)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User, models.VaultData) (int64, error)); ok {
		return rf(ctx, user, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User, models.VaultData) int64); ok {
		r0 = rf(ctx, user, data)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User, models.VaultData) error); ok {
		r1 = rf(ctx, user, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateData provides a mock function with given fields: ctx, user, data
func (_m *DataHandler) CreateData(ctx context.Context, user models.User, data models.VaultData) (string, int64, error) {
	ret := _m.Called(ctx, user, data)

	var r0 string
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User, models.VaultData) (string, int64, error)); ok {
		return rf(ctx, user, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User, models.VaultData) string); ok {
		r0 = rf(ctx, user, data)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User, models.VaultData) int64); ok {
		r1 = rf(ctx, user, data)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, models.User, models.VaultData) error); ok {
		r2 = rf(ctx, user, data)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DeleteData provides a mock function with given fields: ctx, user, data
func (_m *DataHandler) DeleteData(ctx context.Context, user models.User, data models.VaultData) (int64, error) {
	ret := _m.Called(ctx, user, data)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User, models.VaultData) (int64, error)); ok {
		return rf(ctx, user, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User, models.VaultData) int64); ok {
		r0 = rf(ctx, user, data)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User, models.VaultData) error); ok {
		r1 = rf(ctx, user, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllData provides a mock function with given fields: ctx, user
func (_m *DataHandler) GetAllData(ctx context.Context, user models.User) ([]models.VaultData, error) {
	ret := _m.Called(ctx, user)

	var r0 []models.VaultData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) ([]models.VaultData, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.User) []models.VaultData); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.VaultData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDataHandler creates a new instance of DataHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataHandler {
	mock := &DataHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
