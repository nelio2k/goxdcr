// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/v8/base"
	metadata "github.com/couchbase/goxdcr/v8/metadata"

	mock "github.com/stretchr/testify/mock"
)

// PipelineMgrBackfillIface is an autogenerated mock type for the PipelineMgrBackfillIface type
type PipelineMgrBackfillIface struct {
	mock.Mock
}

type PipelineMgrBackfillIface_Expecter struct {
	mock *mock.Mock
}

func (_m *PipelineMgrBackfillIface) EXPECT() *PipelineMgrBackfillIface_Expecter {
	return &PipelineMgrBackfillIface_Expecter{mock: &_m.Mock}
}

// BackfillMappingStatusUpdate provides a mock function with given fields: topic, diffPair, srcManifestDelta
func (_m *PipelineMgrBackfillIface) BackfillMappingStatusUpdate(topic string, diffPair *metadata.CollectionNamespaceMappingsDiffPair, srcManifestDelta []*metadata.CollectionsManifest) error {
	ret := _m.Called(topic, diffPair, srcManifestDelta)

	if len(ret) == 0 {
		panic("no return value specified for BackfillMappingStatusUpdate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *metadata.CollectionNamespaceMappingsDiffPair, []*metadata.CollectionsManifest) error); ok {
		r0 = rf(topic, diffPair, srcManifestDelta)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BackfillMappingStatusUpdate'
type PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call struct {
	*mock.Call
}

// BackfillMappingStatusUpdate is a helper method to define mock.On call
//   - topic string
//   - diffPair *metadata.CollectionNamespaceMappingsDiffPair
//   - srcManifestDelta []*metadata.CollectionsManifest
func (_e *PipelineMgrBackfillIface_Expecter) BackfillMappingStatusUpdate(topic interface{}, diffPair interface{}, srcManifestDelta interface{}) *PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call {
	return &PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call{Call: _e.mock.On("BackfillMappingStatusUpdate", topic, diffPair, srcManifestDelta)}
}

func (_c *PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call) Run(run func(topic string, diffPair *metadata.CollectionNamespaceMappingsDiffPair, srcManifestDelta []*metadata.CollectionsManifest)) *PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*metadata.CollectionNamespaceMappingsDiffPair), args[2].([]*metadata.CollectionsManifest))
	})
	return _c
}

func (_c *PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call) Return(_a0 error) *PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call) RunAndReturn(run func(string, *metadata.CollectionNamespaceMappingsDiffPair, []*metadata.CollectionsManifest) error) *PipelineMgrBackfillIface_BackfillMappingStatusUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// CleanupBackfillCkpts provides a mock function with given fields: topic
func (_m *PipelineMgrBackfillIface) CleanupBackfillCkpts(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for CleanupBackfillCkpts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineMgrBackfillIface_CleanupBackfillCkpts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CleanupBackfillCkpts'
type PipelineMgrBackfillIface_CleanupBackfillCkpts_Call struct {
	*mock.Call
}

// CleanupBackfillCkpts is a helper method to define mock.On call
//   - topic string
func (_e *PipelineMgrBackfillIface_Expecter) CleanupBackfillCkpts(topic interface{}) *PipelineMgrBackfillIface_CleanupBackfillCkpts_Call {
	return &PipelineMgrBackfillIface_CleanupBackfillCkpts_Call{Call: _e.mock.On("CleanupBackfillCkpts", topic)}
}

func (_c *PipelineMgrBackfillIface_CleanupBackfillCkpts_Call) Run(run func(topic string)) *PipelineMgrBackfillIface_CleanupBackfillCkpts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineMgrBackfillIface_CleanupBackfillCkpts_Call) Return(_a0 error) *PipelineMgrBackfillIface_CleanupBackfillCkpts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineMgrBackfillIface_CleanupBackfillCkpts_Call) RunAndReturn(run func(string) error) *PipelineMgrBackfillIface_CleanupBackfillCkpts_Call {
	_c.Call.Return(run)
	return _c
}

// GetMainPipelineThroughSeqnos provides a mock function with given fields: topic
func (_m *PipelineMgrBackfillIface) GetMainPipelineThroughSeqnos(topic string) (map[uint16]uint64, error) {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for GetMainPipelineThroughSeqnos")
	}

	var r0 map[uint16]uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (map[uint16]uint64, error)); ok {
		return rf(topic)
	}
	if rf, ok := ret.Get(0).(func(string) map[uint16]uint64); ok {
		r0 = rf(topic)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uint16]uint64)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(topic)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMainPipelineThroughSeqnos'
type PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call struct {
	*mock.Call
}

// GetMainPipelineThroughSeqnos is a helper method to define mock.On call
//   - topic string
func (_e *PipelineMgrBackfillIface_Expecter) GetMainPipelineThroughSeqnos(topic interface{}) *PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call {
	return &PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call{Call: _e.mock.On("GetMainPipelineThroughSeqnos", topic)}
}

func (_c *PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call) Run(run func(topic string)) *PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call) Return(_a0 map[uint16]uint64, _a1 error) *PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call) RunAndReturn(run func(string) (map[uint16]uint64, error)) *PipelineMgrBackfillIface_GetMainPipelineThroughSeqnos_Call {
	_c.Call.Return(run)
	return _c
}

// HaltBackfill provides a mock function with given fields: topic
func (_m *PipelineMgrBackfillIface) HaltBackfill(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for HaltBackfill")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineMgrBackfillIface_HaltBackfill_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HaltBackfill'
type PipelineMgrBackfillIface_HaltBackfill_Call struct {
	*mock.Call
}

// HaltBackfill is a helper method to define mock.On call
//   - topic string
func (_e *PipelineMgrBackfillIface_Expecter) HaltBackfill(topic interface{}) *PipelineMgrBackfillIface_HaltBackfill_Call {
	return &PipelineMgrBackfillIface_HaltBackfill_Call{Call: _e.mock.On("HaltBackfill", topic)}
}

func (_c *PipelineMgrBackfillIface_HaltBackfill_Call) Run(run func(topic string)) *PipelineMgrBackfillIface_HaltBackfill_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineMgrBackfillIface_HaltBackfill_Call) Return(_a0 error) *PipelineMgrBackfillIface_HaltBackfill_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineMgrBackfillIface_HaltBackfill_Call) RunAndReturn(run func(string) error) *PipelineMgrBackfillIface_HaltBackfill_Call {
	_c.Call.Return(run)
	return _c
}

// HaltBackfillWithCb provides a mock function with given fields: topic, callback, errCb, skipCkpt
func (_m *PipelineMgrBackfillIface) HaltBackfillWithCb(topic string, callback base.StoppedPipelineCallback, errCb base.StoppedPipelineErrCallback, skipCkpt bool) error {
	ret := _m.Called(topic, callback, errCb, skipCkpt)

	if len(ret) == 0 {
		panic("no return value specified for HaltBackfillWithCb")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, base.StoppedPipelineCallback, base.StoppedPipelineErrCallback, bool) error); ok {
		r0 = rf(topic, callback, errCb, skipCkpt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineMgrBackfillIface_HaltBackfillWithCb_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HaltBackfillWithCb'
type PipelineMgrBackfillIface_HaltBackfillWithCb_Call struct {
	*mock.Call
}

// HaltBackfillWithCb is a helper method to define mock.On call
//   - topic string
//   - callback base.StoppedPipelineCallback
//   - errCb base.StoppedPipelineErrCallback
//   - skipCkpt bool
func (_e *PipelineMgrBackfillIface_Expecter) HaltBackfillWithCb(topic interface{}, callback interface{}, errCb interface{}, skipCkpt interface{}) *PipelineMgrBackfillIface_HaltBackfillWithCb_Call {
	return &PipelineMgrBackfillIface_HaltBackfillWithCb_Call{Call: _e.mock.On("HaltBackfillWithCb", topic, callback, errCb, skipCkpt)}
}

func (_c *PipelineMgrBackfillIface_HaltBackfillWithCb_Call) Run(run func(topic string, callback base.StoppedPipelineCallback, errCb base.StoppedPipelineErrCallback, skipCkpt bool)) *PipelineMgrBackfillIface_HaltBackfillWithCb_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(base.StoppedPipelineCallback), args[2].(base.StoppedPipelineErrCallback), args[3].(bool))
	})
	return _c
}

