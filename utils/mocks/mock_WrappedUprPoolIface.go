// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/v8/base"
	mock "github.com/stretchr/testify/mock"
)

// WrappedUprPoolIface is an autogenerated mock type for the WrappedUprPoolIface type
type WrappedUprPoolIface struct {
	mock.Mock
}

type WrappedUprPoolIface_Expecter struct {
	mock *mock.Mock
}

func (_m *WrappedUprPoolIface) EXPECT() *WrappedUprPoolIface_Expecter {
	return &WrappedUprPoolIface_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with no fields
func (_m *WrappedUprPoolIface) Get() *base.WrappedUprEvent {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *base.WrappedUprEvent
	if rf, ok := ret.Get(0).(func() *base.WrappedUprEvent); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*base.WrappedUprEvent)
		}
	}

	return r0
}

// WrappedUprPoolIface_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type WrappedUprPoolIface_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
func (_e *WrappedUprPoolIface_Expecter) Get() *WrappedUprPoolIface_Get_Call {
	return &WrappedUprPoolIface_Get_Call{Call: _e.mock.On("Get")}
}

func (_c *WrappedUprPoolIface_Get_Call) Run(run func()) *WrappedUprPoolIface_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *WrappedUprPoolIface_Get_Call) Return(_a0 *base.WrappedUprEvent) *WrappedUprPoolIface_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *WrappedUprPoolIface_Get_Call) RunAndReturn(run func() *base.WrappedUprEvent) *WrappedUprPoolIface_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Put provides a mock function with given fields: uprEvent
func (_m *WrappedUprPoolIface) Put(uprEvent *base.WrappedUprEvent) {
	_m.Called(uprEvent)
}

// WrappedUprPoolIface_Put_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Put'
type WrappedUprPoolIface_Put_Call struct {
	*mock.Call
}

// Put is a helper method to define mock.On call
//   - uprEvent *base.WrappedUprEvent
func (_e *WrappedUprPoolIface_Expecter) Put(uprEvent interface{}) *WrappedUprPoolIface_Put_Call {
	return &WrappedUprPoolIface_Put_Call{Call: _e.mock.On("Put", uprEvent)}
}

func (_c *WrappedUprPoolIface_Put_Call) Run(run func(uprEvent *base.WrappedUprEvent)) *WrappedUprPoolIface_Put_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*base.WrappedUprEvent))
	})
	return _c
}

func (_c *WrappedUprPoolIface_Put_Call) Return() *WrappedUprPoolIface_Put_Call {
	_c.Call.Return()
	return _c
}

func (_c *WrappedUprPoolIface_Put_Call) RunAndReturn(run func(*base.WrappedUprEvent)) *WrappedUprPoolIface_Put_Call {
	_c.Run(run)
	return _c
}

// NewWrappedUprPoolIface creates a new instance of WrappedUprPoolIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewWrappedUprPoolIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *WrappedUprPoolIface {
	mock := &WrappedUprPoolIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
