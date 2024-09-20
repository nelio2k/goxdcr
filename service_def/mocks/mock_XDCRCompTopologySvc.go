// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/v8/base"
	mock "github.com/stretchr/testify/mock"
)

// XDCRCompTopologySvc is an autogenerated mock type for the XDCRCompTopologySvc type
type XDCRCompTopologySvc struct {
	mock.Mock
}

type XDCRCompTopologySvc_Expecter struct {
	mock *mock.Mock
}

func (_m *XDCRCompTopologySvc) EXPECT() *XDCRCompTopologySvc_Expecter {
	return &XDCRCompTopologySvc_Expecter{mock: &_m.Mock}
}

// ClientCertIsMandatory provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) ClientCertIsMandatory() (bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ClientCertIsMandatory")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_ClientCertIsMandatory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ClientCertIsMandatory'
type XDCRCompTopologySvc_ClientCertIsMandatory_Call struct {
	*mock.Call
}

// ClientCertIsMandatory is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) ClientCertIsMandatory() *XDCRCompTopologySvc_ClientCertIsMandatory_Call {
	return &XDCRCompTopologySvc_ClientCertIsMandatory_Call{Call: _e.mock.On("ClientCertIsMandatory")}
}

func (_c *XDCRCompTopologySvc_ClientCertIsMandatory_Call) Run(run func()) *XDCRCompTopologySvc_ClientCertIsMandatory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_ClientCertIsMandatory_Call) Return(_a0 bool, _a1 error) *XDCRCompTopologySvc_ClientCertIsMandatory_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_ClientCertIsMandatory_Call) RunAndReturn(run func() (bool, error)) *XDCRCompTopologySvc_ClientCertIsMandatory_Call {
	_c.Call.Return(run)
	return _c
}

// GetLocalHostName provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) GetLocalHostName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetLocalHostName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// XDCRCompTopologySvc_GetLocalHostName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLocalHostName'
type XDCRCompTopologySvc_GetLocalHostName_Call struct {
	*mock.Call
}

// GetLocalHostName is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) GetLocalHostName() *XDCRCompTopologySvc_GetLocalHostName_Call {
	return &XDCRCompTopologySvc_GetLocalHostName_Call{Call: _e.mock.On("GetLocalHostName")}
}

func (_c *XDCRCompTopologySvc_GetLocalHostName_Call) Run(run func()) *XDCRCompTopologySvc_GetLocalHostName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_GetLocalHostName_Call) Return(_a0 string) *XDCRCompTopologySvc_GetLocalHostName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *XDCRCompTopologySvc_GetLocalHostName_Call) RunAndReturn(run func() string) *XDCRCompTopologySvc_GetLocalHostName_Call {
	_c.Call.Return(run)
	return _c
}

// IsIpv4Blocked provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) IsIpv4Blocked() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsIpv4Blocked")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// XDCRCompTopologySvc_IsIpv4Blocked_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsIpv4Blocked'
type XDCRCompTopologySvc_IsIpv4Blocked_Call struct {
	*mock.Call
}

// IsIpv4Blocked is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) IsIpv4Blocked() *XDCRCompTopologySvc_IsIpv4Blocked_Call {
	return &XDCRCompTopologySvc_IsIpv4Blocked_Call{Call: _e.mock.On("IsIpv4Blocked")}
}

func (_c *XDCRCompTopologySvc_IsIpv4Blocked_Call) Run(run func()) *XDCRCompTopologySvc_IsIpv4Blocked_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_IsIpv4Blocked_Call) Return(_a0 bool) *XDCRCompTopologySvc_IsIpv4Blocked_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *XDCRCompTopologySvc_IsIpv4Blocked_Call) RunAndReturn(run func() bool) *XDCRCompTopologySvc_IsIpv4Blocked_Call {
	_c.Call.Return(run)
	return _c
}

// IsIpv6Blocked provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) IsIpv6Blocked() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsIpv6Blocked")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// XDCRCompTopologySvc_IsIpv6Blocked_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsIpv6Blocked'
type XDCRCompTopologySvc_IsIpv6Blocked_Call struct {
	*mock.Call
}

