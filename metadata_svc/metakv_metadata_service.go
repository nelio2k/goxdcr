// Copyright (c) 2013-2017 Couchbase, Inc.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
// except in compliance with the License. You may obtain a copy of the License at
//   http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing permissions
// and limitations under the License.

// metadata service implementation leveraging metakv
package metadata_svc

import (
	"fmt"
	"github.com/couchbase/cbauth/metakv"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/service_def"
	utilities "github.com/couchbase/goxdcr/utils"
	"strings"
	"sync"
	"time"
)

type MetaKVMetadataSvc struct {
	logger   *log.CommonLogger
	utils    utilities.UtilsIface
	readOnly bool
}

func NewMetaKVMetadataSvc(logger_ctx *log.LoggerContext, utilsIn utilities.UtilsIface, readOnly bool) (*MetaKVMetadataSvc, error) {
	return &MetaKVMetadataSvc{
		logger:   log.NewLogger("MetadataSvc", logger_ctx),
		utils:    utilsIn,
		readOnly: readOnly,
	}, nil
}

//Wrap metakv.Get with retries
//if the key is not found in metakv, return nil, nil, service_def.MetadataNotFoundErr
//if metakv operation failed after max number of retries, return nil, nil, ErrorFailedAfterRetry
func (meta_svc *MetaKVMetadataSvc) Get(key string) ([]byte, interface{}, error) {
	var value []byte
	var rev interface{}
	var err error
	start_time := time.Now()
	defer meta_svc.logger.Debugf("Took %vs to get %v to metakv\n", time.Since(start_time).Seconds(), key)

	metakvOpGetFunc := func() error {
		value, rev, err = metakv.Get(getPathFromKey(key))
		if value == nil && rev == nil && err == nil {
			meta_svc.logger.Debugf("Can't find key=%v", key)
			err = service_def.MetadataNotFoundErr
			return nil
		} else if err == nil {
			return nil
		} else {
			meta_svc.logger.Warnf("metakv.Get failed. path=%v, err=%v", getPathFromKey(key), err)
			return err
		}
	}

	expOpErr := meta_svc.utils.ExponentialBackoffExecutor("metakv_svc_getOp", base.RetryIntervalMetakv, base.MaxNumOfMetakvRetries,
		base.MetaKvBackoffFactor, metakvOpGetFunc)

	if expOpErr != nil {
		// Executor will return error only if it timed out with ErrorFailedAfterRetry. Log it and override the ret err
		meta_svc.logger.Errorf("metakv.Get failed after max retry. path=%v, err=%v", getPathFromKey(key), err)
		err = expOpErr
	}

	return value, rev, err
}

func (meta_svc *MetaKVMetadataSvc) Add(key string, value []byte) error {
	if meta_svc.readOnly {
		return nil
	}
	return meta_svc.add(key, value, false)
}

func (meta_svc *MetaKVMetadataSvc) AddSensitive(key string, value []byte) error {
	if meta_svc.readOnly {
		return nil
	}
	return meta_svc.add(key, value, true)
}

//Wrap metakv.Add with retries
//if the key is already exist in metakv, return service_def.ErrorKeyAlreadyExist
//if metakv operation failed after max number of retries, return ErrorFailedAfterRetry
func (meta_svc *MetaKVMetadataSvc) add(key string, value []byte, sensitive bool) error {
	var err error
	var redactedValue []byte
	valueToPrint := &value
	start_time := time.Now()
	defer meta_svc.logger.Debugf("Took %vs to add %v to metakv\n", time.Since(start_time).Seconds(), key)

	var redactOnceSync sync.Once
	redactOnce := func() {
		redactOnceSync.Do(func() {
			if sensitive {
				redactedValue = base.DeepCopyByteArray(value)
				redactedValue = base.TagUDBytes(redactedValue)
				valueToPrint = &redactedValue
			}
		})
	}

	metakvOpAddFunc := func() error {
		if sensitive {
			err = metakv.AddSensitive(getPathFromKey(key), value)
		} else {
			err = metakv.Add(getPathFromKey(key), value)
		}
		if err == metakv.ErrRevMismatch {
			err = service_def.ErrorKeyAlreadyExist
			return nil
		} else if err == nil {
			return nil
		} else {
			redactOnce()
			meta_svc.logger.Warnf("metakv.Add failed. key=%v, value=%v, err=%v\n", key, valueToPrint, err)
			return err
		}
	}

	expOpErr := meta_svc.utils.ExponentialBackoffExecutor("metakv_svc_addOp", base.RetryIntervalMetakv, base.MaxNumOfMetakvRetries,
		base.MetaKvBackoffFactor, metakvOpAddFunc)

	if expOpErr != nil {
		// Executor will return error only if it timed out with ErrorFailedAfterRetry. Log it and override the ret err
		redactOnce()
		meta_svc.logger.Errorf("metakv.Add failed after max retry. key=%v, value=%v, err=%v\n", key, valueToPrint, err)
		err = expOpErr
	}

	return err
}

