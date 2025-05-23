// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/v8/base"
	common "github.com/couchbase/goxdcr/v8/common"

	metadata "github.com/couchbase/goxdcr/v8/metadata"

	mock "github.com/stretchr/testify/mock"
)

// Pipeline is an autogenerated mock type for the Pipeline type
type Pipeline struct {
	mock.Mock
}

type Pipeline_Expecter struct {
	mock *mock.Mock
}

func (_m *Pipeline) EXPECT() *Pipeline_Expecter {
	return &Pipeline_Expecter{mock: &_m.Mock}
}

// FullTopic provides a mock function with no fields
func (_m *Pipeline) FullTopic() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FullTopic")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Pipeline_FullTopic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FullTopic'
type Pipeline_FullTopic_Call struct {
	*mock.Call
}

// FullTopic is a helper method to define mock.On call
func (_e *Pipeline_Expecter) FullTopic() *Pipeline_FullTopic_Call {
	return &Pipeline_FullTopic_Call{Call: _e.mock.On("FullTopic")}
}

func (_c *Pipeline_FullTopic_Call) Run(run func()) *Pipeline_FullTopic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_FullTopic_Call) Return(_a0 string) *Pipeline_FullTopic_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_FullTopic_Call) RunAndReturn(run func() string) *Pipeline_FullTopic_Call {
	_c.Call.Return(run)
	return _c
}

// GetAsyncListenerMap provides a mock function with no fields
func (_m *Pipeline) GetAsyncListenerMap() map[string]common.AsyncComponentEventListener {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAsyncListenerMap")
	}

	var r0 map[string]common.AsyncComponentEventListener
	if rf, ok := ret.Get(0).(func() map[string]common.AsyncComponentEventListener); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]common.AsyncComponentEventListener)
		}
	}

	return r0
}

// Pipeline_GetAsyncListenerMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAsyncListenerMap'
type Pipeline_GetAsyncListenerMap_Call struct {
	*mock.Call
}

// GetAsyncListenerMap is a helper method to define mock.On call
func (_e *Pipeline_Expecter) GetAsyncListenerMap() *Pipeline_GetAsyncListenerMap_Call {
	return &Pipeline_GetAsyncListenerMap_Call{Call: _e.mock.On("GetAsyncListenerMap")}
}

func (_c *Pipeline_GetAsyncListenerMap_Call) Run(run func()) *Pipeline_GetAsyncListenerMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_GetAsyncListenerMap_Call) Return(_a0 map[string]common.AsyncComponentEventListener) *Pipeline_GetAsyncListenerMap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_GetAsyncListenerMap_Call) RunAndReturn(run func() map[string]common.AsyncComponentEventListener) *Pipeline_GetAsyncListenerMap_Call {
	_c.Call.Return(run)
	return _c
}

// GetBrokenMapRO provides a mock function with no fields
func (_m *Pipeline) GetBrokenMapRO() (metadata.CollectionNamespaceMapping, func()) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetBrokenMapRO")
	}

	var r0 metadata.CollectionNamespaceMapping
	var r1 func()
	if rf, ok := ret.Get(0).(func() (metadata.CollectionNamespaceMapping, func())); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() metadata.CollectionNamespaceMapping); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metadata.CollectionNamespaceMapping)
		}
	}

	if rf, ok := ret.Get(1).(func() func()); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	return r0, r1
}

// Pipeline_GetBrokenMapRO_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBrokenMapRO'
type Pipeline_GetBrokenMapRO_Call struct {
	*mock.Call
}

// GetBrokenMapRO is a helper method to define mock.On call
func (_e *Pipeline_Expecter) GetBrokenMapRO() *Pipeline_GetBrokenMapRO_Call {
	return &Pipeline_GetBrokenMapRO_Call{Call: _e.mock.On("GetBrokenMapRO")}
}

func (_c *Pipeline_GetBrokenMapRO_Call) Run(run func()) *Pipeline_GetBrokenMapRO_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_GetBrokenMapRO_Call) Return(brokenMapReadOnly metadata.CollectionNamespaceMapping, callOnceDone func()) *Pipeline_GetBrokenMapRO_Call {
	_c.Call.Return(brokenMapReadOnly, callOnceDone)
	return _c
}

