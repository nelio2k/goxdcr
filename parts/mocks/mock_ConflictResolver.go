// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/v8/base"
	crMeta "github.com/couchbase/goxdcr/v8/crMeta"

	gomemcached "github.com/couchbase/gomemcached"

	hlv "github.com/couchbase/goxdcr/v8/hlv"

	log "github.com/couchbase/goxdcr/v8/log"

	mock "github.com/stretchr/testify/mock"
)

// ConflictResolver is an autogenerated mock type for the ConflictResolver type
type ConflictResolver struct {
	mock.Mock
}

type ConflictResolver_Expecter struct {
	mock *mock.Mock
}

func (_m *ConflictResolver) EXPECT() *ConflictResolver_Expecter {
	return &ConflictResolver_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: req, resp, specs, sourceId, targetId, logConflict, xattrEnabled, uncompressFunc, logger
func (_m *ConflictResolver) Execute(req *base.WrappedMCRequest, resp *gomemcached.MCResponse, specs []base.SubdocLookupPathSpec, sourceId hlv.DocumentSourceId, targetId hlv.DocumentSourceId, logConflict bool, xattrEnabled bool, uncompressFunc base.UncompressFunc, logger *log.CommonLogger) (crMeta.ConflictDetectionResult, crMeta.ConflictResolutionResult, error) {
	ret := _m.Called(req, resp, specs, sourceId, targetId, logConflict, xattrEnabled, uncompressFunc, logger)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 crMeta.ConflictDetectionResult
	var r1 crMeta.ConflictResolutionResult
	var r2 error
	if rf, ok := ret.Get(0).(func(*base.WrappedMCRequest, *gomemcached.MCResponse, []base.SubdocLookupPathSpec, hlv.DocumentSourceId, hlv.DocumentSourceId, bool, bool, base.UncompressFunc, *log.CommonLogger) (crMeta.ConflictDetectionResult, crMeta.ConflictResolutionResult, error)); ok {
		return rf(req, resp, specs, sourceId, targetId, logConflict, xattrEnabled, uncompressFunc, logger)
	}
	if rf, ok := ret.Get(0).(func(*base.WrappedMCRequest, *gomemcached.MCResponse, []base.SubdocLookupPathSpec, hlv.DocumentSourceId, hlv.DocumentSourceId, bool, bool, base.UncompressFunc, *log.CommonLogger) crMeta.ConflictDetectionResult); ok {
		r0 = rf(req, resp, specs, sourceId, targetId, logConflict, xattrEnabled, uncompressFunc, logger)
	} else {
		r0 = ret.Get(0).(crMeta.ConflictDetectionResult)
	}

	if rf, ok := ret.Get(1).(func(*base.WrappedMCRequest, *gomemcached.MCResponse, []base.SubdocLookupPathSpec, hlv.DocumentSourceId, hlv.DocumentSourceId, bool, bool, base.UncompressFunc, *log.CommonLogger) crMeta.ConflictResolutionResult); ok {
		r1 = rf(req, resp, specs, sourceId, targetId, logConflict, xattrEnabled, uncompressFunc, logger)
	} else {
		r1 = ret.Get(1).(crMeta.ConflictResolutionResult)
	}

	if rf, ok := ret.Get(2).(func(*base.WrappedMCRequest, *gomemcached.MCResponse, []base.SubdocLookupPathSpec, hlv.DocumentSourceId, hlv.DocumentSourceId, bool, bool, base.UncompressFunc, *log.CommonLogger) error); ok {
		r2 = rf(req, resp, specs, sourceId, targetId, logConflict, xattrEnabled, uncompressFunc, logger)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ConflictResolver_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type ConflictResolver_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - req *base.WrappedMCRequest
//   - resp *gomemcached.MCResponse
//   - specs []base.SubdocLookupPathSpec
//   - sourceId hlv.DocumentSourceId
//   - targetId hlv.DocumentSourceId
//   - logConflict bool
//   - xattrEnabled bool
//   - uncompressFunc base.UncompressFunc
//   - logger *log.CommonLogger
func (_e *ConflictResolver_Expecter) Execute(req interface{}, resp interface{}, specs interface{}, sourceId interface{}, targetId interface{}, logConflict interface{}, xattrEnabled interface{}, uncompressFunc interface{}, logger interface{}) *ConflictResolver_Execute_Call {
	return &ConflictResolver_Execute_Call{Call: _e.mock.On("Execute", req, resp, specs, sourceId, targetId, logConflict, xattrEnabled, uncompressFunc, logger)}
}

func (_c *ConflictResolver_Execute_Call) Run(run func(req *base.WrappedMCRequest, resp *gomemcached.MCResponse, specs []base.SubdocLookupPathSpec, sourceId hlv.DocumentSourceId, targetId hlv.DocumentSourceId, logConflict bool, xattrEnabled bool, uncompressFunc base.UncompressFunc, logger *log.CommonLogger)) *ConflictResolver_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*base.WrappedMCRequest), args[1].(*gomemcached.MCResponse), args[2].([]base.SubdocLookupPathSpec), args[3].(hlv.DocumentSourceId), args[4].(hlv.DocumentSourceId), args[5].(bool), args[6].(bool), args[7].(base.UncompressFunc), args[8].(*log.CommonLogger))
	})
	return _c
}

func (_c *ConflictResolver_Execute_Call) Return(_a0 crMeta.ConflictDetectionResult, _a1 crMeta.ConflictResolutionResult, _a2 error) *ConflictResolver_Execute_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *ConflictResolver_Execute_Call) RunAndReturn(run func(*base.WrappedMCRequest, *gomemcached.MCResponse, []base.SubdocLookupPathSpec, hlv.DocumentSourceId, hlv.DocumentSourceId, bool, bool, base.UncompressFunc, *log.CommonLogger) (crMeta.ConflictDetectionResult, crMeta.ConflictResolutionResult, error)) *ConflictResolver_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewConflictResolver creates a new instance of ConflictResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConflictResolver(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConflictResolver {
	mock := &ConflictResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
