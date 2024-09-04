// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/v8/base"
	common "github.com/couchbase/goxdcr/v8/common"

	mock "github.com/stretchr/testify/mock"
)

// PipelineFactory is an autogenerated mock type for the PipelineFactory type
type PipelineFactory struct {
	mock.Mock
}

type PipelineFactory_Expecter struct {
	mock *mock.Mock
}

func (_m *PipelineFactory) EXPECT() *PipelineFactory_Expecter {
	return &PipelineFactory_Expecter{mock: &_m.Mock}
}

// NewPipeline provides a mock function with given fields: topic, progressRecorder
func (_m *PipelineFactory) NewPipeline(topic string, progressRecorder common.PipelineProgressRecorder) (common.Pipeline, error) {
	ret := _m.Called(topic, progressRecorder)

	if len(ret) == 0 {
		panic("no return value specified for NewPipeline")
	}

	var r0 common.Pipeline
	var r1 error
	if rf, ok := ret.Get(0).(func(string, common.PipelineProgressRecorder) (common.Pipeline, error)); ok {
		return rf(topic, progressRecorder)
	}
	if rf, ok := ret.Get(0).(func(string, common.PipelineProgressRecorder) common.Pipeline); ok {
		r0 = rf(topic, progressRecorder)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Pipeline)
		}
	}

	if rf, ok := ret.Get(1).(func(string, common.PipelineProgressRecorder) error); ok {
		r1 = rf(topic, progressRecorder)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PipelineFactory_NewPipeline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewPipeline'
type PipelineFactory_NewPipeline_Call struct {
	*mock.Call
}

// NewPipeline is a helper method to define mock.On call
//   - topic string
//   - progressRecorder common.PipelineProgressRecorder
func (_e *PipelineFactory_Expecter) NewPipeline(topic interface{}, progressRecorder interface{}) *PipelineFactory_NewPipeline_Call {
	return &PipelineFactory_NewPipeline_Call{Call: _e.mock.On("NewPipeline", topic, progressRecorder)}
}

func (_c *PipelineFactory_NewPipeline_Call) Run(run func(topic string, progressRecorder common.PipelineProgressRecorder)) *PipelineFactory_NewPipeline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(common.PipelineProgressRecorder))
	})
	return _c
}

func (_c *PipelineFactory_NewPipeline_Call) Return(_a0 common.Pipeline, _a1 error) *PipelineFactory_NewPipeline_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PipelineFactory_NewPipeline_Call) RunAndReturn(run func(string, common.PipelineProgressRecorder) (common.Pipeline, error)) *PipelineFactory_NewPipeline_Call {
	_c.Call.Return(run)
	return _c
}

// NewSecondaryPipeline provides a mock function with given fields: topic, primaryPipeline, progress_recorder, pipelineType
func (_m *PipelineFactory) NewSecondaryPipeline(topic string, primaryPipeline common.Pipeline, progress_recorder common.PipelineProgressRecorder, pipelineType common.PipelineType) (common.Pipeline, error) {
	ret := _m.Called(topic, primaryPipeline, progress_recorder, pipelineType)

	if len(ret) == 0 {
		panic("no return value specified for NewSecondaryPipeline")
	}

	var r0 common.Pipeline
	var r1 error
	if rf, ok := ret.Get(0).(func(string, common.Pipeline, common.PipelineProgressRecorder, common.PipelineType) (common.Pipeline, error)); ok {
		return rf(topic, primaryPipeline, progress_recorder, pipelineType)
	}
	if rf, ok := ret.Get(0).(func(string, common.Pipeline, common.PipelineProgressRecorder, common.PipelineType) common.Pipeline); ok {
		r0 = rf(topic, primaryPipeline, progress_recorder, pipelineType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Pipeline)
		}
	}

	if rf, ok := ret.Get(1).(func(string, common.Pipeline, common.PipelineProgressRecorder, common.PipelineType) error); ok {
		r1 = rf(topic, primaryPipeline, progress_recorder, pipelineType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PipelineFactory_NewSecondaryPipeline_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewSecondaryPipeline'
type PipelineFactory_NewSecondaryPipeline_Call struct {
	*mock.Call
}

// NewSecondaryPipeline is a helper method to define mock.On call
//   - topic string
//   - primaryPipeline common.Pipeline
//   - progress_recorder common.PipelineProgressRecorder
//   - pipelineType common.PipelineType
func (_e *PipelineFactory_Expecter) NewSecondaryPipeline(topic interface{}, primaryPipeline interface{}, progress_recorder interface{}, pipelineType interface{}) *PipelineFactory_NewSecondaryPipeline_Call {
	return &PipelineFactory_NewSecondaryPipeline_Call{Call: _e.mock.On("NewSecondaryPipeline", topic, primaryPipeline, progress_recorder, pipelineType)}
}

func (_c *PipelineFactory_NewSecondaryPipeline_Call) Run(run func(topic string, primaryPipeline common.Pipeline, progress_recorder common.PipelineProgressRecorder, pipelineType common.PipelineType)) *PipelineFactory_NewSecondaryPipeline_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(common.Pipeline), args[2].(common.PipelineProgressRecorder), args[3].(common.PipelineType))
	})
	return _c
}

func (_c *PipelineFactory_NewSecondaryPipeline_Call) Return(_a0 common.Pipeline, _a1 error) *PipelineFactory_NewSecondaryPipeline_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PipelineFactory_NewSecondaryPipeline_Call) RunAndReturn(run func(string, common.Pipeline, common.PipelineProgressRecorder, common.PipelineType) (common.Pipeline, error)) *PipelineFactory_NewSecondaryPipeline_Call {
	_c.Call.Return(run)
	return _c
}

// SetPipelineStopCallback provides a mock function with given fields: cbType
func (_m *PipelineFactory) SetPipelineStopCallback(cbType base.PipelineMgrStopCbType) {
	_m.Called(cbType)
}

// PipelineFactory_SetPipelineStopCallback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetPipelineStopCallback'
type PipelineFactory_SetPipelineStopCallback_Call struct {
	*mock.Call
}

// SetPipelineStopCallback is a helper method to define mock.On call
//   - cbType base.PipelineMgrStopCbType
func (_e *PipelineFactory_Expecter) SetPipelineStopCallback(cbType interface{}) *PipelineFactory_SetPipelineStopCallback_Call {
	return &PipelineFactory_SetPipelineStopCallback_Call{Call: _e.mock.On("SetPipelineStopCallback", cbType)}
}

func (_c *PipelineFactory_SetPipelineStopCallback_Call) Run(run func(cbType base.PipelineMgrStopCbType)) *PipelineFactory_SetPipelineStopCallback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(base.PipelineMgrStopCbType))
	})
	return _c
}

func (_c *PipelineFactory_SetPipelineStopCallback_Call) Return() *PipelineFactory_SetPipelineStopCallback_Call {
	_c.Call.Return()
	return _c
}

func (_c *PipelineFactory_SetPipelineStopCallback_Call) RunAndReturn(run func(base.PipelineMgrStopCbType)) *PipelineFactory_SetPipelineStopCallback_Call {
	_c.Call.Return(run)
	return _c
}

// NewPipelineFactory creates a new instance of PipelineFactory. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPipelineFactory(t interface {
	mock.TestingT
	Cleanup(func())
}) *PipelineFactory {
	mock := &PipelineFactory{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
