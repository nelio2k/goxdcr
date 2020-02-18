package mocks

import (
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
