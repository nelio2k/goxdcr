// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/v8/base"
	metadata "github.com/couchbase/goxdcr/v8/metadata"

	mock "github.com/stretchr/testify/mock"

	pipeline "github.com/couchbase/goxdcr/v8/pipeline"
)

// PipelineOpSerializerIface is an autogenerated mock type for the PipelineOpSerializerIface type
type PipelineOpSerializerIface struct {
	mock.Mock
}

type PipelineOpSerializerIface_Expecter struct {
	mock *mock.Mock
}

func (_m *PipelineOpSerializerIface) EXPECT() *PipelineOpSerializerIface_Expecter {
	return &PipelineOpSerializerIface_Expecter{mock: &_m.Mock}
}

// BackfillMappingStatusUpdate provides a mock function with given fields: topic, diffPair, srcManifestsDelta, latestTgtManifestId
func (_m *PipelineOpSerializerIface) BackfillMappingStatusUpdate(topic string, diffPair *metadata.CollectionNamespaceMappingsDiffPair, srcManifestsDelta []*metadata.CollectionsManifest, latestTgtManifestId uint64) error {
	ret := _m.Called(topic, diffPair, srcManifestsDelta, latestTgtManifestId)

	if len(ret) == 0 {
		panic("no return value specified for BackfillMappingStatusUpdate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *metadata.CollectionNamespaceMappingsDiffPair, []*metadata.CollectionsManifest, uint64) error); ok {
		r0 = rf(topic, diffPair, srcManifestsDelta, latestTgtManifestId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BackfillMappingStatusUpdate'
type PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call struct {
	*mock.Call
}

// BackfillMappingStatusUpdate is a helper method to define mock.On call
//   - topic string
//   - diffPair *metadata.CollectionNamespaceMappingsDiffPair
//   - srcManifestsDelta []*metadata.CollectionsManifest
//   - latestTgtManifestId uint64
func (_e *PipelineOpSerializerIface_Expecter) BackfillMappingStatusUpdate(topic interface{}, diffPair interface{}, srcManifestsDelta interface{}, latestTgtManifestId interface{}) *PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call {
	return &PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call{Call: _e.mock.On("BackfillMappingStatusUpdate", topic, diffPair, srcManifestsDelta, latestTgtManifestId)}
}

func (_c *PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call) Run(run func(topic string, diffPair *metadata.CollectionNamespaceMappingsDiffPair, srcManifestsDelta []*metadata.CollectionsManifest, latestTgtManifestId uint64)) *PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*metadata.CollectionNamespaceMappingsDiffPair), args[2].([]*metadata.CollectionsManifest), args[3].(uint64))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call) Return(_a0 error) *PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call) RunAndReturn(run func(string, *metadata.CollectionNamespaceMappingsDiffPair, []*metadata.CollectionsManifest, uint64) error) *PipelineOpSerializerIface_BackfillMappingStatusUpdate_Call {
	_c.Call.Return(run)
	return _c
}

