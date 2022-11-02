// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UIWarnings is an autogenerated mock type for the UIWarnings type
type UIWarnings struct {
	mock.Mock
}

// AddWarning provides a mock function with given fields: key, val
func (_m *UIWarnings) AddWarning(key string, val string) {
	_m.Called(key, val)
}

// AppendGeneric provides a mock function with given fields: warning
func (_m *UIWarnings) AppendGeneric(warning string) {
	_m.Called(warning)
}

// GetFieldWarningsOnly provides a mock function with given fields:
func (_m *UIWarnings) GetFieldWarningsOnly() map[string]interface{} {
	ret := _m.Called()

	var r0 map[string]interface{}
	if rf, ok := ret.Get(0).(func() map[string]interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	return r0
}

// GetSuccessfulWarningStrings provides a mock function with given fields:
func (_m *UIWarnings) GetSuccessfulWarningStrings() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Len provides a mock function with given fields:
func (_m *UIWarnings) Len() int {
	ret := _m.Called()

	var r0 int
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *UIWarnings) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

type mockConstructorTestingTNewUIWarnings interface {
	mock.TestingT
	Cleanup(func())
}

// NewUIWarnings creates a new instance of UIWarnings. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUIWarnings(t mockConstructorTestingTNewUIWarnings) *UIWarnings {
	mock := &UIWarnings{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