func (meta_svc *MetaKVMetadataSvc) Set(key string, value []byte, rev interface{}) error {
	if meta_svc.readOnly {
		return nil
	}
	return meta_svc.set(key, value, rev, false)
}

func (meta_svc *MetaKVMetadataSvc) SetSensitive(key string, value []byte, rev interface{}) error {
	if meta_svc.readOnly {
		return nil
	}
	return meta_svc.set(key, value, rev, true)
}

//Wrap metakv.Set with retries
//if the rev provided doesn't match with the rev metakv has, return service_def.ErrorRevisionMismatch
//if metakv operation failed after max number of retries, return ErrorFailedAfterRetry
func (meta_svc *MetaKVMetadataSvc) set(key string, value []byte, rev interface{}, sensitive bool) error {
	var err error
	var redactedValue []byte
	valueToPrint := &value
	start_time := time.Now()
	defer meta_svc.logger.Debugf("Took %vs to set %v to metakv\n", time.Since(start_time).Seconds(), key)

	var redactOnceSync sync.Once
	redactOnce := func() {
		redactOnceSync.Do(func() {
			if sensitive {
				redactedValue = base.DeepCopyByteArray(value)
				redactedValue = base.TagUDBytes(redactedValue)
				valueToPrint = &redactedValue
			}
		})
	}

	metakvOpSetFunc := func() error {
		if sensitive {
			err = metakv.SetSensitive(getPathFromKey(key), value, rev)
		} else {
			err = metakv.Set(getPathFromKey(key), value, rev)
		}
		if err == metakv.ErrRevMismatch {
			err = service_def.ErrorRevisionMismatch
			return nil
		} else if err == nil {
			return nil
		} else {
			redactOnce()
			meta_svc.logger.Warnf("metakv.Set failed. key=%v, value=%v, err=%v\n", key, valueToPrint, err)
			return err
		}
	}

	expOpErr := meta_svc.utils.ExponentialBackoffExecutor("metakv_svc_setOp", base.RetryIntervalMetakv, base.MaxNumOfMetakvRetries,
		base.MetaKvBackoffFactor, metakvOpSetFunc)

	if expOpErr != nil {
		redactOnce()
		meta_svc.logger.Errorf("metakv.Set failed after max retry. key=%v, value=%v, err=%v\n", key, valueToPrint, err)
		err = expOpErr
	}
	return err
}

//Wrap metakv.Del with retries
//if the rev provided doesn't match with the rev metakv has, return service_def.ErrorRevisionMismatch
//if metakv operation failed after max number of retries, return ErrorFailedAfterRetry
func (meta_svc *MetaKVMetadataSvc) Del(key string, rev interface{}) error {
	if meta_svc.readOnly {
		return nil
	}
	var err error
	start_time := time.Now()
	defer meta_svc.logger.Debugf("Took %vs to delete %v from metakv\n", time.Since(start_time).Seconds(), key)

	metakvOpDelFunc := func() error {
		err = metakv.Delete(getPathFromKey(key), rev)
		if err == metakv.ErrRevMismatch {
			err = service_def.ErrorRevisionMismatch
			return nil
		} else if err == nil {
			return nil
		} else {
			meta_svc.logger.Warnf("metakv.Delete failed. key=%v, rev=%v, err=%v\n", key, rev, err)
			return err
		}
	}

	expOpErr := meta_svc.utils.ExponentialBackoffExecutor("metakv_svc_delOp", base.RetryIntervalMetakv, base.MaxNumOfMetakvRetries,
		base.MetaKvBackoffFactor, metakvOpDelFunc)

	if expOpErr != nil {
		// Executor will return error only if it timed out with ErrorFailedAfterRetry. Log it and override the ret err
		meta_svc.logger.Errorf("metakv.Delete failed. key=%v, rev=%v, err=%v\n", key, rev, err)
		err = expOpErr
	}
	return err
}

