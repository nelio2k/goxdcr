package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"
)

// RemoteClusterSvc is an autogenerated mock type for the RemoteClusterSvc type
type RemoteClusterSvc struct {
	mock.Mock
}

// AddRemoteCluster provides a mock function with given fields: ref, skipConnectivityValidation
func (_m *RemoteClusterSvc) AddRemoteCluster(ref *metadata.RemoteClusterReference, skipConnectivityValidation bool) error {
	ret := _m.Called(ref, skipConnectivityValidation)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.RemoteClusterReference, bool) error); ok {
		r0 = rf(ref, skipConnectivityValidation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckAndUnwrapRemoteClusterError provides a mock function with given fields: err
func (_m *RemoteClusterSvc) CheckAndUnwrapRemoteClusterError(err error) (bool, error) {
	ret := _m.Called(err)

	var r0 bool
	if rf, ok := ret.Get(0).(func(error) bool); ok {
		r0 = rf(err)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(error) error); ok {
		r1 = rf(err)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DelRemoteCluster provides a mock function with given fields: refName
func (_m *RemoteClusterSvc) DelRemoteCluster(refName string) (*metadata.RemoteClusterReference, error) {
	ret := _m.Called(refName)

	var r0 *metadata.RemoteClusterReference
	if rf, ok := ret.Get(0).(func(string) *metadata.RemoteClusterReference); ok {
		r0 = rf(refName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.RemoteClusterReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(refName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCapability provides a mock function with given fields: ref
func (_m *RemoteClusterSvc) GetCapability(ref *metadata.RemoteClusterReference) (metadata.Capability, error) {
	ret := _m.Called(ref)

	var r0 metadata.Capability
	if rf, ok := ret.Get(0).(func(*metadata.RemoteClusterReference) metadata.Capability); ok {
		r0 = rf(ref)
	} else {
		r0 = ret.Get(0).(metadata.Capability)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*metadata.RemoteClusterReference) error); ok {
		r1 = rf(ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConnectionStringForRemoteCluster provides a mock function with given fields: ref, isCapiReplication
func (_m *RemoteClusterSvc) GetConnectionStringForRemoteCluster(ref *metadata.RemoteClusterReference, isCapiReplication bool) (string, error) {
	ret := _m.Called(ref, isCapiReplication)

	var r0 string
	if rf, ok := ret.Get(0).(func(*metadata.RemoteClusterReference, bool) string); ok {
		r0 = rf(ref, isCapiReplication)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*metadata.RemoteClusterReference, bool) error); ok {
		r1 = rf(ref, isCapiReplication)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManifestByUuid provides a mock function with given fields: uuid, bucketName, forceRefresh
func (_m *RemoteClusterSvc) GetManifestByUuid(uuid string, bucketName string, forceRefresh bool) (*metadata.CollectionsManifest, error) {
	ret := _m.Called(uuid, bucketName, forceRefresh)

	var r0 *metadata.CollectionsManifest
	if rf, ok := ret.Get(0).(func(string, string, bool) *metadata.CollectionsManifest); ok {
		r0 = rf(uuid, bucketName, forceRefresh)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.CollectionsManifest)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, bool) error); ok {
		r1 = rf(uuid, bucketName, forceRefresh)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRefListForRestartAndClearState provides a mock function with given fields:
func (_m *RemoteClusterSvc) GetRefListForRestartAndClearState() ([]*metadata.RemoteClusterReference, error) {
	ret := _m.Called()

	var r0 []*metadata.RemoteClusterReference
	if rf, ok := ret.Get(0).(func() []*metadata.RemoteClusterReference); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*metadata.RemoteClusterReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRemoteClusterNameFromClusterUuid provides a mock function with given fields: uuid
func (_m *RemoteClusterSvc) GetRemoteClusterNameFromClusterUuid(uuid string) string {
	ret := _m.Called(uuid)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(uuid)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RemoteClusterByRefId provides a mock function with given fields: refId, refresh
func (_m *RemoteClusterSvc) RemoteClusterByRefId(refId string, refresh bool) (*metadata.RemoteClusterReference, error) {
	ret := _m.Called(refId, refresh)

	var r0 *metadata.RemoteClusterReference
	if rf, ok := ret.Get(0).(func(string, bool) *metadata.RemoteClusterReference); ok {
		r0 = rf(refId, refresh)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.RemoteClusterReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(refId, refresh)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoteClusterByRefName provides a mock function with given fields: refName, refresh
func (_m *RemoteClusterSvc) RemoteClusterByRefName(refName string, refresh bool) (*metadata.RemoteClusterReference, error) {
	ret := _m.Called(refName, refresh)

	var r0 *metadata.RemoteClusterReference
	if rf, ok := ret.Get(0).(func(string, bool) *metadata.RemoteClusterReference); ok {
		r0 = rf(refName, refresh)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.RemoteClusterReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(refName, refresh)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoteClusterByUuid provides a mock function with given fields: uuid, refresh
func (_m *RemoteClusterSvc) RemoteClusterByUuid(uuid string, refresh bool) (*metadata.RemoteClusterReference, error) {
	ret := _m.Called(uuid, refresh)

	var r0 *metadata.RemoteClusterReference
	if rf, ok := ret.Get(0).(func(string, bool) *metadata.RemoteClusterReference); ok {
		r0 = rf(uuid, refresh)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.RemoteClusterReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(uuid, refresh)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoteClusterServiceCallback provides a mock function with given fields: path, value, rev
func (_m *RemoteClusterSvc) RemoteClusterServiceCallback(path string, value []byte, rev interface{}) error {
	ret := _m.Called(path, value, rev)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, interface{}) error); ok {
		r0 = rf(path, value, rev)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoteClusters provides a mock function with given fields:
func (_m *RemoteClusterSvc) RemoteClusters() (map[string]*metadata.RemoteClusterReference, error) {
	ret := _m.Called()

	var r0 map[string]*metadata.RemoteClusterReference
	if rf, ok := ret.Get(0).(func() map[string]*metadata.RemoteClusterReference); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*metadata.RemoteClusterReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RequestRemoteMonitoring provides a mock function with given fields: spec
func (_m *RemoteClusterSvc) RequestRemoteMonitoring(spec *metadata.ReplicationSpecification) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetMetadataChangeHandlerCallback provides a mock function with given fields: callBack
func (_m *RemoteClusterSvc) SetMetadataChangeHandlerCallback(callBack base.MetadataChangeHandlerCallback) {
	_m.Called(callBack)
}

// SetRemoteCluster provides a mock function with given fields: refName, ref
func (_m *RemoteClusterSvc) SetRemoteCluster(refName string, ref *metadata.RemoteClusterReference) error {
	ret := _m.Called(refName, ref)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *metadata.RemoteClusterReference) error); ok {
		r0 = rf(refName, ref)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ShouldUseAlternateAddress provides a mock function with given fields: ref
func (_m *RemoteClusterSvc) ShouldUseAlternateAddress(ref *metadata.RemoteClusterReference) (bool, error) {
	ret := _m.Called(ref)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*metadata.RemoteClusterReference) bool); ok {
		r0 = rf(ref)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*metadata.RemoteClusterReference) error); ok {
		r1 = rf(ref)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnRequestRemoteMonitoring provides a mock function with given fields: spec
func (_m *RemoteClusterSvc) UnRequestRemoteMonitoring(spec *metadata.ReplicationSpecification) error {
	ret := _m.Called(spec)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.ReplicationSpecification) error); ok {
		r0 = rf(spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateAddRemoteCluster provides a mock function with given fields: ref
func (_m *RemoteClusterSvc) ValidateAddRemoteCluster(ref *metadata.RemoteClusterReference) error {
	ret := _m.Called(ref)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.RemoteClusterReference) error); ok {
		r0 = rf(ref)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateRemoteCluster provides a mock function with given fields: ref
func (_m *RemoteClusterSvc) ValidateRemoteCluster(ref *metadata.RemoteClusterReference) error {
	ret := _m.Called(ref)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.RemoteClusterReference) error); ok {
		r0 = rf(ref)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateSetRemoteCluster provides a mock function with given fields: refName, ref
func (_m *RemoteClusterSvc) ValidateSetRemoteCluster(refName string, ref *metadata.RemoteClusterReference) error {
	ret := _m.Called(refName, ref)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *metadata.RemoteClusterReference) error); ok {
		r0 = rf(refName, ref)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
