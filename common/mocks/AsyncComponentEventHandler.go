// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	mock "github.com/stretchr/testify/mock"
)

// AsyncComponentEventHandler is an autogenerated mock type for the AsyncComponentEventHandler type
type AsyncComponentEventHandler struct {
	mock.Mock
}

// Id provides a mock function with given fields:
func (_m *AsyncComponentEventHandler) Id() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ProcessEvent provides a mock function with given fields: event
func (_m *AsyncComponentEventHandler) ProcessEvent(event *common.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(*common.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