func (_c *Pipeline_GetBrokenMapRO_Call) RunAndReturn(run func() (metadata.CollectionNamespaceMapping, func())) *Pipeline_GetBrokenMapRO_Call {
	_c.Call.Return(run)
	return _c
}

// GetRebalanceProgress provides a mock function with no fields
func (_m *Pipeline) GetRebalanceProgress() (string, string) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRebalanceProgress")
	}

	var r0 string
	var r1 string
	if rf, ok := ret.Get(0).(func() (string, string)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() string); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// Pipeline_GetRebalanceProgress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRebalanceProgress'
type Pipeline_GetRebalanceProgress_Call struct {
	*mock.Call
}

// GetRebalanceProgress is a helper method to define mock.On call
func (_e *Pipeline_Expecter) GetRebalanceProgress() *Pipeline_GetRebalanceProgress_Call {
	return &Pipeline_GetRebalanceProgress_Call{Call: _e.mock.On("GetRebalanceProgress")}
}

func (_c *Pipeline_GetRebalanceProgress_Call) Run(run func()) *Pipeline_GetRebalanceProgress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_GetRebalanceProgress_Call) Return(_a0 string, _a1 string) *Pipeline_GetRebalanceProgress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Pipeline_GetRebalanceProgress_Call) RunAndReturn(run func() (string, string)) *Pipeline_GetRebalanceProgress_Call {
	_c.Call.Return(run)
	return _c
}

// InstanceId provides a mock function with no fields
func (_m *Pipeline) InstanceId() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for InstanceId")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Pipeline_InstanceId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InstanceId'
type Pipeline_InstanceId_Call struct {
	*mock.Call
}

// InstanceId is a helper method to define mock.On call
func (_e *Pipeline_Expecter) InstanceId() *Pipeline_InstanceId_Call {
	return &Pipeline_InstanceId_Call{Call: _e.mock.On("InstanceId")}
}

func (_c *Pipeline_InstanceId_Call) Run(run func()) *Pipeline_InstanceId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_InstanceId_Call) Return(_a0 string) *Pipeline_InstanceId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_InstanceId_Call) RunAndReturn(run func() string) *Pipeline_InstanceId_Call {
	_c.Call.Return(run)
	return _c
}

// ReportProgress provides a mock function with given fields: progress
func (_m *Pipeline) ReportProgress(progress string) {
	_m.Called(progress)
}

// Pipeline_ReportProgress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReportProgress'
type Pipeline_ReportProgress_Call struct {
	*mock.Call
}

// ReportProgress is a helper method to define mock.On call
//   - progress string
func (_e *Pipeline_Expecter) ReportProgress(progress interface{}) *Pipeline_ReportProgress_Call {
	return &Pipeline_ReportProgress_Call{Call: _e.mock.On("ReportProgress", progress)}
}

func (_c *Pipeline_ReportProgress_Call) Run(run func(progress string)) *Pipeline_ReportProgress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Pipeline_ReportProgress_Call) Return() *Pipeline_ReportProgress_Call {
	_c.Call.Return()
	return _c
}

func (_c *Pipeline_ReportProgress_Call) RunAndReturn(run func(string)) *Pipeline_ReportProgress_Call {
	_c.Run(run)
	return _c
}

// RuntimeContext provides a mock function with no fields
func (_m *Pipeline) RuntimeContext() common.PipelineRuntimeContext {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RuntimeContext")
	}

	var r0 common.PipelineRuntimeContext
	if rf, ok := ret.Get(0).(func() common.PipelineRuntimeContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.PipelineRuntimeContext)
		}
	}

	return r0
}

// Pipeline_RuntimeContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RuntimeContext'
type Pipeline_RuntimeContext_Call struct {
	*mock.Call
}

// RuntimeContext is a helper method to define mock.On call
func (_e *Pipeline_Expecter) RuntimeContext() *Pipeline_RuntimeContext_Call {
	return &Pipeline_RuntimeContext_Call{Call: _e.mock.On("RuntimeContext")}
}