// IsIpv6Blocked is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) IsIpv6Blocked() *XDCRCompTopologySvc_IsIpv6Blocked_Call {
	return &XDCRCompTopologySvc_IsIpv6Blocked_Call{Call: _e.mock.On("IsIpv6Blocked")}
}

func (_c *XDCRCompTopologySvc_IsIpv6Blocked_Call) Run(run func()) *XDCRCompTopologySvc_IsIpv6Blocked_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_IsIpv6Blocked_Call) Return(_a0 bool) *XDCRCompTopologySvc_IsIpv6Blocked_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *XDCRCompTopologySvc_IsIpv6Blocked_Call) RunAndReturn(run func() bool) *XDCRCompTopologySvc_IsIpv6Blocked_Call {
	_c.Call.Return(run)
	return _c
}

// IsKVNode provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) IsKVNode() (bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsKVNode")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_IsKVNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsKVNode'
type XDCRCompTopologySvc_IsKVNode_Call struct {
	*mock.Call
}

// IsKVNode is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) IsKVNode() *XDCRCompTopologySvc_IsKVNode_Call {
	return &XDCRCompTopologySvc_IsKVNode_Call{Call: _e.mock.On("IsKVNode")}
}

func (_c *XDCRCompTopologySvc_IsKVNode_Call) Run(run func()) *XDCRCompTopologySvc_IsKVNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_IsKVNode_Call) Return(_a0 bool, _a1 error) *XDCRCompTopologySvc_IsKVNode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_IsKVNode_Call) RunAndReturn(run func() (bool, error)) *XDCRCompTopologySvc_IsKVNode_Call {
	_c.Call.Return(run)
	return _c
}

// IsMyClusterDeveloperPreview provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) IsMyClusterDeveloperPreview() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsMyClusterDeveloperPreview")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsMyClusterDeveloperPreview'
type XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call struct {
	*mock.Call
}

// IsMyClusterDeveloperPreview is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) IsMyClusterDeveloperPreview() *XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call {
	return &XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call{Call: _e.mock.On("IsMyClusterDeveloperPreview")}
}

func (_c *XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call) Run(run func()) *XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call) Return(_a0 bool) *XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call) RunAndReturn(run func() bool) *XDCRCompTopologySvc_IsMyClusterDeveloperPreview_Call {
	_c.Call.Return(run)
	return _c
}

// IsMyClusterEncryptionLevelStrict provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) IsMyClusterEncryptionLevelStrict() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsMyClusterEncryptionLevelStrict")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsMyClusterEncryptionLevelStrict'
type XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call struct {
	*mock.Call
}

// IsMyClusterEncryptionLevelStrict is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) IsMyClusterEncryptionLevelStrict() *XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call {
	return &XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call{Call: _e.mock.On("IsMyClusterEncryptionLevelStrict")}
}

func (_c *XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call) Run(run func()) *XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call) Return(_a0 bool) *XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call) RunAndReturn(run func() bool) *XDCRCompTopologySvc_IsMyClusterEncryptionLevelStrict_Call {
	_c.Call.Return(run)
	return _c
}

// IsMyClusterEnterprise provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) IsMyClusterEnterprise() (bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsMyClusterEnterprise")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_IsMyClusterEnterprise_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsMyClusterEnterprise'
type XDCRCompTopologySvc_IsMyClusterEnterprise_Call struct {
	*mock.Call
}

// IsMyClusterEnterprise is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) IsMyClusterEnterprise() *XDCRCompTopologySvc_IsMyClusterEnterprise_Call {
	return &XDCRCompTopologySvc_IsMyClusterEnterprise_Call{Call: _e.mock.On("IsMyClusterEnterprise")}
}

func (_c *XDCRCompTopologySvc_IsMyClusterEnterprise_Call) Run(run func()) *XDCRCompTopologySvc_IsMyClusterEnterprise_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_IsMyClusterEnterprise_Call) Return(_a0 bool, _a1 error) *XDCRCompTopologySvc_IsMyClusterEnterprise_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_IsMyClusterEnterprise_Call) RunAndReturn(run func() (bool, error)) *XDCRCompTopologySvc_IsMyClusterEnterprise_Call {
	_c.Call.Return(run)
	return _c
}

