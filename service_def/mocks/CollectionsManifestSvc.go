// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"

	sync "sync"
)

// CollectionsManifestSvc is an autogenerated mock type for the CollectionsManifestSvc type
type CollectionsManifestSvc struct {
	mock.Mock
}

// CollectionManifestGetter provides a mock function with given fields: bucketName, hasStoredManifest, storedManifestUid, spec
func (_m *CollectionsManifestSvc) CollectionManifestGetter(bucketName string, hasStoredManifest bool, storedManifestUid uint64, spec *metadata.ReplicationSpecification) (*metadata.CollectionsManifest, error) {
	ret := _m.Called(bucketName, hasStoredManifest, storedManifestUid, spec)

	var r0 *metadata.CollectionsManifest
	if rf, ok := ret.Get(0).(func(string, bool, uint64, *metadata.ReplicationSpecification) *metadata.CollectionsManifest); ok {
		r0 = rf(bucketName, hasStoredManifest, storedManifestUid, spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.CollectionsManifest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool, uint64, *metadata.ReplicationSpecification) error); ok {
		r1 = rf(bucketName, hasStoredManifest, storedManifestUid, spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ForceTargetManifestRefresh provides a mock function with given fields: spec
func (_m *CollectionsManifestSvc) ForceTargetManifestRefresh(spec *metadata.ReplicationSpecification) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllCachedManifests provides a mock function with given fields: spec
func (_m *CollectionsManifestSvc) GetAllCachedManifests(spec *metadata.ReplicationSpecification) (map[uint64]*metadata.CollectionsManifest, map[uint64]*metadata.CollectionsManifest, error) {
	ret := _m.Called(spec)

	var r0 map[uint64]*metadata.CollectionsManifest
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) map[uint64]*metadata.CollectionsManifest); ok {
		r0 = rf(spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uint64]*metadata.CollectionsManifest)
		}
	}

	var r1 map[uint64]*metadata.CollectionsManifest
	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification) map[uint64]*metadata.CollectionsManifest); ok {
		r1 = rf(spec)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[uint64]*metadata.CollectionsManifest)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*metadata.ReplicationSpecification) error); ok {
		r2 = rf(spec)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetLastPersistedManifests provides a mock function with given fields: spec
func (_m *CollectionsManifestSvc) GetLastPersistedManifests(spec *metadata.ReplicationSpecification) (*metadata.CollectionsManifestPair, error) {
	ret := _m.Called(spec)

	var r0 *metadata.CollectionsManifestPair
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) *metadata.CollectionsManifestPair); ok {
		r0 = rf(spec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.CollectionsManifestPair)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification) error); ok {
		r1 = rf(spec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestManifests provides a mock function with given fields: spec, specMayNotExist
func (_m *CollectionsManifestSvc) GetLatestManifests(spec *metadata.ReplicationSpecification, specMayNotExist bool) (*metadata.CollectionsManifest, *metadata.CollectionsManifest, error) {
	ret := _m.Called(spec, specMayNotExist)

	var r0 *metadata.CollectionsManifest
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification, bool) *metadata.CollectionsManifest); ok {
		r0 = rf(spec, specMayNotExist)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.CollectionsManifest)
		}
	}

	var r1 *metadata.CollectionsManifest
	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification, bool) *metadata.CollectionsManifest); ok {
		r1 = rf(spec, specMayNotExist)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*metadata.CollectionsManifest)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*metadata.ReplicationSpecification, bool) error); ok {
		r2 = rf(spec, specMayNotExist)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetSpecificSourceManifest provides a mock function with given fields: spec, manifestVersion
func (_m *CollectionsManifestSvc) GetSpecificSourceManifest(spec *metadata.ReplicationSpecification, manifestVersion uint64) (*metadata.CollectionsManifest, error) {
	ret := _m.Called(spec, manifestVersion)

	var r0 *metadata.CollectionsManifest
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification, uint64) *metadata.CollectionsManifest); ok {
		r0 = rf(spec, manifestVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.CollectionsManifest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification, uint64) error); ok {
		r1 = rf(spec, manifestVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSpecificTargetManifest provides a mock function with given fields: spec, manifestVersion
func (_m *CollectionsManifestSvc) GetSpecificTargetManifest(spec *metadata.ReplicationSpecification, manifestVersion uint64) (*metadata.CollectionsManifest, error) {
	ret := _m.Called(spec, manifestVersion)

	var r0 *metadata.CollectionsManifest
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification, uint64) *metadata.CollectionsManifest); ok {
		r0 = rf(spec, manifestVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.CollectionsManifest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*metadata.ReplicationSpecification, uint64) error); ok {
		r1 = rf(spec, manifestVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PersistNeededManifests provides a mock function with given fields: spec
func (_m *CollectionsManifestSvc) PersistNeededManifests(spec *metadata.ReplicationSpecification) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PersistReceivedManifests provides a mock function with given fields: spec, srcManifests, tgtManifests
func (_m *CollectionsManifestSvc) PersistReceivedManifests(spec *metadata.ReplicationSpecification, srcManifests map[uint64]*metadata.CollectionsManifest, tgtManifests map[uint64]*metadata.CollectionsManifest) error {
	ret := _m.Called(spec, srcManifests, tgtManifests)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification, map[uint64]*metadata.CollectionsManifest, map[uint64]*metadata.CollectionsManifest) error); ok {
		r0 = rf(spec, srcManifests, tgtManifests)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReplicationSpecChangeCallback provides a mock function with given fields: id, oldVal, newVal, wg
func (_m *CollectionsManifestSvc) ReplicationSpecChangeCallback(id string, oldVal interface{}, newVal interface{}, wg *sync.WaitGroup) error {
	ret := _m.Called(id, oldVal, newVal, wg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}, *sync.WaitGroup) error); ok {
		r0 = rf(id, oldVal, newVal, wg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetMetadataChangeHandlerCallback provides a mock function with given fields: callBack
func (_m *CollectionsManifestSvc) SetMetadataChangeHandlerCallback(callBack base.MetadataChangeHandlerCallback) {
	_m.Called(callBack)
}

type mockConstructorTestingTNewCollectionsManifestSvc interface {
	mock.TestingT
	Cleanup(func())
}

// NewCollectionsManifestSvc creates a new instance of CollectionsManifestSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCollectionsManifestSvc(t mockConstructorTestingTNewCollectionsManifestSvc) *CollectionsManifestSvc {
	mock := &CollectionsManifestSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
