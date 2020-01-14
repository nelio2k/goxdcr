// Copyright (c) 2013-2020 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package service_def

import (
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/metadata"
)

type BackfillReplSvc interface {
	ReplicationSpec(replicationId string) (*metadata.BackfillReplicationSpec, error)
	AddReplicationSpec(spec *metadata.BackfillReplicationSpec) error
	SetReplicationSpec(spec *metadata.BackfillReplicationSpec) error
	DelReplicationSpec(replicationId string) (*metadata.BackfillReplicationSpec, error)
	//	AllReplicationSpecs() (map[string]*metadata.BackfillReplicationSpec, error)

	SetMetadataChangeHandlerCallback(callBack base.MetadataChangeHandlerCallback)

	//get the derived object (i.e. ReplicationStatus) for the specification
	//this is used to keep the derived object and replication spec in the same cache
	GetDerivedObj(specId string) (interface{}, error)

	//set the derived object (i.e ReplicationStatus) for the specification
	SetDerivedObj(specId string, derivedObj interface{}) error
}
