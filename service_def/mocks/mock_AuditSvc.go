// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	service_def "github.com/couchbase/goxdcr/service_def"
	mock "github.com/stretchr/testify/mock"
)

// AuditSvc is an autogenerated mock type for the AuditSvc type
type AuditSvc struct {
	mock.Mock
}

type AuditSvc_Expecter struct {
	mock *mock.Mock
}

func (_m *AuditSvc) EXPECT() *AuditSvc_Expecter {
	return &AuditSvc_Expecter{mock: &_m.Mock}
}

// Write provides a mock function with given fields: eventId, event
func (_m *AuditSvc) Write(eventId uint32, event service_def.AuditEventIface) error {
	ret := _m.Called(eventId, event)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint32, service_def.AuditEventIface) error); ok {
		r0 = rf(eventId, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AuditSvc_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type AuditSvc_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - eventId uint32
//   - event service_def.AuditEventIface
func (_e *AuditSvc_Expecter) Write(eventId interface{}, event interface{}) *AuditSvc_Write_Call {
	return &AuditSvc_Write_Call{Call: _e.mock.On("Write", eventId, event)}
}

func (_c *AuditSvc_Write_Call) Run(run func(eventId uint32, event service_def.AuditEventIface)) *AuditSvc_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint32), args[1].(service_def.AuditEventIface))
	})
	return _c
}

func (_c *AuditSvc_Write_Call) Return(_a0 error) *AuditSvc_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AuditSvc_Write_Call) RunAndReturn(run func(uint32, service_def.AuditEventIface) error) *AuditSvc_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuditSvc creates a new instance of AuditSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuditSvc(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuditSvc {
	mock := &AuditSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}