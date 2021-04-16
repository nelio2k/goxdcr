// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	log "github.com/couchbase/goxdcr/log"

	mock "github.com/stretchr/testify/mock"
)

// ConflictResolver is an autogenerated mock type for the ConflictResolver type
type ConflictResolver struct {
	mock.Mock
}

// Execute provides a mock function with given fields: doc_metadata_source, doc_metadata_target, source_cr_mode, xattrEnabled, logger
func (_m *ConflictResolver) Execute(doc_metadata_source parts.documentMetadata, doc_metadata_target parts.documentMetadata, source_cr_mode base.ConflictResolutionMode, xattrEnabled bool, logger *log.CommonLogger) bool {
	ret := _m.Called(doc_metadata_source, doc_metadata_target, source_cr_mode, xattrEnabled, logger)

	var r0 bool
	if rf, ok := ret.Get(0).(func(parts.documentMetadata, parts.documentMetadata, base.ConflictResolutionMode, bool, *log.CommonLogger) bool); ok {
		r0 = rf(doc_metadata_source, doc_metadata_target, source_cr_mode, xattrEnabled, logger)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
