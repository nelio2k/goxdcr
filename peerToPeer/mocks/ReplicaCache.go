// Code generated by mockery (devel). DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"
)

// ReplicaCache is an autogenerated mock type for the ReplicaCache type
type ReplicaCache struct {
	mock.Mock
}

// GetReplicaInfo provides a mock function with given fields: spec
func (_m *ReplicaCache) GetReplicaInfo(spec *metadata.ReplicationSpecification) (int, base.VbHostsMapType, base.StringStringMap, func(), error) {
	ret := _m.Called(spec)

	var r0 int
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) int); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 base.VbHostsMapType
	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification) base.VbHostsMapType); ok {
		r1 = rf(spec)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(base.VbHostsMapType)
		}
	}

	var r2 base.StringStringMap
	if rf, ok := ret.Get(2).(func(*metadata.ReplicationSpecification) base.StringStringMap); ok {
		r2 = rf(spec)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(base.StringStringMap)
		}
	}

	var r3 func()
	if rf, ok := ret.Get(3).(func(*metadata.ReplicationSpecification) func()); ok {
		r3 = rf(spec)
	} else {
		if ret.Get(3) != nil {
			r3 = ret.Get(3).(func())
		}
	}

	var r4 error
	if rf, ok := ret.Get(4).(func(*metadata.ReplicationSpecification) error); ok {
		r4 = rf(spec)
	} else {
		r4 = ret.Error(4)
	}

	return r0, r1, r2, r3, r4
}

// HandleSpecCreation provides a mock function with given fields: spec
func (_m *ReplicaCache) HandleSpecCreation(spec *metadata.ReplicationSpecification) {
	_m.Called(spec)
}

// HandleSpecDeletion provides a mock function with given fields: oldSpec
func (_m *ReplicaCache) HandleSpecDeletion(oldSpec *metadata.ReplicationSpecification) {
	_m.Called(oldSpec)
}