func (_c *Pipeline_RuntimeContext_Call) Run(run func()) *Pipeline_RuntimeContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_RuntimeContext_Call) Return(_a0 common.PipelineRuntimeContext) *Pipeline_RuntimeContext_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_RuntimeContext_Call) RunAndReturn(run func() common.PipelineRuntimeContext) *Pipeline_RuntimeContext_Call {
	_c.Call.Return(run)
	return _c
}

// SetAsyncListenerMap provides a mock function with given fields: _a0
func (_m *Pipeline) SetAsyncListenerMap(_a0 map[string]common.AsyncComponentEventListener) {
	_m.Called(_a0)
}

// Pipeline_SetAsyncListenerMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetAsyncListenerMap'
type Pipeline_SetAsyncListenerMap_Call struct {
	*mock.Call
}

// SetAsyncListenerMap is a helper method to define mock.On call
//   - _a0 map[string]common.AsyncComponentEventListener
func (_e *Pipeline_Expecter) SetAsyncListenerMap(_a0 interface{}) *Pipeline_SetAsyncListenerMap_Call {
	return &Pipeline_SetAsyncListenerMap_Call{Call: _e.mock.On("SetAsyncListenerMap", _a0)}
}

func (_c *Pipeline_SetAsyncListenerMap_Call) Run(run func(_a0 map[string]common.AsyncComponentEventListener)) *Pipeline_SetAsyncListenerMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]common.AsyncComponentEventListener))
	})
	return _c
}

func (_c *Pipeline_SetAsyncListenerMap_Call) Return() *Pipeline_SetAsyncListenerMap_Call {
	_c.Call.Return()
	return _c
}

func (_c *Pipeline_SetAsyncListenerMap_Call) RunAndReturn(run func(map[string]common.AsyncComponentEventListener)) *Pipeline_SetAsyncListenerMap_Call {
	_c.Run(run)
	return _c
}

// SetBrokenMap provides a mock function with given fields: brokenMap
func (_m *Pipeline) SetBrokenMap(brokenMap metadata.CollectionNamespaceMapping) {
	_m.Called(brokenMap)
}

// Pipeline_SetBrokenMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetBrokenMap'
type Pipeline_SetBrokenMap_Call struct {
	*mock.Call
}

// SetBrokenMap is a helper method to define mock.On call
//   - brokenMap metadata.CollectionNamespaceMapping
func (_e *Pipeline_Expecter) SetBrokenMap(brokenMap interface{}) *Pipeline_SetBrokenMap_Call {
	return &Pipeline_SetBrokenMap_Call{Call: _e.mock.On("SetBrokenMap", brokenMap)}
}

func (_c *Pipeline_SetBrokenMap_Call) Run(run func(brokenMap metadata.CollectionNamespaceMapping)) *Pipeline_SetBrokenMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.CollectionNamespaceMapping))
	})
	return _c
}

func (_c *Pipeline_SetBrokenMap_Call) Return() *Pipeline_SetBrokenMap_Call {
	_c.Call.Return()
	return _c
}

func (_c *Pipeline_SetBrokenMap_Call) RunAndReturn(run func(metadata.CollectionNamespaceMapping)) *Pipeline_SetBrokenMap_Call {
	_c.Run(run)
	return _c
}

// SetProgressRecorder provides a mock function with given fields: recorder
func (_m *Pipeline) SetProgressRecorder(recorder common.PipelineProgressRecorder) {
	_m.Called(recorder)
}

// Pipeline_SetProgressRecorder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetProgressRecorder'
type Pipeline_SetProgressRecorder_Call struct {
	*mock.Call
}

// SetProgressRecorder is a helper method to define mock.On call
//   - recorder common.PipelineProgressRecorder
func (_e *Pipeline_Expecter) SetProgressRecorder(recorder interface{}) *Pipeline_SetProgressRecorder_Call {
	return &Pipeline_SetProgressRecorder_Call{Call: _e.mock.On("SetProgressRecorder", recorder)}
}

