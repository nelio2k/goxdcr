// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	service_def "github.com/couchbase/goxdcr/service_def"
	mock "github.com/stretchr/testify/mock"
)

// AuditEventIface is an autogenerated mock type for the AuditEventIface type
type AuditEventIface struct {
	mock.Mock
}

// Clone provides a mock function with given fields:
func (_m *AuditEventIface) Clone() service_def.AuditEventIface {
	ret := _m.Called()

	var r0 service_def.AuditEventIface
	if rf, ok := ret.Get(0).(func() service_def.AuditEventIface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service_def.AuditEventIface)
		}
	}

	return r0
}

// Redact provides a mock function with given fields:
func (_m *AuditEventIface) Redact() service_def.AuditEventIface {
	ret := _m.Called()

	var r0 service_def.AuditEventIface
	if rf, ok := ret.Get(0).(func() service_def.AuditEventIface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(service_def.AuditEventIface)
		}
	}

	return r0
}

type mockConstructorTestingTNewAuditEventIface interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuditEventIface creates a new instance of AuditEventIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuditEventIface(t mockConstructorTestingTNewAuditEventIface) *AuditEventIface {
	mock := &AuditEventIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