// CleanBackfill provides a mock function with given fields: topic
func (_m *PipelineOpSerializerIface) CleanBackfill(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for CleanBackfill")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_CleanBackfill_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CleanBackfill'
type PipelineOpSerializerIface_CleanBackfill_Call struct {
	*mock.Call
}

// CleanBackfill is a helper method to define mock.On call
//   - topic string
func (_e *PipelineOpSerializerIface_Expecter) CleanBackfill(topic interface{}) *PipelineOpSerializerIface_CleanBackfill_Call {
	return &PipelineOpSerializerIface_CleanBackfill_Call{Call: _e.mock.On("CleanBackfill", topic)}
}

func (_c *PipelineOpSerializerIface_CleanBackfill_Call) Run(run func(topic string)) *PipelineOpSerializerIface_CleanBackfill_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_CleanBackfill_Call) Return(_a0 error) *PipelineOpSerializerIface_CleanBackfill_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_CleanBackfill_Call) RunAndReturn(run func(string) error) *PipelineOpSerializerIface_CleanBackfill_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: topic
func (_m *PipelineOpSerializerIface) Delete(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type PipelineOpSerializerIface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - topic string
func (_e *PipelineOpSerializerIface_Expecter) Delete(topic interface{}) *PipelineOpSerializerIface_Delete_Call {
	return &PipelineOpSerializerIface_Delete_Call{Call: _e.mock.On("Delete", topic)}
}

func (_c *PipelineOpSerializerIface_Delete_Call) Run(run func(topic string)) *PipelineOpSerializerIface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_Delete_Call) Return(_a0 error) *PipelineOpSerializerIface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_Delete_Call) RunAndReturn(run func(string) error) *PipelineOpSerializerIface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DismissEvent provides a mock function with given fields: topic, eventId
func (_m *PipelineOpSerializerIface) DismissEvent(topic string, eventId int) error {
	ret := _m.Called(topic, eventId)

	if len(ret) == 0 {
		panic("no return value specified for DismissEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int) error); ok {
		r0 = rf(topic, eventId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_DismissEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DismissEvent'
type PipelineOpSerializerIface_DismissEvent_Call struct {
	*mock.Call
}

// DismissEvent is a helper method to define mock.On call
//   - topic string
//   - eventId int
func (_e *PipelineOpSerializerIface_Expecter) DismissEvent(topic interface{}, eventId interface{}) *PipelineOpSerializerIface_DismissEvent_Call {
	return &PipelineOpSerializerIface_DismissEvent_Call{Call: _e.mock.On("DismissEvent", topic, eventId)}
}

func (_c *PipelineOpSerializerIface_DismissEvent_Call) Run(run func(topic string, eventId int)) *PipelineOpSerializerIface_DismissEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(int))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_DismissEvent_Call) Return(_a0 error) *PipelineOpSerializerIface_DismissEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_DismissEvent_Call) RunAndReturn(run func(string, int) error) *PipelineOpSerializerIface_DismissEvent_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrCreateReplicationStatus provides a mock function with given fields: topic, cur_err
func (_m *PipelineOpSerializerIface) GetOrCreateReplicationStatus(topic string, cur_err error) (*pipeline.ReplicationStatus, error) {
	ret := _m.Called(topic, cur_err)

	if len(ret) == 0 {
		panic("no return value specified for GetOrCreateReplicationStatus")
	}

	var r0 *pipeline.ReplicationStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(string, error) (*pipeline.ReplicationStatus, error)); ok {
		return rf(topic, cur_err)
	}
	if rf, ok := ret.Get(0).(func(string, error) *pipeline.ReplicationStatus); ok {
		r0 = rf(topic, cur_err)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pipeline.ReplicationStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(string, error) error); ok {
		r1 = rf(topic, cur_err)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrCreateReplicationStatus'
type PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call struct {
	*mock.Call
}

// GetOrCreateReplicationStatus is a helper method to define mock.On call
//   - topic string
//   - cur_err error
func (_e *PipelineOpSerializerIface_Expecter) GetOrCreateReplicationStatus(topic interface{}, cur_err interface{}) *PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call {
	return &PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call{Call: _e.mock.On("GetOrCreateReplicationStatus", topic, cur_err)}
}

func (_c *PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call) Run(run func(topic string, cur_err error)) *PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(error))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call) Return(_a0 *pipeline.ReplicationStatus, _a1 error) *PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call) RunAndReturn(run func(string, error) (*pipeline.ReplicationStatus, error)) *PipelineOpSerializerIface_GetOrCreateReplicationStatus_Call {
	_c.Call.Return(run)
	return _c
}

// Init provides a mock function with given fields: topic
func (_m *PipelineOpSerializerIface) Init(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for Init")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type PipelineOpSerializerIface_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - topic string
func (_e *PipelineOpSerializerIface_Expecter) Init(topic interface{}) *PipelineOpSerializerIface_Init_Call {
	return &PipelineOpSerializerIface_Init_Call{Call: _e.mock.On("Init", topic)}
}

func (_c *PipelineOpSerializerIface_Init_Call) Run(run func(topic string)) *PipelineOpSerializerIface_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_Init_Call) Return(_a0 error) *PipelineOpSerializerIface_Init_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_Init_Call) RunAndReturn(run func(string) error) *PipelineOpSerializerIface_Init_Call {
	_c.Call.Return(run)
	return _c
}

// Pause provides a mock function with given fields: topic
func (_m *PipelineOpSerializerIface) Pause(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for Pause")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_Pause_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Pause'
type PipelineOpSerializerIface_Pause_Call struct {
	*mock.Call
}

// Pause is a helper method to define mock.On call
//   - topic string
func (_e *PipelineOpSerializerIface_Expecter) Pause(topic interface{}) *PipelineOpSerializerIface_Pause_Call {
	return &PipelineOpSerializerIface_Pause_Call{Call: _e.mock.On("Pause", topic)}
}

func (_c *PipelineOpSerializerIface_Pause_Call) Run(run func(topic string)) *PipelineOpSerializerIface_Pause_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_Pause_Call) Return(_a0 error) *PipelineOpSerializerIface_Pause_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_Pause_Call) RunAndReturn(run func(string) error) *PipelineOpSerializerIface_Pause_Call {
	_c.Call.Return(run)
	return _c
}

