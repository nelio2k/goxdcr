// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	parts "github.com/couchbase/goxdcr/parts"
	mock "github.com/stretchr/testify/mock"
)

// stateCheckFunc is an autogenerated mock type for the stateCheckFunc type
type stateCheckFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: state
func (_m *stateCheckFunc) Execute(state parts.DcpStreamState) bool {
	ret := _m.Called(state)

	var r0 bool
	if rf, ok := ret.Get(0).(func(parts.DcpStreamState) bool); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
