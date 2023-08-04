// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	common "github.com/couchbase/goxdcr/common"

	log "github.com/couchbase/goxdcr/log"

	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"

	parts "github.com/couchbase/goxdcr/parts"
)

// DcpNozzleIface is an autogenerated mock type for the DcpNozzleIface type
type DcpNozzleIface struct {
	mock.Mock
}

type DcpNozzleIface_Expecter struct {
	mock *mock.Mock
}

func (_m *DcpNozzleIface) EXPECT() *DcpNozzleIface_Expecter {
	return &DcpNozzleIface_Expecter{mock: &_m.Mock}
}

// AsyncComponentEventListeners provides a mock function with given fields:
func (_m *DcpNozzleIface) AsyncComponentEventListeners() map[string]common.AsyncComponentEventListener {
	ret := _m.Called()

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

// DcpNozzleIface_AsyncComponentEventListeners_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AsyncComponentEventListeners'
type DcpNozzleIface_AsyncComponentEventListeners_Call struct {
	*mock.Call
}

// AsyncComponentEventListeners is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) AsyncComponentEventListeners() *DcpNozzleIface_AsyncComponentEventListeners_Call {
	return &DcpNozzleIface_AsyncComponentEventListeners_Call{Call: _e.mock.On("AsyncComponentEventListeners")}
}

func (_c *DcpNozzleIface_AsyncComponentEventListeners_Call) Run(run func()) *DcpNozzleIface_AsyncComponentEventListeners_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_AsyncComponentEventListeners_Call) Return(_a0 map[string]common.AsyncComponentEventListener) *DcpNozzleIface_AsyncComponentEventListeners_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_AsyncComponentEventListeners_Call) RunAndReturn(run func() map[string]common.AsyncComponentEventListener) *DcpNozzleIface_AsyncComponentEventListeners_Call {
	_c.Call.Return(run)
	return _c
}

// CheckStuckness provides a mock function with given fields: dcp_stats
func (_m *DcpNozzleIface) CheckStuckness(dcp_stats base.DcpStatsMapType) error {
	ret := _m.Called(dcp_stats)

	var r0 error
	if rf, ok := ret.Get(0).(func(base.DcpStatsMapType) error); ok {
		r0 = rf(dcp_stats)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_CheckStuckness_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckStuckness'
type DcpNozzleIface_CheckStuckness_Call struct {
	*mock.Call
}

// CheckStuckness is a helper method to define mock.On call
//   - dcp_stats base.DcpStatsMapType
func (_e *DcpNozzleIface_Expecter) CheckStuckness(dcp_stats interface{}) *DcpNozzleIface_CheckStuckness_Call {
	return &DcpNozzleIface_CheckStuckness_Call{Call: _e.mock.On("CheckStuckness", dcp_stats)}
}

func (_c *DcpNozzleIface_CheckStuckness_Call) Run(run func(dcp_stats base.DcpStatsMapType)) *DcpNozzleIface_CheckStuckness_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(base.DcpStatsMapType))
	})
	return _c
}

func (_c *DcpNozzleIface_CheckStuckness_Call) Return(_a0 error) *DcpNozzleIface_CheckStuckness_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_CheckStuckness_Call) RunAndReturn(run func(base.DcpStatsMapType) error) *DcpNozzleIface_CheckStuckness_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *DcpNozzleIface) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type DcpNozzleIface_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) Close() *DcpNozzleIface_Close_Call {
	return &DcpNozzleIface_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *DcpNozzleIface_Close_Call) Run(run func()) *DcpNozzleIface_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_Close_Call) Return(_a0 error) *DcpNozzleIface_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_Close_Call) RunAndReturn(run func() error) *DcpNozzleIface_Close_Call {
	_c.Call.Return(run)
	return _c
}

// CollectionEnabled provides a mock function with given fields:
func (_m *DcpNozzleIface) CollectionEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// DcpNozzleIface_CollectionEnabled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CollectionEnabled'
type DcpNozzleIface_CollectionEnabled_Call struct {
	*mock.Call
}

// CollectionEnabled is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) CollectionEnabled() *DcpNozzleIface_CollectionEnabled_Call {
	return &DcpNozzleIface_CollectionEnabled_Call{Call: _e.mock.On("CollectionEnabled")}
}

