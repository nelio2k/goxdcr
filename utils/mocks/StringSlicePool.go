// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StringSlicePool is an autogenerated mock type for the StringSlicePool type
type StringSlicePool struct {
	mock.Mock
}

// Get provides a mock function with given fields:
func (_m *StringSlicePool) Get() *[]string {
	ret := _m.Called()

	var r0 *[]string
	if rf, ok := ret.Get(0).(func() *[]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]string)
		}
	}

	return r0
}

// Put provides a mock function with given fields: _a0
func (_m *StringSlicePool) Put(_a0 *[]string) {
	_m.Called(_a0)
}
