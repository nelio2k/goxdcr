package mocks

import common "github.com/couchbase/goxdcr/common"
import mock "github.com/stretchr/testify/mock"

// ComponentEventListener is an autogenerated mock type for the ComponentEventListener type
type ComponentEventListener struct {
	mock.Mock
}

// OnEvent provides a mock function with given fields: event
func (_m *ComponentEventListener) OnEvent(event *common.Event) {
	_m.Called(event)
}