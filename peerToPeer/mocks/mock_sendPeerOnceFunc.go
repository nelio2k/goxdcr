// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	peerToPeer "github.com/couchbase/goxdcr/v8/peerToPeer"
	mock "github.com/stretchr/testify/mock"
)

// sendPeerOnceFunc is an autogenerated mock type for the sendPeerOnceFunc type
type sendPeerOnceFunc struct {
	mock.Mock
}

type sendPeerOnceFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *sendPeerOnceFunc) EXPECT() *sendPeerOnceFunc_Expecter {
	return &sendPeerOnceFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: opCode, getReqFunc, cbOpts
func (_m *sendPeerOnceFunc) Execute(opCode peerToPeer.OpCode, getReqFunc peerToPeer.GetReqFunc, cbOpts *peerToPeer.SendOpts) (error, map[string]bool) {
	ret := _m.Called(opCode, getReqFunc, cbOpts)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	var r1 map[string]bool
	if rf, ok := ret.Get(0).(func(peerToPeer.OpCode, peerToPeer.GetReqFunc, *peerToPeer.SendOpts) (error, map[string]bool)); ok {
		return rf(opCode, getReqFunc, cbOpts)
	}
	if rf, ok := ret.Get(0).(func(peerToPeer.OpCode, peerToPeer.GetReqFunc, *peerToPeer.SendOpts) error); ok {
		r0 = rf(opCode, getReqFunc, cbOpts)
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func(peerToPeer.OpCode, peerToPeer.GetReqFunc, *peerToPeer.SendOpts) map[string]bool); ok {
		r1 = rf(opCode, getReqFunc, cbOpts)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]bool)
		}
	}

	return r0, r1
}

// sendPeerOnceFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type sendPeerOnceFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - opCode peerToPeer.OpCode
//   - getReqFunc peerToPeer.GetReqFunc
//   - cbOpts *peerToPeer.SendOpts
func (_e *sendPeerOnceFunc_Expecter) Execute(opCode interface{}, getReqFunc interface{}, cbOpts interface{}) *sendPeerOnceFunc_Execute_Call {
	return &sendPeerOnceFunc_Execute_Call{Call: _e.mock.On("Execute", opCode, getReqFunc, cbOpts)}
}

func (_c *sendPeerOnceFunc_Execute_Call) Run(run func(opCode peerToPeer.OpCode, getReqFunc peerToPeer.GetReqFunc, cbOpts *peerToPeer.SendOpts)) *sendPeerOnceFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(peerToPeer.OpCode), args[1].(peerToPeer.GetReqFunc), args[2].(*peerToPeer.SendOpts))
	})
	return _c
}

func (_c *sendPeerOnceFunc_Execute_Call) Return(_a0 error, _a1 map[string]bool) *sendPeerOnceFunc_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *sendPeerOnceFunc_Execute_Call) RunAndReturn(run func(peerToPeer.OpCode, peerToPeer.GetReqFunc, *peerToPeer.SendOpts) (error, map[string]bool)) *sendPeerOnceFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// newSendPeerOnceFunc creates a new instance of sendPeerOnceFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newSendPeerOnceFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *sendPeerOnceFunc {
	mock := &sendPeerOnceFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}