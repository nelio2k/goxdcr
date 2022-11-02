// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	metadata "github.com/couchbase/goxdcr/metadata"

	mock "github.com/stretchr/testify/mock"

	service_def "github.com/couchbase/goxdcr/service_def"

	sync "sync"
)

// CheckpointsService is an autogenerated mock type for the CheckpointsService type
type CheckpointsService struct {
	mock.Mock
}

// BackfillReplicationSpecChangeCallback provides a mock function with given fields: metadataId, oldMetadata, newMetadata
func (_m *CheckpointsService) BackfillReplicationSpecChangeCallback(metadataId string, oldMetadata interface{}, newMetadata interface{}) error {
	ret := _m.Called(metadataId, oldMetadata, newMetadata)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) error); ok {
		r0 = rf(metadataId, oldMetadata, newMetadata)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckpointsDoc provides a mock function with given fields: replicationId, vbno
func (_m *CheckpointsService) CheckpointsDoc(replicationId string, vbno uint16) (*metadata.CheckpointsDoc, error) {
	ret := _m.Called(replicationId, vbno)

	var r0 *metadata.CheckpointsDoc
	if rf, ok := ret.Get(0).(func(string, uint16) *metadata.CheckpointsDoc); ok {
		r0 = rf(replicationId, vbno)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.CheckpointsDoc)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, uint16) error); ok {
		r1 = rf(replicationId, vbno)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CheckpointsDocs provides a mock function with given fields: replicationId, brokenMappingsNeeded
func (_m *CheckpointsService) CheckpointsDocs(replicationId string, brokenMappingsNeeded bool) (map[uint16]*metadata.CheckpointsDoc, error) {
	ret := _m.Called(replicationId, brokenMappingsNeeded)

	var r0 map[uint16]*metadata.CheckpointsDoc
	if rf, ok := ret.Get(0).(func(string, bool) map[uint16]*metadata.CheckpointsDoc); ok {
		r0 = rf(replicationId, brokenMappingsNeeded)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[uint16]*metadata.CheckpointsDoc)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, bool) error); ok {
		r1 = rf(replicationId, brokenMappingsNeeded)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionsManifestChangeCb provides a mock function with given fields: metadataId, oldMetadata, newMetadata
func (_m *CheckpointsService) CollectionsManifestChangeCb(metadataId string, oldMetadata interface{}, newMetadata interface{}) error {
	ret := _m.Called(metadataId, oldMetadata, newMetadata)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}) error); ok {
		r0 = rf(metadataId, oldMetadata, newMetadata)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DelCheckpointsDoc provides a mock function with given fields: replicationId, vbno, specInternalId
func (_m *CheckpointsService) DelCheckpointsDoc(replicationId string, vbno uint16, specInternalId string) error {
	ret := _m.Called(replicationId, vbno, specInternalId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, uint16, string) error); ok {
		r0 = rf(replicationId, vbno, specInternalId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DelCheckpointsDocs provides a mock function with given fields: replicationId
func (_m *CheckpointsService) DelCheckpointsDocs(replicationId string) error {
	ret := _m.Called(replicationId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(replicationId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DisableRefCntDecrement provides a mock function with given fields: topic
func (_m *CheckpointsService) DisableRefCntDecrement(topic string) {
	_m.Called(topic)
}

// EnableRefCntDecrement provides a mock function with given fields: topic
func (_m *CheckpointsService) EnableRefCntDecrement(topic string) {
	_m.Called(topic)
}

// GetCkptsMappingsCleanupCallback provides a mock function with given fields: specId, specInternalId, toBeRemoved
func (_m *CheckpointsService) GetCkptsMappingsCleanupCallback(specId string, specInternalId string, toBeRemoved metadata.ScopesMap) (base.StoppedPipelineCallback, base.StoppedPipelineErrCallback) {
	ret := _m.Called(specId, specInternalId, toBeRemoved)

	var r0 base.StoppedPipelineCallback
	if rf, ok := ret.Get(0).(func(string, string, metadata.ScopesMap) base.StoppedPipelineCallback); ok {
		r0 = rf(specId, specInternalId, toBeRemoved)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(base.StoppedPipelineCallback)
		}
	}

	var r1 base.StoppedPipelineErrCallback
	if rf, ok := ret.Get(1).(func(string, string, metadata.ScopesMap) base.StoppedPipelineErrCallback); ok {
		r1 = rf(specId, specInternalId, toBeRemoved)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(base.StoppedPipelineErrCallback)
		}
	}

	return r0, r1
}

// GetVbnosFromCheckpointDocs provides a mock function with given fields: replicationId
func (_m *CheckpointsService) GetVbnosFromCheckpointDocs(replicationId string) ([]uint16, error) {
	ret := _m.Called(replicationId)

	var r0 []uint16
	if rf, ok := ret.Get(0).(func(string) []uint16); ok {
		r0 = rf(replicationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]uint16)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(replicationId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadBrokenMappings provides a mock function with given fields: replicationId
func (_m *CheckpointsService) LoadBrokenMappings(replicationId string) (metadata.ShaToCollectionNamespaceMap, *metadata.CollectionNsMappingsDoc, service_def.IncrementerFunc, bool, error) {
	ret := _m.Called(replicationId)

	var r0 metadata.ShaToCollectionNamespaceMap
	if rf, ok := ret.Get(0).(func(string) metadata.ShaToCollectionNamespaceMap); ok {
		r0 = rf(replicationId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(metadata.ShaToCollectionNamespaceMap)
		}
	}

	var r1 *metadata.CollectionNsMappingsDoc
	if rf, ok := ret.Get(1).(func(string) *metadata.CollectionNsMappingsDoc); ok {
		r1 = rf(replicationId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*metadata.CollectionNsMappingsDoc)
		}
	}

	var r2 service_def.IncrementerFunc
	if rf, ok := ret.Get(2).(func(string) service_def.IncrementerFunc); ok {
		r2 = rf(replicationId)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(service_def.IncrementerFunc)
		}
	}

	var r3 bool
	if rf, ok := ret.Get(3).(func(string) bool); ok {
		r3 = rf(replicationId)
	} else {
		r3 = ret.Get(3).(bool)
	}

	var r4 error
	if rf, ok := ret.Get(4).(func(string) error); ok {
		r4 = rf(replicationId)
	} else {
		r4 = ret.Error(4)
	}

	return r0, r1, r2, r3, r4
}

// PreUpsertBrokenMapping provides a mock function with given fields: replicationId, specInternalId, oneBrokenMapping
func (_m *CheckpointsService) PreUpsertBrokenMapping(replicationId string, specInternalId string, oneBrokenMapping *metadata.CollectionNamespaceMapping) error {
	ret := _m.Called(replicationId, specInternalId, oneBrokenMapping)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, *metadata.CollectionNamespaceMapping) error); ok {
		r0 = rf(replicationId, specInternalId, oneBrokenMapping)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReplicationSpecChangeCallback provides a mock function with given fields: metadataId, oldMetadata, newMetadata, wg
func (_m *CheckpointsService) ReplicationSpecChangeCallback(metadataId string, oldMetadata interface{}, newMetadata interface{}, wg *sync.WaitGroup) error {
	ret := _m.Called(metadataId, oldMetadata, newMetadata, wg)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}, *sync.WaitGroup) error); ok {
		r0 = rf(metadataId, oldMetadata, newMetadata, wg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertAndReloadCheckpointCompleteSet provides a mock function with given fields: replicationId, mappingDoc, ckptDoc, internalId
func (_m *CheckpointsService) UpsertAndReloadCheckpointCompleteSet(replicationId string, mappingDoc *metadata.CollectionNsMappingsDoc, ckptDoc map[uint16]*metadata.CheckpointsDoc, internalId string) error {
	ret := _m.Called(replicationId, mappingDoc, ckptDoc, internalId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *metadata.CollectionNsMappingsDoc, map[uint16]*metadata.CheckpointsDoc, string) error); ok {
		r0 = rf(replicationId, mappingDoc, ckptDoc, internalId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertBrokenMapping provides a mock function with given fields: replicationId, specInternalId
func (_m *CheckpointsService) UpsertBrokenMapping(replicationId string, specInternalId string) error {
	ret := _m.Called(replicationId, specInternalId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(replicationId, specInternalId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertCheckpoints provides a mock function with given fields: replicationId, specInternalId, vbno, ckpt_record
func (_m *CheckpointsService) UpsertCheckpoints(replicationId string, specInternalId string, vbno uint16, ckpt_record *metadata.CheckpointRecord) (int, error) {
	ret := _m.Called(replicationId, specInternalId, vbno, ckpt_record)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, string, uint16, *metadata.CheckpointRecord) int); ok {
		r0 = rf(replicationId, specInternalId, vbno, ckpt_record)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, uint16, *metadata.CheckpointRecord) error); ok {
		r1 = rf(replicationId, specInternalId, vbno, ckpt_record)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpsertCheckpointsDone provides a mock function with given fields: replicationId
func (_m *CheckpointsService) UpsertCheckpointsDone(replicationId string) error {
	ret := _m.Called(replicationId)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(replicationId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCheckpointsService interface {
	mock.TestingT
	Cleanup(func())
}

// NewCheckpointsService creates a new instance of CheckpointsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCheckpointsService(t mockConstructorTestingTNewCheckpointsService) *CheckpointsService {
	mock := &CheckpointsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
