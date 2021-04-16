// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	common "github.com/couchbase/goxdcr/common"

	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"
)

// PipelineRuntimeContext is an autogenerated mock type for the PipelineRuntimeContext type
type PipelineRuntimeContext struct {
	mock.Mock
}

// Pipeline provides a mock function with given fields:
func (_m *PipelineRuntimeContext) Pipeline() common.Pipeline {
	ret := _m.Called()

	var r0 common.Pipeline
	if rf, ok := ret.Get(0).(func() common.Pipeline); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Pipeline)
		}
	}

	return r0
}

// RegisterService provides a mock function with given fields: svc_name, svc
func (_m *PipelineRuntimeContext) RegisterService(svc_name string, svc common.PipelineService) error {
	ret := _m.Called(svc_name, svc)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, common.PipelineService) error); ok {
		r0 = rf(svc_name, svc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Service provides a mock function with given fields: svc_name
func (_m *PipelineRuntimeContext) Service(svc_name string) common.PipelineService {
	ret := _m.Called(svc_name)

	var r0 common.PipelineService
	if rf, ok := ret.Get(0).(func(string) common.PipelineService); ok {
		r0 = rf(svc_name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.PipelineService)
		}
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *PipelineRuntimeContext) Start(_a0 metadata.ReplicationSettingsMap) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Stop provides a mock function with given fields:
func (_m *PipelineRuntimeContext) Stop() base.ErrorMap {
	ret := _m.Called()

	var r0 base.ErrorMap
	if rf, ok := ret.Get(0).(func() base.ErrorMap); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(base.ErrorMap)
		}
	}

	return r0
}

// UnregisterService provides a mock function with given fields: srv_name
func (_m *PipelineRuntimeContext) UnregisterService(srv_name string) error {
	ret := _m.Called(srv_name)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(srv_name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSettings provides a mock function with given fields: settings
func (_m *PipelineRuntimeContext) UpdateSettings(settings metadata.ReplicationSettingsMap) error {
	ret := _m.Called(settings)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) error); ok {
		r0 = rf(settings)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