// ReInit provides a mock function with given fields: topic
func (_m *PipelineOpSerializerIface) ReInit(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for ReInit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_ReInit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReInit'
type PipelineOpSerializerIface_ReInit_Call struct {
	*mock.Call
}

// ReInit is a helper method to define mock.On call
//   - topic string
func (_e *PipelineOpSerializerIface_Expecter) ReInit(topic interface{}) *PipelineOpSerializerIface_ReInit_Call {
	return &PipelineOpSerializerIface_ReInit_Call{Call: _e.mock.On("ReInit", topic)}
}

func (_c *PipelineOpSerializerIface_ReInit_Call) Run(run func(topic string)) *PipelineOpSerializerIface_ReInit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_ReInit_Call) Return(_a0 error) *PipelineOpSerializerIface_ReInit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_ReInit_Call) RunAndReturn(run func(string) error) *PipelineOpSerializerIface_ReInit_Call {
	_c.Call.Return(run)
	return _c
}

// StartBackfill provides a mock function with given fields: topic
func (_m *PipelineOpSerializerIface) StartBackfill(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for StartBackfill")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_StartBackfill_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartBackfill'
type PipelineOpSerializerIface_StartBackfill_Call struct {
	*mock.Call
}

// StartBackfill is a helper method to define mock.On call
//   - topic string
func (_e *PipelineOpSerializerIface_Expecter) StartBackfill(topic interface{}) *PipelineOpSerializerIface_StartBackfill_Call {
	return &PipelineOpSerializerIface_StartBackfill_Call{Call: _e.mock.On("StartBackfill", topic)}
}

func (_c *PipelineOpSerializerIface_StartBackfill_Call) Run(run func(topic string)) *PipelineOpSerializerIface_StartBackfill_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_StartBackfill_Call) Return(_a0 error) *PipelineOpSerializerIface_StartBackfill_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_StartBackfill_Call) RunAndReturn(run func(string) error) *PipelineOpSerializerIface_StartBackfill_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *PipelineOpSerializerIface) Stop() {
	_m.Called()
}

// PipelineOpSerializerIface_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type PipelineOpSerializerIface_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *PipelineOpSerializerIface_Expecter) Stop() *PipelineOpSerializerIface_Stop_Call {
	return &PipelineOpSerializerIface_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *PipelineOpSerializerIface_Stop_Call) Run(run func()) *PipelineOpSerializerIface_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PipelineOpSerializerIface_Stop_Call) Return() *PipelineOpSerializerIface_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *PipelineOpSerializerIface_Stop_Call) RunAndReturn(run func()) *PipelineOpSerializerIface_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// StopBackfill provides a mock function with given fields: topic
func (_m *PipelineOpSerializerIface) StopBackfill(topic string) error {
	ret := _m.Called(topic)

	if len(ret) == 0 {
		panic("no return value specified for StopBackfill")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(topic)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_StopBackfill_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopBackfill'
type PipelineOpSerializerIface_StopBackfill_Call struct {
	*mock.Call
}

// StopBackfill is a helper method to define mock.On call
//   - topic string
func (_e *PipelineOpSerializerIface_Expecter) StopBackfill(topic interface{}) *PipelineOpSerializerIface_StopBackfill_Call {
	return &PipelineOpSerializerIface_StopBackfill_Call{Call: _e.mock.On("StopBackfill", topic)}
}

func (_c *PipelineOpSerializerIface_StopBackfill_Call) Run(run func(topic string)) *PipelineOpSerializerIface_StopBackfill_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_StopBackfill_Call) Return(_a0 error) *PipelineOpSerializerIface_StopBackfill_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_StopBackfill_Call) RunAndReturn(run func(string) error) *PipelineOpSerializerIface_StopBackfill_Call {
	_c.Call.Return(run)
	return _c
}