func (_c *Pipeline_SetProgressRecorder_Call) Run(run func(recorder common.PipelineProgressRecorder)) *Pipeline_SetProgressRecorder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.PipelineProgressRecorder))
	})
	return _c
}

func (_c *Pipeline_SetProgressRecorder_Call) Return() *Pipeline_SetProgressRecorder_Call {
	_c.Call.Return()
	return _c
}

func (_c *Pipeline_SetProgressRecorder_Call) RunAndReturn(run func(common.PipelineProgressRecorder)) *Pipeline_SetProgressRecorder_Call {
	_c.Run(run)
	return _c
}

// SetRuntimeContext provides a mock function with given fields: ctx
func (_m *Pipeline) SetRuntimeContext(ctx common.PipelineRuntimeContext) {
	_m.Called(ctx)
}

// Pipeline_SetRuntimeContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetRuntimeContext'
type Pipeline_SetRuntimeContext_Call struct {
	*mock.Call
}

// SetRuntimeContext is a helper method to define mock.On call
//   - ctx common.PipelineRuntimeContext
func (_e *Pipeline_Expecter) SetRuntimeContext(ctx interface{}) *Pipeline_SetRuntimeContext_Call {
	return &Pipeline_SetRuntimeContext_Call{Call: _e.mock.On("SetRuntimeContext", ctx)}
}

func (_c *Pipeline_SetRuntimeContext_Call) Run(run func(ctx common.PipelineRuntimeContext)) *Pipeline_SetRuntimeContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.PipelineRuntimeContext))
	})
	return _c
}

func (_c *Pipeline_SetRuntimeContext_Call) Return() *Pipeline_SetRuntimeContext_Call {
	_c.Call.Return()
	return _c
}

func (_c *Pipeline_SetRuntimeContext_Call) RunAndReturn(run func(common.PipelineRuntimeContext)) *Pipeline_SetRuntimeContext_Call {
	_c.Run(run)
	return _c
}

// SetState provides a mock function with given fields: state
func (_m *Pipeline) SetState(state common.PipelineState) error {
	ret := _m.Called(state)

	if len(ret) == 0 {
		panic("no return value specified for SetState")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(common.PipelineState) error); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Pipeline_SetState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetState'
type Pipeline_SetState_Call struct {
	*mock.Call
}

// SetState is a helper method to define mock.On call
//   - state common.PipelineState
func (_e *Pipeline_Expecter) SetState(state interface{}) *Pipeline_SetState_Call {
	return &Pipeline_SetState_Call{Call: _e.mock.On("SetState", state)}
}

func (_c *Pipeline_SetState_Call) Run(run func(state common.PipelineState)) *Pipeline_SetState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.PipelineState))
	})
	return _c
}

func (_c *Pipeline_SetState_Call) Return(_a0 error) *Pipeline_SetState_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_SetState_Call) RunAndReturn(run func(common.PipelineState) error) *Pipeline_SetState_Call {
	_c.Call.Return(run)
	return _c
}

// Settings provides a mock function with no fields
func (_m *Pipeline) Settings() metadata.ReplicationSettingsMap {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Settings")
	}

	var r0 metadata.ReplicationSettingsMap
	if rf, ok := ret.Get(0).(func() metadata.ReplicationSettingsMap); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metadata.ReplicationSettingsMap)
		}
	}

	return r0
}

// Pipeline_Settings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Settings'
type Pipeline_Settings_Call struct {
	*mock.Call
}

// Settings is a helper method to define mock.On call
func (_e *Pipeline_Expecter) Settings() *Pipeline_Settings_Call {
	return &Pipeline_Settings_Call{Call: _e.mock.On("Settings")}
}

func (_c *Pipeline_Settings_Call) Run(run func()) *Pipeline_Settings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_Settings_Call) Return(_a0 metadata.ReplicationSettingsMap) *Pipeline_Settings_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_Settings_Call) RunAndReturn(run func() metadata.ReplicationSettingsMap) *Pipeline_Settings_Call {
	_c.Call.Return(run)
	return _c
}