func (_c *DcpNozzleIface_CollectionEnabled_Call) Run(run func()) *DcpNozzleIface_CollectionEnabled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_CollectionEnabled_Call) Return(_a0 bool) *DcpNozzleIface_CollectionEnabled_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_CollectionEnabled_Call) RunAndReturn(run func() bool) *DcpNozzleIface_CollectionEnabled_Call {
	_c.Call.Return(run)
	return _c
}

// Connector provides a mock function with given fields:
func (_m *DcpNozzleIface) Connector() common.Connector {
	ret := _m.Called()

	var r0 common.Connector
	if rf, ok := ret.Get(0).(func() common.Connector); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Connector)
		}
	}

	return r0
}

// DcpNozzleIface_Connector_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connector'
type DcpNozzleIface_Connector_Call struct {
	*mock.Call
}

// Connector is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) Connector() *DcpNozzleIface_Connector_Call {
	return &DcpNozzleIface_Connector_Call{Call: _e.mock.On("Connector")}
}

func (_c *DcpNozzleIface_Connector_Call) Run(run func()) *DcpNozzleIface_Connector_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_Connector_Call) Return(_a0 common.Connector) *DcpNozzleIface_Connector_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_Connector_Call) RunAndReturn(run func() common.Connector) *DcpNozzleIface_Connector_Call {
	_c.Call.Return(run)
	return _c
}

// GetStreamState provides a mock function with given fields: vbno
func (_m *DcpNozzleIface) GetStreamState(vbno uint16) (parts.DcpStreamState, error) {
	ret := _m.Called(vbno)

	var r0 parts.DcpStreamState
	var r1 error
	if rf, ok := ret.Get(0).(func(uint16) (parts.DcpStreamState, error)); ok {
		return rf(vbno)
	}
	if rf, ok := ret.Get(0).(func(uint16) parts.DcpStreamState); ok {
		r0 = rf(vbno)
	} else {
		r0 = ret.Get(0).(parts.DcpStreamState)
	}

	if rf, ok := ret.Get(1).(func(uint16) error); ok {
		r1 = rf(vbno)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DcpNozzleIface_GetStreamState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetStreamState'
type DcpNozzleIface_GetStreamState_Call struct {
	*mock.Call
}

// GetStreamState is a helper method to define mock.On call
//   - vbno uint16
func (_e *DcpNozzleIface_Expecter) GetStreamState(vbno interface{}) *DcpNozzleIface_GetStreamState_Call {
	return &DcpNozzleIface_GetStreamState_Call{Call: _e.mock.On("GetStreamState", vbno)}
}

func (_c *DcpNozzleIface_GetStreamState_Call) Run(run func(vbno uint16)) *DcpNozzleIface_GetStreamState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(uint16))
	})
	return _c
}

func (_c *DcpNozzleIface_GetStreamState_Call) Return(_a0 parts.DcpStreamState, _a1 error) *DcpNozzleIface_GetStreamState_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DcpNozzleIface_GetStreamState_Call) RunAndReturn(run func(uint16) (parts.DcpStreamState, error)) *DcpNozzleIface_GetStreamState_Call {
	_c.Call.Return(run)
	return _c
}

// Id provides a mock function with given fields:
func (_m *DcpNozzleIface) Id() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// DcpNozzleIface_Id_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Id'
type DcpNozzleIface_Id_Call struct {
	*mock.Call
}

// Id is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) Id() *DcpNozzleIface_Id_Call {
	return &DcpNozzleIface_Id_Call{Call: _e.mock.On("Id")}
}

func (_c *DcpNozzleIface_Id_Call) Run(run func()) *DcpNozzleIface_Id_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_Id_Call) Return(_a0 string) *DcpNozzleIface_Id_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_Id_Call) RunAndReturn(run func() string) *DcpNozzleIface_Id_Call {
	_c.Call.Return(run)
	return _c
}

// IsOpen provides a mock function with given fields:
func (_m *DcpNozzleIface) IsOpen() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// DcpNozzleIface_IsOpen_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsOpen'
type DcpNozzleIface_IsOpen_Call struct {
	*mock.Call
}

// IsOpen is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) IsOpen() *DcpNozzleIface_IsOpen_Call {
	return &DcpNozzleIface_IsOpen_Call{Call: _e.mock.On("IsOpen")}
}

