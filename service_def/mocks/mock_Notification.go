// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	mock "github.com/stretchr/testify/mock"
)

// Notification is an autogenerated mock type for the Notification type
type Notification struct {
	mock.Mock
}

type Notification_Expecter struct {
	mock *mock.Mock
}

func (_m *Notification) EXPECT() *Notification_Expecter {
	return &Notification_Expecter{mock: &_m.Mock}
}

// Clone provides a mock function with given fields: numOfReaders
func (_m *Notification) Clone(numOfReaders int) interface{} {
	ret := _m.Called(numOfReaders)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int) interface{}); ok {
		r0 = rf(numOfReaders)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Notification_Clone_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clone'
type Notification_Clone_Call struct {
	*mock.Call
}

// Clone is a helper method to define mock.On call
//   - numOfReaders int
func (_e *Notification_Expecter) Clone(numOfReaders interface{}) *Notification_Clone_Call {
	return &Notification_Clone_Call{Call: _e.mock.On("Clone", numOfReaders)}
}

func (_c *Notification_Clone_Call) Run(run func(numOfReaders int)) *Notification_Clone_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *Notification_Clone_Call) Return(_a0 interface{}) *Notification_Clone_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Notification_Clone_Call) RunAndReturn(run func(int) interface{}) *Notification_Clone_Call {
	_c.Call.Return(run)
	return _c
}

// GetReplicasInfo provides a mock function with given fields:
func (_m *Notification) GetReplicasInfo() (int, *base.VbHostsMapType, *base.StringStringMap, []uint16) {
	ret := _m.Called()

	var r0 int
	var r1 *base.VbHostsMapType
	var r2 *base.StringStringMap
	var r3 []uint16
	if rf, ok := ret.Get(0).(func() (int, *base.VbHostsMapType, *base.StringStringMap, []uint16)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() *base.VbHostsMapType); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*base.VbHostsMapType)
		}
	}

	if rf, ok := ret.Get(2).(func() *base.StringStringMap); ok {
		r2 = rf()
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*base.StringStringMap)
		}
	}

	if rf, ok := ret.Get(3).(func() []uint16); ok {
		r3 = rf()
	} else {
		if ret.Get(3) != nil {
			r3 = ret.Get(3).([]uint16)
		}
	}

	return r0, r1, r2, r3
}

// Notification_GetReplicasInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReplicasInfo'
type Notification_GetReplicasInfo_Call struct {
	*mock.Call
}

// GetReplicasInfo is a helper method to define mock.On call
func (_e *Notification_Expecter) GetReplicasInfo() *Notification_GetReplicasInfo_Call {
	return &Notification_GetReplicasInfo_Call{Call: _e.mock.On("GetReplicasInfo")}
}

func (_c *Notification_GetReplicasInfo_Call) Run(run func()) *Notification_GetReplicasInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Notification_GetReplicasInfo_Call) Return(_a0 int, _a1 *base.VbHostsMapType, _a2 *base.StringStringMap, _a3 []uint16) *Notification_GetReplicasInfo_Call {
	_c.Call.Return(_a0, _a1, _a2, _a3)
	return _c
}

func (_c *Notification_GetReplicasInfo_Call) RunAndReturn(run func() (int, *base.VbHostsMapType, *base.StringStringMap, []uint16)) *Notification_GetReplicasInfo_Call {
	_c.Call.Return(run)
	return _c
}

// IsSourceNotification provides a mock function with given fields:
func (_m *Notification) IsSourceNotification() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Notification_IsSourceNotification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsSourceNotification'
type Notification_IsSourceNotification_Call struct {
	*mock.Call
}

// IsSourceNotification is a helper method to define mock.On call
func (_e *Notification_Expecter) IsSourceNotification() *Notification_IsSourceNotification_Call {
	return &Notification_IsSourceNotification_Call{Call: _e.mock.On("IsSourceNotification")}
}

func (_c *Notification_IsSourceNotification_Call) Run(run func()) *Notification_IsSourceNotification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Notification_IsSourceNotification_Call) Return(_a0 bool) *Notification_IsSourceNotification_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Notification_IsSourceNotification_Call) RunAndReturn(run func() bool) *Notification_IsSourceNotification_Call {
	_c.Call.Return(run)
	return _c
}

// Recycle provides a mock function with given fields:
func (_m *Notification) Recycle() {
	_m.Called()
}

// Notification_Recycle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Recycle'
type Notification_Recycle_Call struct {
	*mock.Call
}

// Recycle is a helper method to define mock.On call
func (_e *Notification_Expecter) Recycle() *Notification_Recycle_Call {
	return &Notification_Recycle_Call{Call: _e.mock.On("Recycle")}
}

func (_c *Notification_Recycle_Call) Run(run func()) *Notification_Recycle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Notification_Recycle_Call) Return() *Notification_Recycle_Call {
	_c.Call.Return()
	return _c
}

func (_c *Notification_Recycle_Call) RunAndReturn(run func()) *Notification_Recycle_Call {
	_c.Call.Return(run)
	return _c
}

// NewNotification creates a new instance of Notification. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewNotification(t interface {
	mock.TestingT
	Cleanup(func())
}) *Notification {
	mock := &Notification{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}