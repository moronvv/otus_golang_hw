// Code generated by mockery v2.34.2. DO NOT EDIT.

package mockedstorage

import (
	context "context"

	models "github.com/moronvv/otus_golang_hw/hw12_13_14_15_calendar/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockEventStorage is an autogenerated mock type for the EventStorage type
type MockEventStorage struct {
	mock.Mock
}

type MockEventStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEventStorage) EXPECT() *MockEventStorage_Expecter {
	return &MockEventStorage_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *MockEventStorage) Create(_a0 context.Context, _a1 *models.Event) (*models.Event, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Event) (*models.Event, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Event) *models.Event); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Event) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEventStorage_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockEventStorage_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *models.Event
func (_e *MockEventStorage_Expecter) Create(_a0 interface{}, _a1 interface{}) *MockEventStorage_Create_Call {
	return &MockEventStorage_Create_Call{Call: _e.mock.On("Create", _a0, _a1)}
}

func (_c *MockEventStorage_Create_Call) Run(run func(_a0 context.Context, _a1 *models.Event)) *MockEventStorage_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Event))
	})
	return _c
}

func (_c *MockEventStorage_Create_Call) Return(_a0 *models.Event, _a1 error) *MockEventStorage_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEventStorage_Create_Call) RunAndReturn(run func(context.Context, *models.Event) (*models.Event, error)) *MockEventStorage_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *MockEventStorage) Delete(_a0 context.Context, _a1 int64) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockEventStorage_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockEventStorage_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int64
func (_e *MockEventStorage_Expecter) Delete(_a0 interface{}, _a1 interface{}) *MockEventStorage_Delete_Call {
	return &MockEventStorage_Delete_Call{Call: _e.mock.On("Delete", _a0, _a1)}
}

func (_c *MockEventStorage_Delete_Call) Run(run func(_a0 context.Context, _a1 int64)) *MockEventStorage_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockEventStorage_Delete_Call) Return(_a0 error) *MockEventStorage_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEventStorage_Delete_Call) RunAndReturn(run func(context.Context, int64) error) *MockEventStorage_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetMany provides a mock function with given fields: _a0
func (_m *MockEventStorage) GetMany(_a0 context.Context) ([]models.Event, error) {
	ret := _m.Called(_a0)

	var r0 []models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Event, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Event); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEventStorage_GetMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMany'
type MockEventStorage_GetMany_Call struct {
	*mock.Call
}

// GetMany is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *MockEventStorage_Expecter) GetMany(_a0 interface{}) *MockEventStorage_GetMany_Call {
	return &MockEventStorage_GetMany_Call{Call: _e.mock.On("GetMany", _a0)}
}

func (_c *MockEventStorage_GetMany_Call) Run(run func(_a0 context.Context)) *MockEventStorage_GetMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockEventStorage_GetMany_Call) Return(_a0 []models.Event, _a1 error) *MockEventStorage_GetMany_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEventStorage_GetMany_Call) RunAndReturn(run func(context.Context) ([]models.Event, error)) *MockEventStorage_GetMany_Call {
	_c.Call.Return(run)
	return _c
}

// GetOne provides a mock function with given fields: _a0, _a1
func (_m *MockEventStorage) GetOne(_a0 context.Context, _a1 int64) (*models.Event, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*models.Event, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *models.Event); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEventStorage_GetOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOne'
type MockEventStorage_GetOne_Call struct {
	*mock.Call
}

// GetOne is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int64
func (_e *MockEventStorage_Expecter) GetOne(_a0 interface{}, _a1 interface{}) *MockEventStorage_GetOne_Call {
	return &MockEventStorage_GetOne_Call{Call: _e.mock.On("GetOne", _a0, _a1)}
}

func (_c *MockEventStorage_GetOne_Call) Run(run func(_a0 context.Context, _a1 int64)) *MockEventStorage_GetOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockEventStorage_GetOne_Call) Return(_a0 *models.Event, _a1 error) *MockEventStorage_GetOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEventStorage_GetOne_Call) RunAndReturn(run func(context.Context, int64) (*models.Event, error)) *MockEventStorage_GetOne_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockEventStorage) Update(_a0 context.Context, _a1 int64, _a2 *models.Event) (*models.Event, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *models.Event
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, *models.Event) (*models.Event, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, *models.Event) *models.Event); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Event)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, *models.Event) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockEventStorage_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockEventStorage_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int64
//   - _a2 *models.Event
func (_e *MockEventStorage_Expecter) Update(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockEventStorage_Update_Call {
	return &MockEventStorage_Update_Call{Call: _e.mock.On("Update", _a0, _a1, _a2)}
}

func (_c *MockEventStorage_Update_Call) Run(run func(_a0 context.Context, _a1 int64, _a2 *models.Event)) *MockEventStorage_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(*models.Event))
	})
	return _c
}

func (_c *MockEventStorage_Update_Call) Return(_a0 *models.Event, _a1 error) *MockEventStorage_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockEventStorage_Update_Call) RunAndReturn(run func(context.Context, int64, *models.Event) (*models.Event, error)) *MockEventStorage_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockEventStorage creates a new instance of MockEventStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockEventStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockEventStorage {
	mock := &MockEventStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