func (_c *DcpNozzleIface_IsOpen_Call) Run(run func()) *DcpNozzleIface_IsOpen_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_IsOpen_Call) Return(_a0 bool) *DcpNozzleIface_IsOpen_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_IsOpen_Call) RunAndReturn(run func() bool) *DcpNozzleIface_IsOpen_Call {
	_c.Call.Return(run)
	return _c
}

// Logger provides a mock function with given fields:
func (_m *DcpNozzleIface) Logger() *log.CommonLogger {
	ret := _m.Called()

	var r0 *log.CommonLogger
	if rf, ok := ret.Get(0).(func() *log.CommonLogger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*log.CommonLogger)
		}
	}

	return r0
}

// DcpNozzleIface_Logger_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Logger'
type DcpNozzleIface_Logger_Call struct {
	*mock.Call
}

// Logger is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) Logger() *DcpNozzleIface_Logger_Call {
	return &DcpNozzleIface_Logger_Call{Call: _e.mock.On("Logger")}
}

func (_c *DcpNozzleIface_Logger_Call) Run(run func()) *DcpNozzleIface_Logger_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_Logger_Call) Return(_a0 *log.CommonLogger) *DcpNozzleIface_Logger_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_Logger_Call) RunAndReturn(run func() *log.CommonLogger) *DcpNozzleIface_Logger_Call {
	_c.Call.Return(run)
	return _c
}

// Open provides a mock function with given fields:
func (_m *DcpNozzleIface) Open() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_Open_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Open'
type DcpNozzleIface_Open_Call struct {
	*mock.Call
}

// Open is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) Open() *DcpNozzleIface_Open_Call {
	return &DcpNozzleIface_Open_Call{Call: _e.mock.On("Open")}
}

func (_c *DcpNozzleIface_Open_Call) Run(run func()) *DcpNozzleIface_Open_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_Open_Call) Return(_a0 error) *DcpNozzleIface_Open_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_Open_Call) RunAndReturn(run func() error) *DcpNozzleIface_Open_Call {
	_c.Call.Return(run)
	return _c
}

// PrintStatusSummary provides a mock function with given fields:
func (_m *DcpNozzleIface) PrintStatusSummary() {
	_m.Called()
}

// DcpNozzleIface_PrintStatusSummary_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PrintStatusSummary'
type DcpNozzleIface_PrintStatusSummary_Call struct {
	*mock.Call
}

// PrintStatusSummary is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) PrintStatusSummary() *DcpNozzleIface_PrintStatusSummary_Call {
	return &DcpNozzleIface_PrintStatusSummary_Call{Call: _e.mock.On("PrintStatusSummary")}
}

func (_c *DcpNozzleIface_PrintStatusSummary_Call) Run(run func()) *DcpNozzleIface_PrintStatusSummary_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_PrintStatusSummary_Call) Return() *DcpNozzleIface_PrintStatusSummary_Call {
	_c.Call.Return()
	return _c
}

func (_c *DcpNozzleIface_PrintStatusSummary_Call) RunAndReturn(run func()) *DcpNozzleIface_PrintStatusSummary_Call {
	_c.Call.Return(run)
	return _c
}

// RaiseEvent provides a mock function with given fields: event
func (_m *DcpNozzleIface) RaiseEvent(event *common.Event) {
	_m.Called(event)
}

// DcpNozzleIface_RaiseEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RaiseEvent'
type DcpNozzleIface_RaiseEvent_Call struct {
	*mock.Call
}

// RaiseEvent is a helper method to define mock.On call
//   - event *common.Event
func (_e *DcpNozzleIface_Expecter) RaiseEvent(event interface{}) *DcpNozzleIface_RaiseEvent_Call {
	return &DcpNozzleIface_RaiseEvent_Call{Call: _e.mock.On("RaiseEvent", event)}
}

func (_c *DcpNozzleIface_RaiseEvent_Call) Run(run func(event *common.Event)) *DcpNozzleIface_RaiseEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*common.Event))
	})
	return _c
}

func (_c *DcpNozzleIface_RaiseEvent_Call) Return() *DcpNozzleIface_RaiseEvent_Call {
	_c.Call.Return()
	return _c
}

func (_c *DcpNozzleIface_RaiseEvent_Call) RunAndReturn(run func(*common.Event)) *DcpNozzleIface_RaiseEvent_Call {
	_c.Call.Return(run)
	return _c
}

