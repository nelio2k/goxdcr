// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/v8/common"
	metadata "github.com/couchbase/goxdcr/v8/metadata"

	mock "github.com/stretchr/testify/mock"
)

// PipelineSupervisorSvc is an autogenerated mock type for the PipelineSupervisorSvc type
type PipelineSupervisorSvc struct {
	mock.Mock
}

type PipelineSupervisorSvc_Expecter struct {
	mock *mock.Mock
}

func (_m *PipelineSupervisorSvc) EXPECT() *PipelineSupervisorSvc_Expecter {
	return &PipelineSupervisorSvc_Expecter{mock: &_m.Mock}
}

// Attach provides a mock function with given fields: pipeline
func (_m *PipelineSupervisorSvc) Attach(pipeline common.Pipeline) error {
	ret := _m.Called(pipeline)

	if len(ret) == 0 {
		panic("no return value specified for Attach")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(pipeline)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineSupervisorSvc_Attach_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Attach'
type PipelineSupervisorSvc_Attach_Call struct {
	*mock.Call
}

// Attach is a helper method to define mock.On call
//   - pipeline common.Pipeline
func (_e *PipelineSupervisorSvc_Expecter) Attach(pipeline interface{}) *PipelineSupervisorSvc_Attach_Call {
	return &PipelineSupervisorSvc_Attach_Call{Call: _e.mock.On("Attach", pipeline)}
}

func (_c *PipelineSupervisorSvc_Attach_Call) Run(run func(pipeline common.Pipeline)) *PipelineSupervisorSvc_Attach_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Pipeline))
	})
	return _c
}

func (_c *PipelineSupervisorSvc_Attach_Call) Return(_a0 error) *PipelineSupervisorSvc_Attach_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineSupervisorSvc_Attach_Call) RunAndReturn(run func(common.Pipeline) error) *PipelineSupervisorSvc_Attach_Call {
	_c.Call.Return(run)
	return _c
}

// Detach provides a mock function with given fields: pipeline
func (_m *PipelineSupervisorSvc) Detach(pipeline common.Pipeline) error {
	ret := _m.Called(pipeline)

	if len(ret) == 0 {
		panic("no return value specified for Detach")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(pipeline)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineSupervisorSvc_Detach_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Detach'
type PipelineSupervisorSvc_Detach_Call struct {
	*mock.Call
}

// Detach is a helper method to define mock.On call
//   - pipeline common.Pipeline
func (_e *PipelineSupervisorSvc_Expecter) Detach(pipeline interface{}) *PipelineSupervisorSvc_Detach_Call {
	return &PipelineSupervisorSvc_Detach_Call{Call: _e.mock.On("Detach", pipeline)}
}

func (_c *PipelineSupervisorSvc_Detach_Call) Run(run func(pipeline common.Pipeline)) *PipelineSupervisorSvc_Detach_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Pipeline))
	})
	return _c
}

func (_c *PipelineSupervisorSvc_Detach_Call) Return(_a0 error) *PipelineSupervisorSvc_Detach_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineSupervisorSvc_Detach_Call) RunAndReturn(run func(common.Pipeline) error) *PipelineSupervisorSvc_Detach_Call {
	_c.Call.Return(run)
	return _c
}

// IsSharable provides a mock function with given fields:
func (_m *PipelineSupervisorSvc) IsSharable() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsSharable")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// PipelineSupervisorSvc_IsSharable_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsSharable'
type PipelineSupervisorSvc_IsSharable_Call struct {
	*mock.Call
}

// IsSharable is a helper method to define mock.On call
func (_e *PipelineSupervisorSvc_Expecter) IsSharable() *PipelineSupervisorSvc_IsSharable_Call {
	return &PipelineSupervisorSvc_IsSharable_Call{Call: _e.mock.On("IsSharable")}
}

func (_c *PipelineSupervisorSvc_IsSharable_Call) Run(run func()) *PipelineSupervisorSvc_IsSharable_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PipelineSupervisorSvc_IsSharable_Call) Return(_a0 bool) *PipelineSupervisorSvc_IsSharable_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineSupervisorSvc_IsSharable_Call) RunAndReturn(run func() bool) *PipelineSupervisorSvc_IsSharable_Call {
	_c.Call.Return(run)
	return _c
}

// OnEvent provides a mock function with given fields: event
func (_m *PipelineSupervisorSvc) OnEvent(event *common.Event) {
	_m.Called(event)
}