// IsMyClusterIpv6 provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) IsMyClusterIpv6() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsMyClusterIpv6")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// XDCRCompTopologySvc_IsMyClusterIpv6_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsMyClusterIpv6'
type XDCRCompTopologySvc_IsMyClusterIpv6_Call struct {
	*mock.Call
}

// IsMyClusterIpv6 is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) IsMyClusterIpv6() *XDCRCompTopologySvc_IsMyClusterIpv6_Call {
	return &XDCRCompTopologySvc_IsMyClusterIpv6_Call{Call: _e.mock.On("IsMyClusterIpv6")}
}

func (_c *XDCRCompTopologySvc_IsMyClusterIpv6_Call) Run(run func()) *XDCRCompTopologySvc_IsMyClusterIpv6_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_IsMyClusterIpv6_Call) Return(_a0 bool) *XDCRCompTopologySvc_IsMyClusterIpv6_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *XDCRCompTopologySvc_IsMyClusterIpv6_Call) RunAndReturn(run func() bool) *XDCRCompTopologySvc_IsMyClusterIpv6_Call {
	_c.Call.Return(run)
	return _c
}

// MyAdminPort provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyAdminPort() (uint16, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyAdminPort")
	}

	var r0 uint16
	var r1 error
	if rf, ok := ret.Get(0).(func() (uint16, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() uint16); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint16)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyAdminPort_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyAdminPort'
type XDCRCompTopologySvc_MyAdminPort_Call struct {
	*mock.Call
}

// MyAdminPort is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyAdminPort() *XDCRCompTopologySvc_MyAdminPort_Call {
	return &XDCRCompTopologySvc_MyAdminPort_Call{Call: _e.mock.On("MyAdminPort")}
}

func (_c *XDCRCompTopologySvc_MyAdminPort_Call) Run(run func()) *XDCRCompTopologySvc_MyAdminPort_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyAdminPort_Call) Return(_a0 uint16, _a1 error) *XDCRCompTopologySvc_MyAdminPort_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyAdminPort_Call) RunAndReturn(run func() (uint16, error)) *XDCRCompTopologySvc_MyAdminPort_Call {
	_c.Call.Return(run)
	return _c
}

// MyClusterCompatibility provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyClusterCompatibility() (int, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyClusterCompatibility")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func() (int, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyClusterCompatibility_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyClusterCompatibility'
type XDCRCompTopologySvc_MyClusterCompatibility_Call struct {
	*mock.Call
}

// MyClusterCompatibility is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyClusterCompatibility() *XDCRCompTopologySvc_MyClusterCompatibility_Call {
	return &XDCRCompTopologySvc_MyClusterCompatibility_Call{Call: _e.mock.On("MyClusterCompatibility")}
}

func (_c *XDCRCompTopologySvc_MyClusterCompatibility_Call) Run(run func()) *XDCRCompTopologySvc_MyClusterCompatibility_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyClusterCompatibility_Call) Return(_a0 int, _a1 error) *XDCRCompTopologySvc_MyClusterCompatibility_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyClusterCompatibility_Call) RunAndReturn(run func() (int, error)) *XDCRCompTopologySvc_MyClusterCompatibility_Call {
	_c.Call.Return(run)
	return _c
}

// MyClusterUuid provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyClusterUuid() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyClusterUuid")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyClusterUuid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyClusterUuid'
type XDCRCompTopologySvc_MyClusterUuid_Call struct {
	*mock.Call
}

// MyClusterUuid is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyClusterUuid() *XDCRCompTopologySvc_MyClusterUuid_Call {
	return &XDCRCompTopologySvc_MyClusterUuid_Call{Call: _e.mock.On("MyClusterUuid")}
}