// Receive provides a mock function with given fields: data
func (_m *DcpNozzleIface) Receive(data interface{}) error {
	ret := _m.Called(data)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_Receive_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Receive'
type DcpNozzleIface_Receive_Call struct {
	*mock.Call
}

// Receive is a helper method to define mock.On call
//   - data interface{}
func (_e *DcpNozzleIface_Expecter) Receive(data interface{}) *DcpNozzleIface_Receive_Call {
	return &DcpNozzleIface_Receive_Call{Call: _e.mock.On("Receive", data)}
}

func (_c *DcpNozzleIface_Receive_Call) Run(run func(data interface{})) *DcpNozzleIface_Receive_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *DcpNozzleIface_Receive_Call) Return(_a0 error) *DcpNozzleIface_Receive_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_Receive_Call) RunAndReturn(run func(interface{}) error) *DcpNozzleIface_Receive_Call {
	_c.Call.Return(run)
	return _c
}

// RecycleDataObj provides a mock function with given fields: obj
func (_m *DcpNozzleIface) RecycleDataObj(obj interface{}) {
	_m.Called(obj)
}

// DcpNozzleIface_RecycleDataObj_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecycleDataObj'
type DcpNozzleIface_RecycleDataObj_Call struct {
	*mock.Call
}

// RecycleDataObj is a helper method to define mock.On call
//   - obj interface{}
func (_e *DcpNozzleIface_Expecter) RecycleDataObj(obj interface{}) *DcpNozzleIface_RecycleDataObj_Call {
	return &DcpNozzleIface_RecycleDataObj_Call{Call: _e.mock.On("RecycleDataObj", obj)}
}

func (_c *DcpNozzleIface_RecycleDataObj_Call) Run(run func(obj interface{})) *DcpNozzleIface_RecycleDataObj_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interface{}))
	})
	return _c
}

func (_c *DcpNozzleIface_RecycleDataObj_Call) Return() *DcpNozzleIface_RecycleDataObj_Call {
	_c.Call.Return()
	return _c
}

func (_c *DcpNozzleIface_RecycleDataObj_Call) RunAndReturn(run func(interface{})) *DcpNozzleIface_RecycleDataObj_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterComponentEventListener provides a mock function with given fields: eventType, listener
func (_m *DcpNozzleIface) RegisterComponentEventListener(eventType common.ComponentEventType, listener common.ComponentEventListener) error {
	ret := _m.Called(eventType, listener)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.ComponentEventType, common.ComponentEventListener) error); ok {
		r0 = rf(eventType, listener)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_RegisterComponentEventListener_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterComponentEventListener'
type DcpNozzleIface_RegisterComponentEventListener_Call struct {
	*mock.Call
}

// RegisterComponentEventListener is a helper method to define mock.On call
//   - eventType common.ComponentEventType
//   - listener common.ComponentEventListener
func (_e *DcpNozzleIface_Expecter) RegisterComponentEventListener(eventType interface{}, listener interface{}) *DcpNozzleIface_RegisterComponentEventListener_Call {
	return &DcpNozzleIface_RegisterComponentEventListener_Call{Call: _e.mock.On("RegisterComponentEventListener", eventType, listener)}
}

func (_c *DcpNozzleIface_RegisterComponentEventListener_Call) Run(run func(eventType common.ComponentEventType, listener common.ComponentEventListener)) *DcpNozzleIface_RegisterComponentEventListener_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.ComponentEventType), args[1].(common.ComponentEventListener))
	})
	return _c
}

func (_c *DcpNozzleIface_RegisterComponentEventListener_Call) Return(_a0 error) *DcpNozzleIface_RegisterComponentEventListener_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_RegisterComponentEventListener_Call) RunAndReturn(run func(common.ComponentEventType, common.ComponentEventListener) error) *DcpNozzleIface_RegisterComponentEventListener_Call {
	_c.Call.Return(run)
	return _c
}

// ResponsibleVBs provides a mock function with given fields:
func (_m *DcpNozzleIface) ResponsibleVBs() []uint16 {
	ret := _m.Called()

	var r0 []uint16
	if rf, ok := ret.Get(0).(func() []uint16); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint16)
		}
	}

	return r0
}

// DcpNozzleIface_ResponsibleVBs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResponsibleVBs'
type DcpNozzleIface_ResponsibleVBs_Call struct {
	*mock.Call
}

// ResponsibleVBs is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) ResponsibleVBs() *DcpNozzleIface_ResponsibleVBs_Call {
	return &DcpNozzleIface_ResponsibleVBs_Call{Call: _e.mock.On("ResponsibleVBs")}
}

