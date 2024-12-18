// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/v8/base"
	mock "github.com/stretchr/testify/mock"
)

// SortedSeqnoListWithLock is an autogenerated mock type for the SortedSeqnoListWithLock type
type SortedSeqnoListWithLock struct {
	mock.Mock
}

type SortedSeqnoListWithLock_Expecter struct {
	mock *mock.Mock
}

func (_m *SortedSeqnoListWithLock) EXPECT() *SortedSeqnoListWithLock_Expecter {
	return &SortedSeqnoListWithLock_Expecter{mock: &_m.Mock}
}

// AppendSeqnoGlobal provides a mock function with given fields: _a0, _a1, _a2
func (_m *SortedSeqnoListWithLock) AppendSeqnoGlobal(_a0 uint16, _a1 string, _a2 uint64) {
	_m.Called(_a0, _a1, _a2)
}

// SortedSeqnoListWithLock_AppendSeqnoGlobal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppendSeqnoGlobal'
type SortedSeqnoListWithLock_AppendSeqnoGlobal_Call struct {
	*mock.Call
}

// AppendSeqnoGlobal is a helper method to define mock.On call
//   - _a0 uint16
//   - _a1 string
//   - _a2 uint64
func (_e *SortedSeqnoListWithLock_Expecter) AppendSeqnoGlobal(_a0 interface{}, _a1 interface{}, _a2 interface{}) *SortedSeqnoListWithLock_AppendSeqnoGlobal_Call {
	return &SortedSeqnoListWithLock_AppendSeqnoGlobal_Call{Call: _e.mock.On("AppendSeqnoGlobal", _a0, _a1, _a2)}
}

func (_c *SortedSeqnoListWithLock_AppendSeqnoGlobal_Call) Run(run func(_a0 uint16, _a1 string, _a2 uint64)) *SortedSeqnoListWithLock_AppendSeqnoGlobal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint16), args[1].(string), args[2].(uint64))
	})
	return _c
}

func (_c *SortedSeqnoListWithLock_AppendSeqnoGlobal_Call) Return() *SortedSeqnoListWithLock_AppendSeqnoGlobal_Call {
	_c.Call.Return()
	return _c
}

func (_c *SortedSeqnoListWithLock_AppendSeqnoGlobal_Call) RunAndReturn(run func(uint16, string, uint64)) *SortedSeqnoListWithLock_AppendSeqnoGlobal_Call {
	_c.Call.Return(run)
	return _c
}

// AppendSeqnoTraditional provides a mock function with given fields: _a0, _a1
func (_m *SortedSeqnoListWithLock) AppendSeqnoTraditional(_a0 string, _a1 uint64) {
	_m.Called(_a0, _a1)
}

// SortedSeqnoListWithLock_AppendSeqnoTraditional_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppendSeqnoTraditional'
type SortedSeqnoListWithLock_AppendSeqnoTraditional_Call struct {
	*mock.Call
}

// AppendSeqnoTraditional is a helper method to define mock.On call
//   - _a0 string
//   - _a1 uint64
func (_e *SortedSeqnoListWithLock_Expecter) AppendSeqnoTraditional(_a0 interface{}, _a1 interface{}) *SortedSeqnoListWithLock_AppendSeqnoTraditional_Call {
	return &SortedSeqnoListWithLock_AppendSeqnoTraditional_Call{Call: _e.mock.On("AppendSeqnoTraditional", _a0, _a1)}
}

func (_c *SortedSeqnoListWithLock_AppendSeqnoTraditional_Call) Run(run func(_a0 string, _a1 uint64)) *SortedSeqnoListWithLock_AppendSeqnoTraditional_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(uint64))
	})
	return _c
}

func (_c *SortedSeqnoListWithLock_AppendSeqnoTraditional_Call) Return() *SortedSeqnoListWithLock_AppendSeqnoTraditional_Call {
	_c.Call.Return()
	return _c
}

func (_c *SortedSeqnoListWithLock_AppendSeqnoTraditional_Call) RunAndReturn(run func(string, uint64)) *SortedSeqnoListWithLock_AppendSeqnoTraditional_Call {
	_c.Call.Return(run)
	return _c
}

// GetterGlobal provides a mock function with given fields: _a0, _a1
func (_m *SortedSeqnoListWithLock) GetterGlobal(_a0 uint16, _a1 string) *base.SortedSeqnoListWithLock {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetterGlobal")
	}

	var r0 *base.SortedSeqnoListWithLock
	if rf, ok := ret.Get(0).(func(uint16, string) *base.SortedSeqnoListWithLock); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*base.SortedSeqnoListWithLock)
		}
	}

	return r0
}