// Sources provides a mock function with no fields
func (_m *Pipeline) Sources() map[string]common.SourceNozzle {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Sources")
	}

	var r0 map[string]common.SourceNozzle
	if rf, ok := ret.Get(0).(func() map[string]common.SourceNozzle); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]common.SourceNozzle)
		}
	}

	return r0
}

// Pipeline_Sources_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Sources'
type Pipeline_Sources_Call struct {
	*mock.Call
}

// Sources is a helper method to define mock.On call
func (_e *Pipeline_Expecter) Sources() *Pipeline_Sources_Call {
	return &Pipeline_Sources_Call{Call: _e.mock.On("Sources")}
}

func (_c *Pipeline_Sources_Call) Run(run func()) *Pipeline_Sources_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_Sources_Call) Return(_a0 map[string]common.SourceNozzle) *Pipeline_Sources_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_Sources_Call) RunAndReturn(run func() map[string]common.SourceNozzle) *Pipeline_Sources_Call {
	_c.Call.Return(run)
	return _c
}

// Specification provides a mock function with no fields
func (_m *Pipeline) Specification() metadata.GenericSpecification {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Specification")
	}

	var r0 metadata.GenericSpecification
	if rf, ok := ret.Get(0).(func() metadata.GenericSpecification); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metadata.GenericSpecification)
		}
	}

	return r0
}

// Pipeline_Specification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Specification'
type Pipeline_Specification_Call struct {
	*mock.Call
}

// Specification is a helper method to define mock.On call
func (_e *Pipeline_Expecter) Specification() *Pipeline_Specification_Call {
	return &Pipeline_Specification_Call{Call: _e.mock.On("Specification")}
}

func (_c *Pipeline_Specification_Call) Run(run func()) *Pipeline_Specification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_Specification_Call) Return(_a0 metadata.GenericSpecification) *Pipeline_Specification_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_Specification_Call) RunAndReturn(run func() metadata.GenericSpecification) *Pipeline_Specification_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: settings
func (_m *Pipeline) Start(settings metadata.ReplicationSettingsMap) base.ErrorMap {
	ret := _m.Called(settings)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 base.ErrorMap
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) base.ErrorMap); ok {
		r0 = rf(settings)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(base.ErrorMap)
		}
	}

	return r0
}

// Pipeline_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type Pipeline_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - settings metadata.ReplicationSettingsMap
func (_e *Pipeline_Expecter) Start(settings interface{}) *Pipeline_Start_Call {
	return &Pipeline_Start_Call{Call: _e.mock.On("Start", settings)}
}

func (_c *Pipeline_Start_Call) Run(run func(settings metadata.ReplicationSettingsMap)) *Pipeline_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.ReplicationSettingsMap))
	})
	return _c
}

func (_c *Pipeline_Start_Call) Return(_a0 base.ErrorMap) *Pipeline_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_Start_Call) RunAndReturn(run func(metadata.ReplicationSettingsMap) base.ErrorMap) *Pipeline_Start_Call {
	_c.Call.Return(run)
	return _c
}

// State provides a mock function with no fields
func (_m *Pipeline) State() common.PipelineState {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for State")
	}

	var r0 common.PipelineState
	if rf, ok := ret.Get(0).(func() common.PipelineState); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.PipelineState)
	}

	return r0
}

// Pipeline_State_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'State'
type Pipeline_State_Call struct {
	*mock.Call
}

// State is a helper method to define mock.On call
func (_e *Pipeline_Expecter) State() *Pipeline_State_Call {
	return &Pipeline_State_Call{Call: _e.mock.On("State")}
}

func (_c *Pipeline_State_Call) Run(run func()) *Pipeline_State_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_State_Call) Return(_a0 common.PipelineState) *Pipeline_State_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_State_Call) RunAndReturn(run func() common.PipelineState) *Pipeline_State_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with no fields
func (_m *Pipeline) Stop() base.ErrorMap {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Stop")
	}

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

// Pipeline_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type Pipeline_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *Pipeline_Expecter) Stop() *Pipeline_Stop_Call {
	return &Pipeline_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *Pipeline_Stop_Call) Run(run func()) *Pipeline_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_Stop_Call) Return(_a0 base.ErrorMap) *Pipeline_Stop_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_Stop_Call) RunAndReturn(run func() base.ErrorMap) *Pipeline_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// Targets provides a mock function with no fields
