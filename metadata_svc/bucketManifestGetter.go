// Copyright (c) 2013-2020 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

package metadata_svc

import (
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/service_def"
	"sync/atomic"
	"time"
)

/**
 * The bucket manifest getter's purpose is to allow other service to specify a getter function
 * and then to do burst-control to prevent too many calls within a time period to overload
 * the manifest provider (ns_server). It will still release the latest pulled manifest
 * to every single caller. Because manifests are eventual consistent (i.e. KV can send down
 * events before ns_server has it for pulling), a little lag shouldn't hurt anyone
 */
type BucketManifestGetter struct {
	bucketName          string
	getterFunc          func(string) (*metadata.CollectionsManifest, error)
	lastQueryTime       time.Time
	lastStoredManifest  *metadata.CollectionsManifest
	firstCallerCh       chan bool
	checkInterval       time.Duration
	getterFuncOperating uint32
}

func NewBucketManifestGetter(bucketName string, manifestOps service_def.CollectionsManifestOps, checkInterval time.Duration) *BucketManifestGetter {
	getter := &BucketManifestGetter{
		bucketName:    bucketName,
		getterFunc:    manifestOps.CollectionManifestGetter,
		firstCallerCh: make(chan bool, 1),
		checkInterval: checkInterval,
	}
	// Allow one caller to call it
	getter.firstCallerCh <- true
	return getter
}

func (s *BucketManifestGetter) GetManifest() *metadata.CollectionsManifest {
	switch {
	case <-s.firstCallerCh:
		atomic.StoreUint32(&s.getterFuncOperating, 1)
		if time.Now().Sub(s.lastQueryTime) > (s.checkInterval) {
			// Prevent overwhelming the ns_server, only query every "checkInterval" seconds
			manifest, err := s.getterFunc(s.bucketName)
			if err == nil {
				s.lastStoredManifest = manifest
				s.lastQueryTime = time.Now()
			}
		}
		defer func() {
			atomic.StoreUint32(&s.getterFuncOperating, 0)
			s.firstCallerCh <- true
		}()
		return s.lastStoredManifest
	default:
		// Concurrent calls are really not common
		// so this rudimentary synchronization method should be ok
		// This is so that if a second caller calls after
		// the first caller, the second caller will be able to receive
		// the most up-to-date lastStoredManifest
		for atomic.LoadUint32(&s.getterFuncOperating) == 1 {
			time.Sleep(300 * time.Nanosecond)
		}
		return s.lastStoredManifest
	}
}