func (_c *DcpNozzleIface_ResponsibleVBs_Call) Run(run func()) *DcpNozzleIface_ResponsibleVBs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_ResponsibleVBs_Call) Return(_a0 []uint16) *DcpNozzleIface_ResponsibleVBs_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_ResponsibleVBs_Call) RunAndReturn(run func() []uint16) *DcpNozzleIface_ResponsibleVBs_Call {
	_c.Call.Return(run)
	return _c
}

// SetConnector provides a mock function with given fields: connector
func (_m *DcpNozzleIface) SetConnector(connector common.Connector) error {
	ret := _m.Called(connector)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Connector) error); ok {
		r0 = rf(connector)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_SetConnector_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetConnector'
type DcpNozzleIface_SetConnector_Call struct {
	*mock.Call
}

// SetConnector is a helper method to define mock.On call
//   - connector common.Connector
func (_e *DcpNozzleIface_Expecter) SetConnector(connector interface{}) *DcpNozzleIface_SetConnector_Call {
	return &DcpNozzleIface_SetConnector_Call{Call: _e.mock.On("SetConnector", connector)}
}

func (_c *DcpNozzleIface_SetConnector_Call) Run(run func(connector common.Connector)) *DcpNozzleIface_SetConnector_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Connector))
	})
	return _c
}

func (_c *DcpNozzleIface_SetConnector_Call) Return(_a0 error) *DcpNozzleIface_SetConnector_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_SetConnector_Call) RunAndReturn(run func(common.Connector) error) *DcpNozzleIface_SetConnector_Call {
	_c.Call.Return(run)
	return _c
}

// SetMaxMissCount provides a mock function with given fields: max_dcp_miss_count
func (_m *DcpNozzleIface) SetMaxMissCount(max_dcp_miss_count int) {
	_m.Called(max_dcp_miss_count)
}

// DcpNozzleIface_SetMaxMissCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetMaxMissCount'
type DcpNozzleIface_SetMaxMissCount_Call struct {
	*mock.Call
}

// SetMaxMissCount is a helper method to define mock.On call
//   - max_dcp_miss_count int
func (_e *DcpNozzleIface_Expecter) SetMaxMissCount(max_dcp_miss_count interface{}) *DcpNozzleIface_SetMaxMissCount_Call {
	return &DcpNozzleIface_SetMaxMissCount_Call{Call: _e.mock.On("SetMaxMissCount", max_dcp_miss_count)}
}

func (_c *DcpNozzleIface_SetMaxMissCount_Call) Run(run func(max_dcp_miss_count int)) *DcpNozzleIface_SetMaxMissCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *DcpNozzleIface_SetMaxMissCount_Call) Return() *DcpNozzleIface_SetMaxMissCount_Call {
	_c.Call.Return()
	return _c
}

func (_c *DcpNozzleIface_SetMaxMissCount_Call) RunAndReturn(run func(int)) *DcpNozzleIface_SetMaxMissCount_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields: settings
func (_m *DcpNozzleIface) Start(settings metadata.ReplicationSettingsMap) error {
	ret := _m.Called(settings)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) error); ok {
		r0 = rf(settings)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type DcpNozzleIface_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
//   - settings metadata.ReplicationSettingsMap
func (_e *DcpNozzleIface_Expecter) Start(settings interface{}) *DcpNozzleIface_Start_Call {
	return &DcpNozzleIface_Start_Call{Call: _e.mock.On("Start", settings)}
}

func (_c *DcpNozzleIface_Start_Call) Run(run func(settings metadata.ReplicationSettingsMap)) *DcpNozzleIface_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.ReplicationSettingsMap))
	})
	return _c
}

func (_c *DcpNozzleIface_Start_Call) Return(_a0 error) *DcpNozzleIface_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_Start_Call) RunAndReturn(run func(metadata.ReplicationSettingsMap) error) *DcpNozzleIface_Start_Call {
	_c.Call.Return(run)
	return _c
}

// State provides a mock function with given fields:
func (_m *DcpNozzleIface) State() common.PartState {
	ret := _m.Called()

	var r0 common.PartState
	if rf, ok := ret.Get(0).(func() common.PartState); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(common.PartState)
	}

	return r0
}