// SortedSeqnoListWithLock_GetterGlobal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetterGlobal'
type SortedSeqnoListWithLock_GetterGlobal_Call struct {
	*mock.Call
}

// GetterGlobal is a helper method to define mock.On call
//   - _a0 uint16
//   - _a1 string
func (_e *SortedSeqnoListWithLock_Expecter) GetterGlobal(_a0 interface{}, _a1 interface{}) *SortedSeqnoListWithLock_GetterGlobal_Call {
	return &SortedSeqnoListWithLock_GetterGlobal_Call{Call: _e.mock.On("GetterGlobal", _a0, _a1)}
}

func (_c *SortedSeqnoListWithLock_GetterGlobal_Call) Run(run func(_a0 uint16, _a1 string)) *SortedSeqnoListWithLock_GetterGlobal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint16), args[1].(string))
	})
	return _c
}

func (_c *SortedSeqnoListWithLock_GetterGlobal_Call) Return(_a0 *base.SortedSeqnoListWithLock) *SortedSeqnoListWithLock_GetterGlobal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SortedSeqnoListWithLock_GetterGlobal_Call) RunAndReturn(run func(uint16, string) *base.SortedSeqnoListWithLock) *SortedSeqnoListWithLock_GetterGlobal_Call {
	_c.Call.Return(run)
	return _c
}

// GetterTraditional provides a mock function with given fields: _a0
func (_m *SortedSeqnoListWithLock) GetterTraditional(_a0 string) *base.SortedSeqnoListWithLock {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetterTraditional")
	}

	var r0 *base.SortedSeqnoListWithLock
	if rf, ok := ret.Get(0).(func(string) *base.SortedSeqnoListWithLock); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*base.SortedSeqnoListWithLock)
		}
	}

	return r0
}

// SortedSeqnoListWithLock_GetterTraditional_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetterTraditional'
type SortedSeqnoListWithLock_GetterTraditional_Call struct {
	*mock.Call
}

// GetterTraditional is a helper method to define mock.On call
//   - _a0 string
func (_e *SortedSeqnoListWithLock_Expecter) GetterTraditional(_a0 interface{}) *SortedSeqnoListWithLock_GetterTraditional_Call {
	return &SortedSeqnoListWithLock_GetterTraditional_Call{Call: _e.mock.On("GetterTraditional", _a0)}
}

func (_c *SortedSeqnoListWithLock_GetterTraditional_Call) Run(run func(_a0 string)) *SortedSeqnoListWithLock_GetterTraditional_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *SortedSeqnoListWithLock_GetterTraditional_Call) Return(_a0 *base.SortedSeqnoListWithLock) *SortedSeqnoListWithLock_GetterTraditional_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SortedSeqnoListWithLock_GetterTraditional_Call) RunAndReturn(run func(string) *base.SortedSeqnoListWithLock) *SortedSeqnoListWithLock_GetterTraditional_Call {
	_c.Call.Return(run)
	return _c
}

// IsTraditional provides a mock function with given fields:
func (_m *SortedSeqnoListWithLock) IsTraditional() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsTraditional")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SortedSeqnoListWithLock_IsTraditional_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsTraditional'
type SortedSeqnoListWithLock_IsTraditional_Call struct {
	*mock.Call
}

// IsTraditional is a helper method to define mock.On call
func (_e *SortedSeqnoListWithLock_Expecter) IsTraditional() *SortedSeqnoListWithLock_IsTraditional_Call {
	return &SortedSeqnoListWithLock_IsTraditional_Call{Call: _e.mock.On("IsTraditional")}
}

func (_c *SortedSeqnoListWithLock_IsTraditional_Call) Run(run func()) *SortedSeqnoListWithLock_IsTraditional_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *SortedSeqnoListWithLock_IsTraditional_Call) Return(_a0 bool) *SortedSeqnoListWithLock_IsTraditional_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SortedSeqnoListWithLock_IsTraditional_Call) RunAndReturn(run func() bool) *SortedSeqnoListWithLock_IsTraditional_Call {
	_c.Call.Return(run)
	return _c
}

// NewSortedSeqnoListWithLock creates a new instance of SortedSeqnoListWithLock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSortedSeqnoListWithLock(t interface {
	mock.TestingT
	Cleanup(func())
}) *SortedSeqnoListWithLock {
	mock := &SortedSeqnoListWithLock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