// PipelineSupervisorSvc_OnEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnEvent'
type PipelineSupervisorSvc_OnEvent_Call struct {
	*mock.Call
}

// OnEvent is a helper method to define mock.On call
//   - event *common.Event
func (_e *PipelineSupervisorSvc_Expecter) OnEvent(event interface{}) *PipelineSupervisorSvc_OnEvent_Call {
	return &PipelineSupervisorSvc_OnEvent_Call{Call: _e.mock.On("OnEvent", event)}
}

func (_c *PipelineSupervisorSvc_OnEvent_Call) Run(run func(event *common.Event)) *PipelineSupervisorSvc_OnEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*common.Event))
	})
	return _c
}

func (_c *PipelineSupervisorSvc_OnEvent_Call) Return() *PipelineSupervisorSvc_OnEvent_Call {
	_c.Call.Return()
	return _c
}

func (_c *PipelineSupervisorSvc_OnEvent_Call) RunAndReturn(run func(*common.Event)) *PipelineSupervisorSvc_OnEvent_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: _a0
func (_m *PipelineSupervisorSvc) Start(_a0 metadata.ReplicationSettingsMap) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineSupervisorSvc_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type PipelineSupervisorSvc_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - _a0 metadata.ReplicationSettingsMap
func (_e *PipelineSupervisorSvc_Expecter) Start(_a0 interface{}) *PipelineSupervisorSvc_Start_Call {
	return &PipelineSupervisorSvc_Start_Call{Call: _e.mock.On("Start", _a0)}
}

func (_c *PipelineSupervisorSvc_Start_Call) Run(run func(_a0 metadata.ReplicationSettingsMap)) *PipelineSupervisorSvc_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.ReplicationSettingsMap))
	})
	return _c
}

func (_c *PipelineSupervisorSvc_Start_Call) Return(_a0 error) *PipelineSupervisorSvc_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineSupervisorSvc_Start_Call) RunAndReturn(run func(metadata.ReplicationSettingsMap) error) *PipelineSupervisorSvc_Start_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *PipelineSupervisorSvc) Stop() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Stop")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PipelineSupervisorSvc_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type PipelineSupervisorSvc_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *PipelineSupervisorSvc_Expecter) Stop() *PipelineSupervisorSvc_Stop_Call {
	return &PipelineSupervisorSvc_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *PipelineSupervisorSvc_Stop_Call) Run(run func()) *PipelineSupervisorSvc_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *PipelineSupervisorSvc_Stop_Call) Return(_a0 error) *PipelineSupervisorSvc_Stop_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineSupervisorSvc_Stop_Call) RunAndReturn(run func() error) *PipelineSupervisorSvc_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSettings provides a mock function with given fields: settings
func (_m *PipelineSupervisorSvc) UpdateSettings(settings metadata.ReplicationSettingsMap) error {
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

// PipelineSupervisorSvc_UpdateSettings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSettings'
type PipelineSupervisorSvc_UpdateSettings_Call struct {
	*mock.Call
}

// UpdateSettings is a helper method to define mock.On call
//   - settings metadata.ReplicationSettingsMap
func (_e *PipelineSupervisorSvc_Expecter) UpdateSettings(settings interface{}) *PipelineSupervisorSvc_UpdateSettings_Call {
	return &PipelineSupervisorSvc_UpdateSettings_Call{Call: _e.mock.On("UpdateSettings", settings)}
}

func (_c *PipelineSupervisorSvc_UpdateSettings_Call) Run(run func(settings metadata.ReplicationSettingsMap)) *PipelineSupervisorSvc_UpdateSettings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.ReplicationSettingsMap))
	})
	return _c
}

func (_c *PipelineSupervisorSvc_UpdateSettings_Call) Return(_a0 error) *PipelineSupervisorSvc_UpdateSettings_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PipelineSupervisorSvc_UpdateSettings_Call) RunAndReturn(run func(metadata.ReplicationSettingsMap) error) *PipelineSupervisorSvc_UpdateSettings_Call {
	_c.Call.Return(run)
	return _c
}

// NewPipelineSupervisorSvc creates a new instance of PipelineSupervisorSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPipelineSupervisorSvc(t interface {
	mock.TestingT
	Cleanup(func())
}) *PipelineSupervisorSvc {
	mock := &PipelineSupervisorSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
