// Copyright (c) 2013-2019 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package metadata

import (
	"fmt"
	"github.com/couchbase/goxdcr/base"
	"reflect"
	"strings"
)

type ReplicationSpecApi interface {
	Id() string
	InternalId() string
	SourceBucketName() string
	SourceBucketUUID() string
	TargetClusterUUID() string
	TargetBucketName() string
	TargetBucketUUID() string
	Settings() *ReplicationSettings
}

/************************************
/* struct ReplicationSpecification
*************************************/
type ReplicationSpecification struct {
	//id of the replication
	Id_ string `json:"id"`

	// internal id, used to detect the case when replication spec has been deleted and recreated
	InternalId_ string `json:"internalId"`

	// Source Bucket Name
	SourceBucketName_ string `json:"sourceBucketName"`

	//Source Bucket UUID
	SourceBucketUUID_ string `json:"sourceBucketUUID"`

	//Target Cluster UUID
	TargetClusterUUID_ string `json:"targetClusterUUID"`

	// Target Bucket Name
	TargetBucketName_ string `json:"targetBucketName"`

	TargetBucketUUID_ string `json:"targetBucketUUID"`

	Settings_ *ReplicationSettings `json:"replicationSettings"`

	// revision number to be used by metadata service. not included in json
	Revision interface{}
}

func NewReplicationSpecification(sourceBucketName string, sourceBucketUUID string, targetClusterUUID string, targetBucketName string, targetBucketUUID string) (*ReplicationSpecification, error) {
	randId, err := base.GenerateRandomId(base.LengthOfRandomId, base.MaxRetryForRandomIdGeneration)
	if err != nil {
		return nil, err
	}
	return &ReplicationSpecification{Id_: ReplicationId(sourceBucketName, targetClusterUUID, targetBucketName),
		InternalId_:        randId,
		SourceBucketName_:  sourceBucketName,
		SourceBucketUUID_:  sourceBucketUUID,
		TargetClusterUUID_: targetClusterUUID,
		TargetBucketName_:  targetBucketName,
		TargetBucketUUID_:  targetBucketUUID,
		Settings_:          DefaultReplicationSettings()}, nil
}

func (spec *ReplicationSpecification) String() string {
	if spec == nil {
		return ""
	}
	var specSettingsMap ReplicationSettingsMap
	if spec.Settings() != nil {
		specSettingsMap = spec.Settings().CloneAndRedact().ToMap(false /*defaultSettings_*/)
	}
	return fmt.Sprintf("Id_: %v InternalId_: %v SourceBucketName_: %v SourceBucketUUID_: %v TargetClusterUUID_: %v TargetBucketName_: %v TargetBucketUUID_: %v Settings_: %v",
		spec.Id(), spec.InternalId(), spec.SourceBucketName(), spec.SourceBucketUUID(), spec.TargetClusterUUID(), spec.TargetBucketName(), spec.TargetBucketUUID(), specSettingsMap)
}

func (s *ReplicationSpecification) Id() string {
	return s.Id_
}
func (s *ReplicationSpecification) InternalId() string {
	return s.InternalId_
}
func (s *ReplicationSpecification) SourceBucketName() string {
	return s.SourceBucketName_
}
func (s *ReplicationSpecification) SourceBucketUUID() string {
	return s.SourceBucketUUID_
}
func (s *ReplicationSpecification) TargetClusterUUID() string {
	return s.TargetClusterUUID_
}
func (s *ReplicationSpecification) TargetBucketName() string {
	return s.TargetBucketName_
}
func (s *ReplicationSpecification) TargetBucketUUID() string {
	return s.TargetBucketUUID_
}
func (s *ReplicationSpecification) Settings() *ReplicationSettings {
	return s.Settings_
}

// checks if the passed in spec is the same as the current spec
// used to check if a spec in cache needs to be refreshed
func (spec *ReplicationSpecification) SameSpec(spec2 *ReplicationSpecification) bool {
	if spec == nil {
		return spec2 == nil
	}
	if spec2 == nil {
		return false
	}
	// note that settings in spec are not compared. The assumption is that if settings are different, Revision will have to be different
	return spec.Id() == spec2.Id() && spec.InternalId() == spec2.InternalId() &&
		spec.SourceBucketName() == spec2.SourceBucketName() &&
		spec.SourceBucketUUID() == spec2.SourceBucketUUID() &&
		spec.TargetClusterUUID() == spec2.TargetClusterUUID() && spec.TargetBucketName() == spec2.TargetBucketName() &&
		spec.TargetBucketUUID() == spec2.TargetBucketUUID() && reflect.DeepEqual(spec.Revision, spec2.Revision)
}

func (spec *ReplicationSpecification) Clone() *ReplicationSpecification {
	if spec == nil {
		return nil
	}
	return &ReplicationSpecification{Id_: spec.Id(),
		InternalId_:        spec.InternalId(),
		SourceBucketName_:  spec.SourceBucketName(),
		SourceBucketUUID_:  spec.SourceBucketUUID(),
		TargetClusterUUID_: spec.TargetClusterUUID(),
		TargetBucketName_:  spec.TargetBucketName(),
		TargetBucketUUID_:  spec.TargetBucketUUID(),
		Settings_:          spec.Settings().Clone(),
		// !!! shallow copy of revision.
		// spec.Revision should only be passed along and should never be modified
		Revision: spec.Revision}
}

func (spec *ReplicationSpecification) Redact() *ReplicationSpecification {
	if spec != nil {
		// Currently only the Settings_ has user identifiable data in filtered expression
		spec.Settings().Redact()
	}
	return spec
}

func (spec *ReplicationSpecification) CloneAndRedact() *ReplicationSpecification {
	if spec != nil {
		return spec.Clone().Redact()
	}
	return spec
}

func (spec *ReplicationSpecification) SameSpecGeneric(other GenericSpecification) bool {
	return spec.SameSpec(other.(*ReplicationSpecification))
}

func (spec *ReplicationSpecification) CloneGeneric() GenericSpecification {
	return spec.Clone()
}

func (spec *ReplicationSpecification) RedactGeneric() GenericSpecification {
	return spec.Redact()
}

func ReplicationId(sourceBucketName string, targetClusterUUID string, targetBucketName string) string {
	parts := []string{targetClusterUUID, sourceBucketName, targetBucketName}
	return strings.Join(parts, base.KeyPartsDelimiter)
}

func IsReplicationIdForSourceBucket(replicationId_ string, sourceBucketName string) (bool, error) {
	replBucketName, err := GetSourceBucketNameFromReplicationId(replicationId_)
	if err != nil {
		return false, err
	} else {
		return replBucketName == sourceBucketName, nil
	}
}

func IsReplicationIdForTargetBucket(replicationId_ string, targetBucketName string) (bool, error) {
	replBucketName, err := GetTargetBucketNameFromReplicationId(replicationId_)
	if err != nil {
		return false, err
	} else {
		return replBucketName == targetBucketName, nil
	}
}

func GetSourceBucketNameFromReplicationId(replicationId_ string) (string, error) {
	parts := strings.Split(replicationId_, base.KeyPartsDelimiter)
	if len(parts) == 3 {
		return parts[1], nil
	} else {
		return "", fmt.Errorf("Invalid replication id: %v", replicationId_)
	}
}

func GetTargetBucketNameFromReplicationId(replicationId_ string) (string, error) {
	parts := strings.Split(replicationId_, base.KeyPartsDelimiter)
	if len(parts) == 3 {
		return parts[2], nil
	} else {
		return "", fmt.Errorf("Invalid replication id: %v", replicationId_)
	}
}
