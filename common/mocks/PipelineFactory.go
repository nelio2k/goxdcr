// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	common "github.com/couchbase/goxdcr/common"

	mock "github.com/stretchr/testify/mock"
)

// PipelineFactory is an autogenerated mock type for the PipelineFactory type
type PipelineFactory struct {
	mock.Mock
}

// NewPipeline provides a mock function with given fields: topic, progressRecorder
func (_m *PipelineFactory) NewPipeline(topic string, progressRecorder common.PipelineProgressRecorder) (common.Pipeline, error) {
	ret := _m.Called(topic, progressRecorder)

	var r0 common.Pipeline
	if rf, ok := ret.Get(0).(func(string, common.PipelineProgressRecorder) common.Pipeline); ok {
		r0 = rf(topic, progressRecorder)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Pipeline)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, common.PipelineProgressRecorder) error); ok {
		r1 = rf(topic, progressRecorder)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSecondaryPipeline provides a mock function with given fields: topic, primaryPipeline, progress_recorder, pipelineType
func (_m *PipelineFactory) NewSecondaryPipeline(topic string, primaryPipeline common.Pipeline, progress_recorder common.PipelineProgressRecorder, pipelineType common.PipelineType) (common.Pipeline, error) {
	ret := _m.Called(topic, primaryPipeline, progress_recorder, pipelineType)

	var r0 common.Pipeline
	if rf, ok := ret.Get(0).(func(string, common.Pipeline, common.PipelineProgressRecorder, common.PipelineType) common.Pipeline); ok {
		r0 = rf(topic, primaryPipeline, progress_recorder, pipelineType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Pipeline)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, common.Pipeline, common.PipelineProgressRecorder, common.PipelineType) error); ok {
		r1 = rf(topic, primaryPipeline, progress_recorder, pipelineType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPipelineStopCallback provides a mock function with given fields: cbType
func (_m *PipelineFactory) SetPipelineStopCallback(cbType base.PipelineMgrStopCbType) {
	_m.Called(cbType)
}