func (_c *XDCRCompTopologySvc_MyClusterUuid_Call) Run(run func()) *XDCRCompTopologySvc_MyClusterUuid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyClusterUuid_Call) Return(_a0 string, _a1 error) *XDCRCompTopologySvc_MyClusterUuid_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyClusterUuid_Call) RunAndReturn(run func() (string, error)) *XDCRCompTopologySvc_MyClusterUuid_Call {
	_c.Call.Return(run)
	return _c
}

// MyConnectionStr provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyConnectionStr() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyConnectionStr")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyConnectionStr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyConnectionStr'
type XDCRCompTopologySvc_MyConnectionStr_Call struct {
	*mock.Call
}

// MyConnectionStr is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyConnectionStr() *XDCRCompTopologySvc_MyConnectionStr_Call {
	return &XDCRCompTopologySvc_MyConnectionStr_Call{Call: _e.mock.On("MyConnectionStr")}
}

func (_c *XDCRCompTopologySvc_MyConnectionStr_Call) Run(run func()) *XDCRCompTopologySvc_MyConnectionStr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyConnectionStr_Call) Return(_a0 string, _a1 error) *XDCRCompTopologySvc_MyConnectionStr_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyConnectionStr_Call) RunAndReturn(run func() (string, error)) *XDCRCompTopologySvc_MyConnectionStr_Call {
	_c.Call.Return(run)
	return _c
}

// MyCredentials provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyCredentials() (string, string, base.HttpAuthMech, []byte, bool, []byte, []byte, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyCredentials")
	}

	var r0 string
	var r1 string
	var r2 base.HttpAuthMech
	var r3 []byte
	var r4 bool
	var r5 []byte
	var r6 []byte
	var r7 error
	if rf, ok := ret.Get(0).(func() (string, string, base.HttpAuthMech, []byte, bool, []byte, []byte, error)); ok {
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

	if rf, ok := ret.Get(2).(func() base.HttpAuthMech); ok {
		r2 = rf()
	} else {
		r2 = ret.Get(2).(base.HttpAuthMech)
	}

	if rf, ok := ret.Get(3).(func() []byte); ok {
		r3 = rf()
	} else {
		if ret.Get(3) != nil {
			r3 = ret.Get(3).([]byte)
		}
	}

	if rf, ok := ret.Get(4).(func() bool); ok {
		r4 = rf()
	} else {
		r4 = ret.Get(4).(bool)
	}

	if rf, ok := ret.Get(5).(func() []byte); ok {
		r5 = rf()
	} else {
		if ret.Get(5) != nil {
			r5 = ret.Get(5).([]byte)
		}
	}

	if rf, ok := ret.Get(6).(func() []byte); ok {
		r6 = rf()
	} else {
		if ret.Get(6) != nil {
			r6 = ret.Get(6).([]byte)
		}
	}

	if rf, ok := ret.Get(7).(func() error); ok {
		r7 = rf()
	} else {
		r7 = ret.Error(7)
	}

	return r0, r1, r2, r3, r4, r5, r6, r7
}

// XDCRCompTopologySvc_MyCredentials_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyCredentials'
type XDCRCompTopologySvc_MyCredentials_Call struct {
	*mock.Call
}

// MyCredentials is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyCredentials() *XDCRCompTopologySvc_MyCredentials_Call {
	return &XDCRCompTopologySvc_MyCredentials_Call{Call: _e.mock.On("MyCredentials")}
}

func (_c *XDCRCompTopologySvc_MyCredentials_Call) Run(run func()) *XDCRCompTopologySvc_MyCredentials_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyCredentials_Call) Return(_a0 string, _a1 string, _a2 base.HttpAuthMech, _a3 []byte, _a4 bool, _a5 []byte, _a6 []byte, _a7 error) *XDCRCompTopologySvc_MyCredentials_Call {
	_c.Call.Return(_a0, _a1, _a2, _a3, _a4, _a5, _a6, _a7)
	return _c
}

func (_c *XDCRCompTopologySvc_MyCredentials_Call) RunAndReturn(run func() (string, string, base.HttpAuthMech, []byte, bool, []byte, []byte, error)) *XDCRCompTopologySvc_MyCredentials_Call {
	_c.Call.Return(run)
	return _c
}

