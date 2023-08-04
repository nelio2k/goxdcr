// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	mock "github.com/stretchr/testify/mock"
)

// AsyncComponentEventHandler is an autogenerated mock type for the AsyncComponentEventHandler type
type AsyncComponentEventHandler struct {
	mock.Mock
}

type AsyncComponentEventHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *AsyncComponentEventHandler) EXPECT() *AsyncComponentEventHandler_Expecter {
	return &AsyncComponentEventHandler_Expecter{mock: &_m.Mock}
}

// Id provides a mock function with given fields:
func (_m *AsyncComponentEventHandler) Id() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AsyncComponentEventHandler_Id_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Id'
type AsyncComponentEventHandler_Id_Call struct {
	*mock.Call
}

// Id is a helper method to define mock.On call
func (_e *AsyncComponentEventHandler_Expecter) Id() *AsyncComponentEventHandler_Id_Call {
	return &AsyncComponentEventHandler_Id_Call{Call: _e.mock.On("Id")}
}

func (_c *AsyncComponentEventHandler_Id_Call) Run(run func()) *AsyncComponentEventHandler_Id_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AsyncComponentEventHandler_Id_Call) Return(_a0 string) *AsyncComponentEventHandler_Id_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AsyncComponentEventHandler_Id_Call) RunAndReturn(run func() string) *AsyncComponentEventHandler_Id_Call {
	_c.Call.Return(run)
	return _c
}

// ProcessEvent provides a mock function with given fields: event
func (_m *AsyncComponentEventHandler) ProcessEvent(event *common.Event) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(*common.Event) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AsyncComponentEventHandler_ProcessEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProcessEvent'
type AsyncComponentEventHandler_ProcessEvent_Call struct {
	*mock.Call
}

// ProcessEvent is a helper method to define mock.On call
//   - event *common.Event
func (_e *AsyncComponentEventHandler_Expecter) ProcessEvent(event interface{}) *AsyncComponentEventHandler_ProcessEvent_Call {
	return &AsyncComponentEventHandler_ProcessEvent_Call{Call: _e.mock.On("ProcessEvent", event)}
}

func (_c *AsyncComponentEventHandler_ProcessEvent_Call) Run(run func(event *common.Event)) *AsyncComponentEventHandler_ProcessEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*common.Event))
	})
	return _c
}

func (_c *AsyncComponentEventHandler_ProcessEvent_Call) Return(_a0 error) *AsyncComponentEventHandler_ProcessEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AsyncComponentEventHandler_ProcessEvent_Call) RunAndReturn(run func(*common.Event) error) *AsyncComponentEventHandler_ProcessEvent_Call {
	_c.Call.Return(run)
	return _c
}

// NewAsyncComponentEventHandler creates a new instance of AsyncComponentEventHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAsyncComponentEventHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *AsyncComponentEventHandler {
	mock := &AsyncComponentEventHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}