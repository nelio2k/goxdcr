// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	service_def "github.com/couchbase/goxdcr/service_def"
	mock "github.com/stretchr/testify/mock"
)

// SecChangeCallback is an autogenerated mock type for the SecChangeCallback type
type SecChangeCallback struct {
	mock.Mock
}

// Execute provides a mock function with given fields: old, new
func (_m *SecChangeCallback) Execute(old service_def.EncryptionSettingIface, new service_def.EncryptionSettingIface) {
	_m.Called(old, new)
}