// MyHost provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyHost() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyHost")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyHost_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyHost'
type XDCRCompTopologySvc_MyHost_Call struct {
	*mock.Call
}

// MyHost is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyHost() *XDCRCompTopologySvc_MyHost_Call {
	return &XDCRCompTopologySvc_MyHost_Call{Call: _e.mock.On("MyHost")}
}

func (_c *XDCRCompTopologySvc_MyHost_Call) Run(run func()) *XDCRCompTopologySvc_MyHost_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyHost_Call) Return(_a0 string, _a1 error) *XDCRCompTopologySvc_MyHost_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyHost_Call) RunAndReturn(run func() (string, error)) *XDCRCompTopologySvc_MyHost_Call {
	_c.Call.Return(run)
	return _c
}

// MyHostAddr provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyHostAddr() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyHostAddr")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyHostAddr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyHostAddr'
type XDCRCompTopologySvc_MyHostAddr_Call struct {
	*mock.Call
}

// MyHostAddr is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyHostAddr() *XDCRCompTopologySvc_MyHostAddr_Call {
	return &XDCRCompTopologySvc_MyHostAddr_Call{Call: _e.mock.On("MyHostAddr")}
}

func (_c *XDCRCompTopologySvc_MyHostAddr_Call) Run(run func()) *XDCRCompTopologySvc_MyHostAddr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyHostAddr_Call) Return(_a0 string, _a1 error) *XDCRCompTopologySvc_MyHostAddr_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyHostAddr_Call) RunAndReturn(run func() (string, error)) *XDCRCompTopologySvc_MyHostAddr_Call {
	_c.Call.Return(run)
	return _c
}

// MyKVNodes provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyKVNodes() ([]string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyKVNodes")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyKVNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyKVNodes'
type XDCRCompTopologySvc_MyKVNodes_Call struct {
	*mock.Call
}

// MyKVNodes is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyKVNodes() *XDCRCompTopologySvc_MyKVNodes_Call {
	return &XDCRCompTopologySvc_MyKVNodes_Call{Call: _e.mock.On("MyKVNodes")}
}

func (_c *XDCRCompTopologySvc_MyKVNodes_Call) Run(run func()) *XDCRCompTopologySvc_MyKVNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyKVNodes_Call) Return(_a0 []string, _a1 error) *XDCRCompTopologySvc_MyKVNodes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyKVNodes_Call) RunAndReturn(run func() ([]string, error)) *XDCRCompTopologySvc_MyKVNodes_Call {
	_c.Call.Return(run)
	return _c
}

// MyMemcachedAddr provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyMemcachedAddr() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyMemcachedAddr")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyMemcachedAddr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyMemcachedAddr'
type XDCRCompTopologySvc_MyMemcachedAddr_Call struct {
	*mock.Call
}

// MyMemcachedAddr is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyMemcachedAddr() *XDCRCompTopologySvc_MyMemcachedAddr_Call {
	return &XDCRCompTopologySvc_MyMemcachedAddr_Call{Call: _e.mock.On("MyMemcachedAddr")}
}

func (_c *XDCRCompTopologySvc_MyMemcachedAddr_Call) Run(run func()) *XDCRCompTopologySvc_MyMemcachedAddr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyMemcachedAddr_Call) Return(_a0 string, _a1 error) *XDCRCompTopologySvc_MyMemcachedAddr_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyMemcachedAddr_Call) RunAndReturn(run func() (string, error)) *XDCRCompTopologySvc_MyMemcachedAddr_Call {
	_c.Call.Return(run)
	return _c
}

// MyNodeVersion provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) MyNodeVersion() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for MyNodeVersion")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_MyNodeVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'MyNodeVersion'
type XDCRCompTopologySvc_MyNodeVersion_Call struct {
	*mock.Call
}

// MyNodeVersion is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) MyNodeVersion() *XDCRCompTopologySvc_MyNodeVersion_Call {
	return &XDCRCompTopologySvc_MyNodeVersion_Call{Call: _e.mock.On("MyNodeVersion")}
}

