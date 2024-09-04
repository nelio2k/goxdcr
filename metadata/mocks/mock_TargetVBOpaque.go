// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	metadata "github.com/couchbase/goxdcr/v8/metadata"
	mock "github.com/stretchr/testify/mock"
)

// TargetVBOpaque is an autogenerated mock type for the TargetVBOpaque type
type TargetVBOpaque struct {
	mock.Mock
}

type TargetVBOpaque_Expecter struct {
	mock *mock.Mock
}

func (_m *TargetVBOpaque) EXPECT() *TargetVBOpaque_Expecter {
	return &TargetVBOpaque_Expecter{mock: &_m.Mock}
}

// IsSame provides a mock function with given fields: targetVBOpaque
func (_m *TargetVBOpaque) IsSame(targetVBOpaque metadata.TargetVBOpaque) bool {
	ret := _m.Called(targetVBOpaque)

	if len(ret) == 0 {
		panic("no return value specified for IsSame")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(metadata.TargetVBOpaque) bool); ok {
		r0 = rf(targetVBOpaque)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TargetVBOpaque_IsSame_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsSame'
type TargetVBOpaque_IsSame_Call struct {
	*mock.Call
}

// IsSame is a helper method to define mock.On call
//   - targetVBOpaque metadata.TargetVBOpaque
func (_e *TargetVBOpaque_Expecter) IsSame(targetVBOpaque interface{}) *TargetVBOpaque_IsSame_Call {
	return &TargetVBOpaque_IsSame_Call{Call: _e.mock.On("IsSame", targetVBOpaque)}
}

func (_c *TargetVBOpaque_IsSame_Call) Run(run func(targetVBOpaque metadata.TargetVBOpaque)) *TargetVBOpaque_IsSame_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.TargetVBOpaque))
	})
	return _c
}

func (_c *TargetVBOpaque_IsSame_Call) Return(_a0 bool) *TargetVBOpaque_IsSame_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetVBOpaque_IsSame_Call) RunAndReturn(run func(metadata.TargetVBOpaque) bool) *TargetVBOpaque_IsSame_Call {
	_c.Call.Return(run)
	return _c
}

// Size provides a mock function with given fields:
func (_m *TargetVBOpaque) Size() int {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Size")
	}

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// TargetVBOpaque_Size_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Size'
type TargetVBOpaque_Size_Call struct {
	*mock.Call
}

// Size is a helper method to define mock.On call
func (_e *TargetVBOpaque_Expecter) Size() *TargetVBOpaque_Size_Call {
	return &TargetVBOpaque_Size_Call{Call: _e.mock.On("Size")}
}

func (_c *TargetVBOpaque_Size_Call) Run(run func()) *TargetVBOpaque_Size_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetVBOpaque_Size_Call) Return(_a0 int) *TargetVBOpaque_Size_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetVBOpaque_Size_Call) RunAndReturn(run func() int) *TargetVBOpaque_Size_Call {
	_c.Call.Return(run)
	return _c
}

// Value provides a mock function with given fields:
func (_m *TargetVBOpaque) Value() interface{} {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Value")
	}

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// TargetVBOpaque_Value_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Value'
type TargetVBOpaque_Value_Call struct {
	*mock.Call
}

// Value is a helper method to define mock.On call
func (_e *TargetVBOpaque_Expecter) Value() *TargetVBOpaque_Value_Call {
	return &TargetVBOpaque_Value_Call{Call: _e.mock.On("Value")}
}

func (_c *TargetVBOpaque_Value_Call) Run(run func()) *TargetVBOpaque_Value_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetVBOpaque_Value_Call) Return(_a0 interface{}) *TargetVBOpaque_Value_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetVBOpaque_Value_Call) RunAndReturn(run func() interface{}) *TargetVBOpaque_Value_Call {
	_c.Call.Return(run)
	return _c
}

// NewTargetVBOpaque creates a new instance of TargetVBOpaque. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTargetVBOpaque(t interface {
	mock.TestingT
	Cleanup(func())
}) *TargetVBOpaque {
	mock := &TargetVBOpaque{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