// DcpNozzleIface_State_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'State'
type DcpNozzleIface_State_Call struct {
	*mock.Call
}

// State is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) State() *DcpNozzleIface_State_Call {
	return &DcpNozzleIface_State_Call{Call: _e.mock.On("State")}
}

func (_c *DcpNozzleIface_State_Call) Run(run func()) *DcpNozzleIface_State_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_State_Call) Return(_a0 common.PartState) *DcpNozzleIface_State_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_State_Call) RunAndReturn(run func() common.PartState) *DcpNozzleIface_State_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *DcpNozzleIface) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type DcpNozzleIface_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *DcpNozzleIface_Expecter) Stop() *DcpNozzleIface_Stop_Call {
	return &DcpNozzleIface_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *DcpNozzleIface_Stop_Call) Run(run func()) *DcpNozzleIface_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DcpNozzleIface_Stop_Call) Return(_a0 error) *DcpNozzleIface_Stop_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_Stop_Call) RunAndReturn(run func() error) *DcpNozzleIface_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// UnRegisterComponentEventListener provides a mock function with given fields: eventType, listener
func (_m *DcpNozzleIface) UnRegisterComponentEventListener(eventType common.ComponentEventType, listener common.ComponentEventListener) error {
	ret := _m.Called(eventType, listener)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.ComponentEventType, common.ComponentEventListener) error); ok {
		r0 = rf(eventType, listener)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_UnRegisterComponentEventListener_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnRegisterComponentEventListener'
type DcpNozzleIface_UnRegisterComponentEventListener_Call struct {
	*mock.Call
}

// UnRegisterComponentEventListener is a helper method to define mock.On call
//   - eventType common.ComponentEventType
//   - listener common.ComponentEventListener
func (_e *DcpNozzleIface_Expecter) UnRegisterComponentEventListener(eventType interface{}, listener interface{}) *DcpNozzleIface_UnRegisterComponentEventListener_Call {
	return &DcpNozzleIface_UnRegisterComponentEventListener_Call{Call: _e.mock.On("UnRegisterComponentEventListener", eventType, listener)}
}

func (_c *DcpNozzleIface_UnRegisterComponentEventListener_Call) Run(run func(eventType common.ComponentEventType, listener common.ComponentEventListener)) *DcpNozzleIface_UnRegisterComponentEventListener_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.ComponentEventType), args[1].(common.ComponentEventListener))
	})
	return _c
}

func (_c *DcpNozzleIface_UnRegisterComponentEventListener_Call) Return(_a0 error) *DcpNozzleIface_UnRegisterComponentEventListener_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_UnRegisterComponentEventListener_Call) RunAndReturn(run func(common.ComponentEventType, common.ComponentEventListener) error) *DcpNozzleIface_UnRegisterComponentEventListener_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSettings provides a mock function with given fields: settings
func (_m *DcpNozzleIface) UpdateSettings(settings metadata.ReplicationSettingsMap) error {
	ret := _m.Called(settings)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.ReplicationSettingsMap) error); ok {
		r0 = rf(settings)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DcpNozzleIface_UpdateSettings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSettings'
type DcpNozzleIface_UpdateSettings_Call struct {
	*mock.Call
}

// UpdateSettings is a helper method to define mock.On call
//   - settings metadata.ReplicationSettingsMap
func (_e *DcpNozzleIface_Expecter) UpdateSettings(settings interface{}) *DcpNozzleIface_UpdateSettings_Call {
	return &DcpNozzleIface_UpdateSettings_Call{Call: _e.mock.On("UpdateSettings", settings)}
}

func (_c *DcpNozzleIface_UpdateSettings_Call) Run(run func(settings metadata.ReplicationSettingsMap)) *DcpNozzleIface_UpdateSettings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(metadata.ReplicationSettingsMap))
	})
	return _c
}

func (_c *DcpNozzleIface_UpdateSettings_Call) Return(_a0 error) *DcpNozzleIface_UpdateSettings_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DcpNozzleIface_UpdateSettings_Call) RunAndReturn(run func(metadata.ReplicationSettingsMap) error) *DcpNozzleIface_UpdateSettings_Call {
	_c.Call.Return(run)
	return _c
}

// NewDcpNozzleIface creates a new instance of DcpNozzleIface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDcpNozzleIface(t interface {
	mock.TestingT
	Cleanup(func())
}) *DcpNozzleIface {
	mock := &DcpNozzleIface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}