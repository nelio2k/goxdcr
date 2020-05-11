// Copyright (c) 2013 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package common

import (
	"github.com/couchbase/goxdcr/metadata"
)

//PipelineService can be any component that monitors, does logging, keeps state for the pipeline
//Each PipelineService is a goroutine that run parallelly
type PipelineService interface {
	Attach(pipeline Pipeline) error

	Start(metadata.ReplicationSettingsMap) error
	Stop() error

	UpdateSettings(settings metadata.ReplicationSettingsMap) error

	// If a service can be shared by more than one pipeline in the same replication
	IsSharable() bool
	// If sharable, allow service to be detached
	Detach(pipeline Pipeline) error
}