func (_c *PipelineMgrBackfillIface_HaltBackfillWithCb_Call) Return(_a0 error) *PipelineMgrBackfillIface_HaltBackfillWithCb_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineMgrBackfillIface_HaltBackfillWithCb_Call) RunAndReturn(run func(string, base.StoppedPipelineCallback, base.StoppedPipelineErrCallback, bool) error) *PipelineMgrBackfillIface_HaltBackfillWithCb_Call {
	_c.Call.Return(run)
	return _c
}

// ReInitStreams provides a mock function with given fields: pipelineName
func (_m *PipelineMgrBackfillIface) ReInitStreams(pipelineName string) error {
	ret := _m.Called(pipelineName)

	if len(ret) == 0 {
		panic("no return value specified for ReInitStreams")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(pipelineName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineMgrBackfillIface_ReInitStreams_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReInitStreams'
type PipelineMgrBackfillIface_ReInitStreams_Call struct {
	*mock.Call
}

// ReInitStreams is a helper method to define mock.On call
//   - pipelineName string
func (_e *PipelineMgrBackfillIface_Expecter) ReInitStreams(pipelineName interface{}) *PipelineMgrBackfillIface_ReInitStreams_Call {
	return &PipelineMgrBackfillIface_ReInitStreams_Call{Call: _e.mock.On("ReInitStreams", pipelineName)}
}

func (_c *PipelineMgrBackfillIface_ReInitStreams_Call) Run(run func(pipelineName string)) *PipelineMgrBackfillIface_ReInitStreams_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineMgrBackfillIface_ReInitStreams_Call) Return(_a0 error) *PipelineMgrBackfillIface_ReInitStreams_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineMgrBackfillIface_ReInitStreams_Call) RunAndReturn(run func(string) error) *PipelineMgrBackfillIface_ReInitStreams_Call {
	_c.Call.Return(run)
	return _c
}

// RequestBackfill provides a mock function with given fields: topic
func (_m *PipelineMgrBackfillIface) RequestBackfill(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for RequestBackfill")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineMgrBackfillIface_RequestBackfill_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestBackfill'
type PipelineMgrBackfillIface_RequestBackfill_Call struct {
	*mock.Call
}

// RequestBackfill is a helper method to define mock.On call
//   - topic string
func (_e *PipelineMgrBackfillIface_Expecter) RequestBackfill(topic interface{}) *PipelineMgrBackfillIface_RequestBackfill_Call {
	return &PipelineMgrBackfillIface_RequestBackfill_Call{Call: _e.mock.On("RequestBackfill", topic)}
}

func (_c *PipelineMgrBackfillIface_RequestBackfill_Call) Run(run func(topic string)) *PipelineMgrBackfillIface_RequestBackfill_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineMgrBackfillIface_RequestBackfill_Call) Return(_a0 error) *PipelineMgrBackfillIface_RequestBackfill_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineMgrBackfillIface_RequestBackfill_Call) RunAndReturn(run func(string) error) *PipelineMgrBackfillIface_RequestBackfill_Call {
	_c.Call.Return(run)
	return _c
}

// WaitForMainPipelineCkptMgrToStop provides a mock function with given fields: topic, internalID
func (_m *PipelineMgrBackfillIface) WaitForMainPipelineCkptMgrToStop(topic string, internalID string) {
	_m.Called(topic, internalID)
}

// PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WaitForMainPipelineCkptMgrToStop'
type PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call struct {
	*mock.Call
}

// WaitForMainPipelineCkptMgrToStop is a helper method to define mock.On call
//   - topic string
//   - internalID string
func (_e *PipelineMgrBackfillIface_Expecter) WaitForMainPipelineCkptMgrToStop(topic interface{}, internalID interface{}) *PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call {
	return &PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call{Call: _e.mock.On("WaitForMainPipelineCkptMgrToStop", topic, internalID)}
}

func (_c *PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call) Run(run func(topic string, internalID string)) *PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call) Return() *PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call {
	_c.Call.Return()
	return _c
}

func (_c *PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call) RunAndReturn(run func(string, string)) *PipelineMgrBackfillIface_WaitForMainPipelineCkptMgrToStop_Call {
	_c.Call.Return(run)
	return _c
}

// NewPipelineMgrBackfillIface creates a new instance of PipelineMgrBackfillIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPipelineMgrBackfillIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *PipelineMgrBackfillIface {
	mock := &PipelineMgrBackfillIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
