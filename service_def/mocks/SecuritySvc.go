// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	service_def "github.com/couchbase/goxdcr/service_def"
	mock "github.com/stretchr/testify/mock"

	x509 "crypto/x509"
)

// SecuritySvc is an autogenerated mock type for the SecuritySvc type
type SecuritySvc struct {
	mock.Mock
}

// EncryptData provides a mock function with given fields:
func (_m *SecuritySvc) EncryptData() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetCACertificates provides a mock function with given fields:
func (_m *SecuritySvc) GetCACertificates() []byte {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	return r0
}

// GetCaPool provides a mock function with given fields:
func (_m *SecuritySvc) GetCaPool() *x509.CertPool {
	ret := _m.Called()

	var r0 *x509.CertPool
	if rf, ok := ret.Get(0).(func() *x509.CertPool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*x509.CertPool)
		}
	}

	return r0
}

// IsClusterEncryptionLevelStrict provides a mock function with given fields:
func (_m *SecuritySvc) IsClusterEncryptionLevelStrict() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SetEncryptionLevelChangeCallback provides a mock function with given fields: key, callback
func (_m *SecuritySvc) SetEncryptionLevelChangeCallback(key string, callback service_def.SecChangeCallback) {
	_m.Called(key, callback)
}

// Start provides a mock function with given fields:
func (_m *SecuritySvc) Start() {
	_m.Called()
}