// StopBackfillWithCb provides a mock function with given fields: pipelineName, cb, errCb, skipCkpt
func (_m *PipelineOpSerializerIface) StopBackfillWithCb(pipelineName string, cb base.StoppedPipelineCallback, errCb base.StoppedPipelineErrCallback, skipCkpt bool) error {
	ret := _m.Called(pipelineName, cb, errCb, skipCkpt)

	if len(ret) == 0 {
		panic("no return value specified for StopBackfillWithCb")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, base.StoppedPipelineCallback, base.StoppedPipelineErrCallback, bool) error); ok {
		r0 = rf(pipelineName, cb, errCb, skipCkpt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_StopBackfillWithCb_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopBackfillWithCb'
type PipelineOpSerializerIface_StopBackfillWithCb_Call struct {
	*mock.Call
}

// StopBackfillWithCb is a helper method to define mock.On call
//   - pipelineName string
//   - cb base.StoppedPipelineCallback
//   - errCb base.StoppedPipelineErrCallback
//   - skipCkpt bool
func (_e *PipelineOpSerializerIface_Expecter) StopBackfillWithCb(pipelineName interface{}, cb interface{}, errCb interface{}, skipCkpt interface{}) *PipelineOpSerializerIface_StopBackfillWithCb_Call {
	return &PipelineOpSerializerIface_StopBackfillWithCb_Call{Call: _e.mock.On("StopBackfillWithCb", pipelineName, cb, errCb, skipCkpt)}
}

func (_c *PipelineOpSerializerIface_StopBackfillWithCb_Call) Run(run func(pipelineName string, cb base.StoppedPipelineCallback, errCb base.StoppedPipelineErrCallback, skipCkpt bool)) *PipelineOpSerializerIface_StopBackfillWithCb_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(base.StoppedPipelineCallback), args[2].(base.StoppedPipelineErrCallback), args[3].(bool))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_StopBackfillWithCb_Call) Return(_a0 error) *PipelineOpSerializerIface_StopBackfillWithCb_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_StopBackfillWithCb_Call) RunAndReturn(run func(string, base.StoppedPipelineCallback, base.StoppedPipelineErrCallback, bool) error) *PipelineOpSerializerIface_StopBackfillWithCb_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: topic, err
func (_m *PipelineOpSerializerIface) Update(topic string, err error) error {
	ret := _m.Called(topic, err)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, error) error); ok {
		r0 = rf(topic, err)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type PipelineOpSerializerIface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - topic string
//   - err error
func (_e *PipelineOpSerializerIface_Expecter) Update(topic interface{}, err interface{}) *PipelineOpSerializerIface_Update_Call {
	return &PipelineOpSerializerIface_Update_Call{Call: _e.mock.On("Update", topic, err)}
}

func (_c *PipelineOpSerializerIface_Update_Call) Run(run func(topic string, err error)) *PipelineOpSerializerIface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(error))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_Update_Call) Return(_a0 error) *PipelineOpSerializerIface_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_Update_Call) RunAndReturn(run func(string, error) error) *PipelineOpSerializerIface_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateWithStoppedCb provides a mock function with given fields: topic, callback, errCb
func (_m *PipelineOpSerializerIface) UpdateWithStoppedCb(topic string, callback base.StoppedPipelineCallback, errCb base.StoppedPipelineErrCallback) error {
	ret := _m.Called(topic, callback, errCb)

	if len(ret) == 0 {
		panic("no return value specified for UpdateWithStoppedCb")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, base.StoppedPipelineCallback, base.StoppedPipelineErrCallback) error); ok {
		r0 = rf(topic, callback, errCb)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineOpSerializerIface_UpdateWithStoppedCb_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateWithStoppedCb'
type PipelineOpSerializerIface_UpdateWithStoppedCb_Call struct {
	*mock.Call
}

// UpdateWithStoppedCb is a helper method to define mock.On call
//   - topic string
//   - callback base.StoppedPipelineCallback
//   - errCb base.StoppedPipelineErrCallback
func (_e *PipelineOpSerializerIface_Expecter) UpdateWithStoppedCb(topic interface{}, callback interface{}, errCb interface{}) *PipelineOpSerializerIface_UpdateWithStoppedCb_Call {
	return &PipelineOpSerializerIface_UpdateWithStoppedCb_Call{Call: _e.mock.On("UpdateWithStoppedCb", topic, callback, errCb)}
}

func (_c *PipelineOpSerializerIface_UpdateWithStoppedCb_Call) Run(run func(topic string, callback base.StoppedPipelineCallback, errCb base.StoppedPipelineErrCallback)) *PipelineOpSerializerIface_UpdateWithStoppedCb_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(base.StoppedPipelineCallback), args[2].(base.StoppedPipelineErrCallback))
	})
	return _c
}

func (_c *PipelineOpSerializerIface_UpdateWithStoppedCb_Call) Return(_a0 error) *PipelineOpSerializerIface_UpdateWithStoppedCb_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineOpSerializerIface_UpdateWithStoppedCb_Call) RunAndReturn(run func(string, base.StoppedPipelineCallback, base.StoppedPipelineErrCallback) error) *PipelineOpSerializerIface_UpdateWithStoppedCb_Call {
	_c.Call.Return(run)
	return _c
}

// NewPipelineOpSerializerIface creates a new instance of PipelineOpSerializerIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPipelineOpSerializerIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *PipelineOpSerializerIface {
	mock := &PipelineOpSerializerIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