//Wrap metakv.RecursiveDelete with retries
//if metakv operation failed after max number of retries, return ErrorFailedAfterRetry
func (meta_svc *MetaKVMetadataSvc) DelAllFromCatalog(catalogKey string) error {
	if meta_svc.readOnly {
		return nil
	}
	start_time := time.Now()
	defer meta_svc.logger.Debugf("Took %vs to RecursiveDelete for catalogKey=%v to metakv\n", time.Since(start_time).Seconds(), catalogKey)

	metakvOpDelRFunc := func() error {
		err := metakv.RecursiveDelete(GetCatalogPathFromCatalogKey(catalogKey))
		if err == nil {
			return err
		} else {
			meta_svc.logger.Warnf("metakv.RecursiveDelete failed. catalogKey=%v, err=%v\n", catalogKey, err)
			return err
		}
	}

	expOpErr := meta_svc.utils.ExponentialBackoffExecutor("metakv_svc_delROp", base.RetryIntervalMetakv, base.MaxNumOfMetakvRetries,
		base.MetaKvBackoffFactor, metakvOpDelRFunc)

	if expOpErr != nil {
		// Executor will return error only if it timed out with ErrorFailedAfterRetry. Log it and override the ret err
		meta_svc.logger.Errorf("metakv.RecursiveDelete failed after max retry. catalogKey=%v\n", catalogKey)
	}
	return expOpErr
}

//Wrap metakv.ListAllChildren with retries
//if metakv operation failed after max number of retries, return ErrorFailedAfterRetry
func (meta_svc *MetaKVMetadataSvc) GetAllMetadataFromCatalog(catalogKey string) ([]*service_def.MetadataEntry, error) {
	var entries = make([]*service_def.MetadataEntry, 0)
	start_time := time.Now()
	defer meta_svc.logger.Debugf("Took %vs to ListAllChildren for catalogKey=%v to metakv\n", time.Since(start_time).Seconds(), catalogKey)

	metakvOpGetAllMetadataFunc := func() error {
		kvEntries, err := metakv.ListAllChildren(GetCatalogPathFromCatalogKey(catalogKey))
		if err == nil {
			for _, kvEntry := range kvEntries {
				entries = append(entries, &service_def.MetadataEntry{GetKeyFromPath(kvEntry.Path), kvEntry.Value, kvEntry.Rev})
			}
			return err
		} else {
			meta_svc.logger.Warnf("metakv.ListAllChildren failed. path=%v, err=%v\n", GetCatalogPathFromCatalogKey(catalogKey), err)
			return err
		}
	}

	expOpErr := meta_svc.utils.ExponentialBackoffExecutor("metakv_svc_getAllMetaOp", base.RetryIntervalMetakv, base.MaxNumOfMetakvRetries,
		base.MetaKvBackoffFactor, metakvOpGetAllMetadataFunc)

	if expOpErr != nil {
		// Executor will return error only if it timed out with ErrorFailedAfterRetry. Log it and override the ret err
		meta_svc.logger.Errorf("metakv.ListAllChildren failed after max retry. path=%v\n", GetCatalogPathFromCatalogKey(catalogKey))
	}
	return entries, expOpErr
}

// get all keys from a catalog
func (meta_svc *MetaKVMetadataSvc) GetAllKeysFromCatalog(catalogKey string) ([]string, error) {
	keys := make([]string, 0)

	metaEntries, err := meta_svc.GetAllMetadataFromCatalog(catalogKey)
	if err != nil {
		return nil, err
	}
	for _, metaEntry := range metaEntries {
		keys = append(keys, metaEntry.Key)
	}
	return keys, nil
}

// metakv requires that all paths start with "/"
func getPathFromKey(key string) string {
	return base.KeyPartsDelimiter + key
}

// the following are exposed since they are needed by metakv call back function

// metakv requires that all parent paths start and end with "/"
func GetCatalogPathFromCatalogKey(catalogKey string) string {
	return base.KeyPartsDelimiter + catalogKey + base.KeyPartsDelimiter
}

func GetKeyFromPath(path string) string {
	if strings.HasPrefix(path, base.KeyPartsDelimiter) {
		return path[len(base.KeyPartsDelimiter):]
	} else {
		// should never get here
		panic(fmt.Sprintf("path=%v doesn't start with '/'", path))
	}
}
