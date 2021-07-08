// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Notification is an autogenerated mock type for the Notification type
type Notification struct {
	mock.Mock
}

// CloneRO provides a mock function with given fields:
func (_m *Notification) CloneRO() interface{} {
	ret := _m.Called()

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