func (_c *XDCRCompTopologySvc_MyNodeVersion_Call) Run(run func()) *XDCRCompTopologySvc_MyNodeVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_MyNodeVersion_Call) Return(_a0 string, _a1 error) *XDCRCompTopologySvc_MyNodeVersion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_MyNodeVersion_Call) RunAndReturn(run func() (string, error)) *XDCRCompTopologySvc_MyNodeVersion_Call {
	_c.Call.Return(run)
	return _c
}

// NumberOfKVNodes provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) NumberOfKVNodes() (int, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for NumberOfKVNodes")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func() (int, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_NumberOfKVNodes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NumberOfKVNodes'
type XDCRCompTopologySvc_NumberOfKVNodes_Call struct {
	*mock.Call
}

// NumberOfKVNodes is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) NumberOfKVNodes() *XDCRCompTopologySvc_NumberOfKVNodes_Call {
	return &XDCRCompTopologySvc_NumberOfKVNodes_Call{Call: _e.mock.On("NumberOfKVNodes")}
}

func (_c *XDCRCompTopologySvc_NumberOfKVNodes_Call) Run(run func()) *XDCRCompTopologySvc_NumberOfKVNodes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_NumberOfKVNodes_Call) Return(_a0 int, _a1 error) *XDCRCompTopologySvc_NumberOfKVNodes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_NumberOfKVNodes_Call) RunAndReturn(run func() (int, error)) *XDCRCompTopologySvc_NumberOfKVNodes_Call {
	_c.Call.Return(run)
	return _c
}

// PeerNodesAdminAddrs provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) PeerNodesAdminAddrs() ([]string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PeerNodesAdminAddrs")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_PeerNodesAdminAddrs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PeerNodesAdminAddrs'
type XDCRCompTopologySvc_PeerNodesAdminAddrs_Call struct {
	*mock.Call
}

// PeerNodesAdminAddrs is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) PeerNodesAdminAddrs() *XDCRCompTopologySvc_PeerNodesAdminAddrs_Call {
	return &XDCRCompTopologySvc_PeerNodesAdminAddrs_Call{Call: _e.mock.On("PeerNodesAdminAddrs")}
}

func (_c *XDCRCompTopologySvc_PeerNodesAdminAddrs_Call) Run(run func()) *XDCRCompTopologySvc_PeerNodesAdminAddrs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_PeerNodesAdminAddrs_Call) Return(_a0 []string, _a1 error) *XDCRCompTopologySvc_PeerNodesAdminAddrs_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_PeerNodesAdminAddrs_Call) RunAndReturn(run func() ([]string, error)) *XDCRCompTopologySvc_PeerNodesAdminAddrs_Call {
	_c.Call.Return(run)
	return _c
}

// XDCRCompToKVNodeMap provides a mock function with given fields:
func (_m *XDCRCompTopologySvc) XDCRCompToKVNodeMap() (map[string][]string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for XDCRCompToKVNodeMap")
	}

	var r0 map[string][]string
	var r1 error
	if rf, ok := ret.Get(0).(func() (map[string][]string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() map[string][]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]string)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'XDCRCompToKVNodeMap'
type XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call struct {
	*mock.Call
}

// XDCRCompToKVNodeMap is a helper method to define mock.On call
func (_e *XDCRCompTopologySvc_Expecter) XDCRCompToKVNodeMap() *XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call {
	return &XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call{Call: _e.mock.On("XDCRCompToKVNodeMap")}
}

func (_c *XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call) Run(run func()) *XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call) Return(_a0 map[string][]string, _a1 error) *XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call) RunAndReturn(run func() (map[string][]string, error)) *XDCRCompTopologySvc_XDCRCompToKVNodeMap_Call {
	_c.Call.Return(run)
	return _c
}

// NewXDCRCompTopologySvc creates a new instance of XDCRCompTopologySvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewXDCRCompTopologySvc(t interface {
	mock.TestingT
	Cleanup(func())
}) *XDCRCompTopologySvc {
	mock := &XDCRCompTopologySvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