func (_m *Pipeline) Targets() map[string]common.Nozzle {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Targets")
	}

	var r0 map[string]common.Nozzle
	if rf, ok := ret.Get(0).(func() map[string]common.Nozzle); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]common.Nozzle)
		}
	}

	return r0
}

// Pipeline_Targets_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Targets'
type Pipeline_Targets_Call struct {
	*mock.Call
}

// Targets is a helper method to define mock.On call
func (_e *Pipeline_Expecter) Targets() *Pipeline_Targets_Call {
	return &Pipeline_Targets_Call{Call: _e.mock.On("Targets")}
}

func (_c *Pipeline_Targets_Call) Run(run func()) *Pipeline_Targets_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_Targets_Call) Return(_a0 map[string]common.Nozzle) *Pipeline_Targets_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_Targets_Call) RunAndReturn(run func() map[string]common.Nozzle) *Pipeline_Targets_Call {
	_c.Call.Return(run)
	return _c
}

// Topic provides a mock function with no fields
func (_m *Pipeline) Topic() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Topic")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Pipeline_Topic_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Topic'
type Pipeline_Topic_Call struct {
	*mock.Call
}

// Topic is a helper method to define mock.On call
func (_e *Pipeline_Expecter) Topic() *Pipeline_Topic_Call {
	return &Pipeline_Topic_Call{Call: _e.mock.On("Topic")}
}

func (_c *Pipeline_Topic_Call) Run(run func()) *Pipeline_Topic_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_Topic_Call) Return(_a0 string) *Pipeline_Topic_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_Topic_Call) RunAndReturn(run func() string) *Pipeline_Topic_Call {
	_c.Call.Return(run)
	return _c
}

// Type provides a mock function with no fields
func (_m *Pipeline) Type() common.PipelineType {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Type")
	}

	var r0 common.PipelineType
	if rf, ok := ret.Get(0).(func() common.PipelineType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.PipelineType)
	}

	return r0
}

// Pipeline_Type_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Type'
type Pipeline_Type_Call struct {
	*mock.Call
}

// Type is a helper method to define mock.On call
func (_e *Pipeline_Expecter) Type() *Pipeline_Type_Call {
	return &Pipeline_Type_Call{Call: _e.mock.On("Type")}
}

func (_c *Pipeline_Type_Call) Run(run func()) *Pipeline_Type_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Pipeline_Type_Call) Return(_a0 common.PipelineType) *Pipeline_Type_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_Type_Call) RunAndReturn(run func() common.PipelineType) *Pipeline_Type_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSettings provides a mock function with given fields: settings
func (_m *Pipeline) UpdateSettings(settings metadata.ReplicationSettingsMap) error {
	ret := _m.Called(settings)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSettings")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) error); ok {
		r0 = rf(settings)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Pipeline_UpdateSettings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSettings'
type Pipeline_UpdateSettings_Call struct {
	*mock.Call
}

// UpdateSettings is a helper method to define mock.On call
//   - settings metadata.ReplicationSettingsMap
func (_e *Pipeline_Expecter) UpdateSettings(settings interface{}) *Pipeline_UpdateSettings_Call {
	return &Pipeline_UpdateSettings_Call{Call: _e.mock.On("UpdateSettings", settings)}
}

func (_c *Pipeline_UpdateSettings_Call) Run(run func(settings metadata.ReplicationSettingsMap)) *Pipeline_UpdateSettings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.ReplicationSettingsMap))
	})
	return _c
}

func (_c *Pipeline_UpdateSettings_Call) Return(_a0 error) *Pipeline_UpdateSettings_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Pipeline_UpdateSettings_Call) RunAndReturn(run func(metadata.ReplicationSettingsMap) error) *Pipeline_UpdateSettings_Call {
	_c.Call.Return(run)
	return _c
}

// NewPipeline creates a new instance of Pipeline. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPipeline(t interface {
	mock.TestingT
	Cleanup(func())
}) *Pipeline {
	mock := &Pipeline{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
