// Copyright 2013-Present Couchbase, Inc.
//
// Use of this software is governed by the Business Source License included in
// the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
// file, in accordance with the Business Source License, use of this software
// will be governed by the Apache License, Version 2.0, included in the file
// licenses/APL2.txt.

package metadata_svc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/common"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/service_def"
	utilities "github.com/couchbase/goxdcr/utils"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const (
	// the key to the metadata that stores the keys of all remote clusters
	CheckpointsCatalogKeyPrefix = "ckpt"
	CheckpointsKeyPrefix        = CheckpointsCatalogKeyPrefix
	BrokenMappingKey            = "brokenMappings"
)

// Used to keep track of brokenmapping SHA and the count of checkpoint records referring to it
type BrokenMapShaRefCounter struct {
	lock           sync.RWMutex
	refCnt         map[string]uint64
	shaToMapping   metadata.ShaToCollectionNamespaceMap
	needToSync     bool // needs to sync refCnt to shaMap and then also persist to metakv
	internalSpecId string
}

func NewBrokenMapShaRefCounter() *BrokenMapShaRefCounter {
	return &BrokenMapShaRefCounter{refCnt: make(map[string]uint64),
		shaToMapping: make(metadata.ShaToCollectionNamespaceMap),
	}
}

func (b *BrokenMapShaRefCounter) RegisterCkptDoc(doc *metadata.CheckpointsDoc) {
	if b == nil || doc == nil {
		return
	}

	b.lock.Lock()
	for _, record := range doc.Checkpoint_records {
		if record == nil || len(record.BrokenMappingSha256) == 0 {
			continue
		}
		b.refCnt[record.BrokenMappingSha256]++
		b.shaToMapping[record.BrokenMappingSha256] = record.BrokenMappings()
	}
	b.lock.Unlock()
}

type CheckpointsService struct {
	*ShaRefCounterService
	metadata_svc service_def.MetadataSvc
	logger       *log.CommonLogger
	utils        utilities.UtilsIface

	specsMtx    sync.RWMutex
	cachedSpecs map[string]*metadata.ReplicationSpecification
	// Certain situations such as merging checkpoints from other nodes require coarse locking
	stopTheWorldMtx map[string]*sync.RWMutex

	backfillSpecsMtx    sync.RWMutex
	cachedBackfillSpecs map[string]*metadata.BackfillReplicationSpec

	replicationSpecSvc service_def.ReplicationSpecSvc
}

func NewCheckpointsService(metadata_svc service_def.MetadataSvc, logger_ctx *log.LoggerContext, utils utilities.UtilsIface, replicationSpecService service_def.ReplicationSpecSvc) (*CheckpointsService, error) {
	logger := log.NewLogger("CheckpointSvc", logger_ctx)
	ckptSvc := &CheckpointsService{metadata_svc: metadata_svc,
		logger:               logger,
		ShaRefCounterService: NewShaRefCounterService(getCollectionNsMappingsDocKey, metadata_svc, logger),
		cachedSpecs:          make(map[string]*metadata.ReplicationSpecification),
		cachedBackfillSpecs:  make(map[string]*metadata.BackfillReplicationSpec),
		utils:                utils,
		replicationSpecSvc:   replicationSpecService,
		stopTheWorldMtx:      map[string]*sync.RWMutex{},
	}
	return ckptSvc, ckptSvc.initWithSpecs()
}

func (ckpt_svc *CheckpointsService) CheckpointsDoc(replicationId string, vbno uint16) (*metadata.CheckpointsDoc, error) {
	key := ckpt_svc.getCheckpointDocKey(replicationId, vbno)
	result, rev, err := ckpt_svc.metadata_svc.Get(key)
	if err != nil {
		return nil, err
	}

	mtx := ckpt_svc.getStopTheWorldMtx(replicationId)
	mtx.RLock()
	defer mtx.RUnlock()

	// Should exist because checkpoint manager must finish loading before allowing ckpt operations
	shaMap, _ := ckpt_svc.GetShaNamespaceMap(replicationId)

	ckpt_doc, err := ckpt_svc.constructCheckpointDoc(result, rev, shaMap)
	if err == service_def.MetadataNotFoundErr {
		ckpt_svc.logger.Errorf("Unable to construct ckpt doc from metakv given replication: %v vbno: %v key: %v",
			replicationId, vbno, key)
	}
	return ckpt_doc, err
}

func (ckpt_svc *CheckpointsService) getCheckpointCatalogKey(replicationId string) string {
	return CheckpointsCatalogKeyPrefix + base.KeyPartsDelimiter + replicationId
}

// Get a unique key to access metakv for brokenMappings
func getCollectionNsMappingsDocKey(replicationId string) string {
	return fmt.Sprintf("%v", CheckpointsKeyPrefix+base.KeyPartsDelimiter+replicationId+base.KeyPartsDelimiter+BrokenMappingKey)
}

// Get a unique key to access metakv for checkpoints
func (ckpt_svc *CheckpointsService) getCheckpointDocKey(replicationId string, vbno uint16) string {
	return fmt.Sprintf("%v%v", CheckpointsKeyPrefix+base.KeyPartsDelimiter+replicationId+base.KeyPartsDelimiter, vbno)
}

func (ckpt_svc *CheckpointsService) decodeVbnoFromCkptDocKey(ckptDocKey string) (uint16, error) {
	parts := strings.Split(ckptDocKey, base.KeyPartsDelimiter)
	vbnoStr := parts[len(parts)-1]
	vbno, err := strconv.Atoi(vbnoStr)
	if err != nil {
		return 0, err
	}
	return uint16(vbno), nil
}

func (ckpt_svc *CheckpointsService) isBrokenMappingDoc(ckptDocKey string) bool {
	return strings.Contains(ckptDocKey, BrokenMappingKey)
}

func (ckpt_svc *CheckpointsService) DelCheckpointsDocs(replicationId string) error {
	ckpt_svc.logger.Infof("DelCheckpointsDocs for replication %v...", replicationId)
	catalogKey := ckpt_svc.getCheckpointCatalogKey(replicationId)

	// No need for stop the world because replication would be deleted

	errMap := make(base.ErrorMap)
	var errMtx sync.Mutex

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err_ret := ckpt_svc.metadata_svc.DelAllFromCatalog(catalogKey)
		if err_ret != nil {
			errMsg := fmt.Sprintf("Failed to delete checkpoints docs for %v - manual clean up may be required\n", replicationId)
			ckpt_svc.logger.Errorf(errMsg)
			errMtx.Lock()
			errMap["1"] = fmt.Errorf(errMsg)
			errMtx.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		cleanMapErr := ckpt_svc.CleanupMapping(replicationId, ckpt_svc.utils)
		if cleanMapErr != nil {
			errMsg := fmt.Sprintf("DelCheckpointsDoc for brokenmapping had error: %v - manual clean up may be required", cleanMapErr)
			ckpt_svc.logger.Errorf(errMsg)
			errMtx.Lock()
			errMap["2"] = fmt.Errorf(errMsg)
			errMtx.Unlock()
		}
	}()

	wg.Wait()
	ckpt_svc.logger.Infof("DelCheckpointsDocs is done for %v\n", replicationId)
	if len(errMap) == 0 {
		return nil
	}
	return fmt.Errorf(base.FlattenErrorMap(errMap))
}

// Need to have correct accounting after deleting checkpoingsDocs
func (ckpt_svc *CheckpointsService) PostDelCheckpointsDoc(replicationId string, doc *metadata.CheckpointsDoc) (modified bool, err error) {
	if doc == nil {
		return
	}

	mtx := ckpt_svc.getStopTheWorldMtx(replicationId)
	mtx.Lock()
	defer mtx.Unlock()
	decrementerFunc, err := ckpt_svc.GetDecrementerFunc(replicationId)
	if err != nil {
		return
	}

	for _, ckptRecord := range doc.Checkpoint_records {
		if ckptRecord == nil {
			continue
		}
		shaInRecord := ckptRecord.BrokenMappingSha256
		if shaInRecord != "" {
			decrementerFunc(shaInRecord)
		}
		modified = true
	}

	return
}

func (ckpt_svc *CheckpointsService) DelCheckpointsDoc(replicationId string, vbno uint16) error {
	ckpt_svc.logger.Debugf("DelCheckpointsDoc for replication %v and vbno %v...", replicationId, vbno)
	key := ckpt_svc.getCheckpointDocKey(replicationId, vbno)
	_, rev, err := ckpt_svc.metadata_svc.Get(key)
	if err != nil {
		return err
	}

	mtx := ckpt_svc.getStopTheWorldMtx(replicationId)
	mtx.RLock()
	defer mtx.RUnlock()

	catalogKey := ckpt_svc.getCheckpointCatalogKey(replicationId)
	err = ckpt_svc.metadata_svc.DelWithCatalog(catalogKey, key, rev)
	if err != nil {
		ckpt_svc.logger.Errorf("Failed to delete checkpoints doc for replication %v and vbno %v\n", replicationId, vbno)
	} else {
		ckpt_svc.logger.Debugf("DelCheckpointsDoc is done for replication %v and vbno %v\n", replicationId, vbno)
	}
	return err
}

// in addition to upserting checkpoint record, this method may also update xattr seqno
// and target cluster version in checkpoint doc
// these operations are done in the same metakv operation to ensure that they succeed and fail together
// Returns size of the checkpoint
func (ckpt_svc *CheckpointsService) UpsertCheckpoints(replicationId string, specInternalId string, vbno uint16, ckpt_record *metadata.CheckpointRecord) (int, error) {
	ckpt_svc.logger.Debugf("Persisting checkpoint record=%v for vbno=%v replication=%v\n", ckpt_record, vbno, replicationId)
	var size int

	if ckpt_record == nil {
		return size, errors.New("nil checkpoint record")
	}
	key := ckpt_svc.getCheckpointDocKey(replicationId, vbno)
	ckpt_doc, err := ckpt_svc.CheckpointsDoc(replicationId, vbno)
	if err == service_def.MetadataNotFoundErr {
		ckpt_doc = metadata.NewCheckpointsDoc(specInternalId)
		err = nil
	}
	if err != nil {
		return size, err
	}

	added, removedRecords := ckpt_doc.AddRecord(ckpt_record)
	if !added {
		ckpt_svc.logger.Debug("the ckpt record to be added is the same as the current ckpt record in the ckpt doc. no-op.")
	} else {
		mtx := ckpt_svc.getStopTheWorldMtx(replicationId)
		mtx.RLock()
		defer mtx.RUnlock()

		ckpt_json, err := json.Marshal(ckpt_doc)
		if err != nil {
			return size, err
		}

		//always update the checkpoint without revision
		err = ckpt_svc.metadata_svc.Set(key, ckpt_json, nil)

		if err != nil {
			ckpt_svc.logger.Errorf("Failed to set checkpoint doc key=%v, err=%v\n", key, err)
		} else {
			ckpt_svc.logger.Debugf("Wrote checkpoint doc key=%v, Size=%v\n", key, ckpt_doc.Size())
			size = ckpt_doc.Size()
			err = ckpt_svc.RecordMappings(replicationId, ckpt_record, removedRecords)
			if err != nil {
				ckpt_svc.logger.Errorf("Failed to record broken mapping err=%v\n", err)
			}
		}
	}
	return size, err
}

// Upserting Checkpoint doc requires a "stop-the-world" situation where refcount needs to be re-initialized
func (ckpt_svc *CheckpointsService) UpsertCheckpointsDoc(replicationId string, ckptDocs map[uint16]*metadata.CheckpointsDoc, internalId string) (bool, error) {
	var atLeastOneWritten bool

	err := ckpt_svc.validateSpecIsValid(replicationId, internalId)
	if err != nil {
		return false, err
	}

	mtx := ckpt_svc.getStopTheWorldMtx(replicationId)
	mtx.Lock()
	defer mtx.Unlock()
	errMap := make(base.ErrorMap)
	for vbno, ckptDoc := range ckptDocs {
		key := ckpt_svc.getCheckpointDocKey(replicationId, vbno)

		ckpt_json, err := json.Marshal(ckptDoc)
		if err != nil {
			return false, err
		}

		//always update the checkpoint without revision
		err = ckpt_svc.metadata_svc.Set(key, ckpt_json, nil)
		if err != nil {
			ckpt_svc.logger.Errorf("Failed to set checkpoint doc key=%v, err=%v\n", key, err)
			errMap[fmt.Sprintf("vb %v", vbno)] = err
			continue
		} else {
			atLeastOneWritten = true
		}

		for _, ckptRecord := range ckptDoc.Checkpoint_records {
			err = ckpt_svc.RecordMappings(replicationId, ckptRecord, nil)
			if err != nil {
				errMap[fmt.Sprintf("vb %v recordMapping", vbno)] = err
			}
		}
	}

	if atLeastOneWritten {
		// If checkpoint docs have been written, we need to restart ref-count the whole checkpoint record
		// This is already done as part of the initialization process
		ckpt_svc.logger.Infof("Replication checkpoint %v has been rewritten. Reload to ensure correct mapping refcnt", replicationId)
	}

	if len(errMap) > 0 {
		return atLeastOneWritten, fmt.Errorf(base.FlattenErrorMap(errMap))
	} else {
		return atLeastOneWritten, nil
	}
}

// Ensure that one single broken mapping that will be used for most, if not all, of the checkpoints, are persisted
func (ckpt_svc *CheckpointsService) PreUpsertBrokenMapping(replicationId string, specInternalId string, oneBrokenMapping *metadata.CollectionNamespaceMapping) error {
	if oneBrokenMapping == nil {
		return nil
	}

	return ckpt_svc.RegisterMapping(replicationId, specInternalId, oneBrokenMapping)
}

func (ckpt_svc *CheckpointsService) UpsertBrokenMapping(replicationId string, specInternalId string) error {
	return ckpt_svc.UpsertMapping(replicationId, specInternalId)
}

func (ckpt_svc *CheckpointsService) LoadBrokenMappings(replicationId string) (metadata.ShaToCollectionNamespaceMap, *metadata.CollectionNsMappingsDoc, service_def.IncrementerFunc, bool, error) {
	alreadyExists := ckpt_svc.InitTopicShaCounterIfNeeded(replicationId)
	mappingsDoc, err := ckpt_svc.GetMappingsDoc(replicationId, !alreadyExists /*initIfNotFound*/)
	if err != nil {
		var emptyMap metadata.ShaToCollectionNamespaceMap
		return emptyMap, nil, nil, false, err
	}

	shaToNamespaceMap, err := ckpt_svc.GetShaToCollectionNsMap(replicationId, mappingsDoc)
	if err != nil {
		var emptyMap metadata.ShaToCollectionNamespaceMap
		return emptyMap, nil, nil, false, err
	}

	incrementerFunc, err := ckpt_svc.GetIncrementerFunc(replicationId)
	if err != nil {
		var emptyMap metadata.ShaToCollectionNamespaceMap
		return emptyMap, nil, nil, false, err
	}

	return shaToNamespaceMap, mappingsDoc, incrementerFunc, alreadyExists, nil
}

// When each vb does checkpointing and adds a new checkpoint record
// Ensure that the old, bumped out ones are refcounted correctly
// And do refcount for the newly added record's count as well
func (ckpt_svc *CheckpointsService) RecordMappings(replicationId string, ckptRecord *metadata.CheckpointRecord, removedRecords []*metadata.CheckpointRecord) error {
	incrementerFunc, err := ckpt_svc.GetIncrementerFunc(replicationId)
	if err != nil {
		return err
	}
	decrementerFunc, err := ckpt_svc.GetDecrementerFunc(replicationId)
	if err != nil {
		return err
	}

	clonedBrokenMap := ckptRecord.BrokenMappings()
	if clonedBrokenMap != nil && len(*(clonedBrokenMap)) > 0 {
		incrementerFunc(ckptRecord.BrokenMappingSha256, clonedBrokenMap)
		for _, removedRecord := range removedRecords {
			if removedRecord != nil && len(removedRecord.BrokenMappingSha256) > 0 {
				decrementerFunc(removedRecord.BrokenMappingSha256)
			}
		}
	}
	return nil
}

// Should be called non-concurrently per pipeline if brokenMappingsNeeded is true
// When brokenMappingsNeeded is true, the ckpts will have mappings populated, which takes more work
func (ckpt_svc *CheckpointsService) CheckpointsDocs(replicationId string, brokenMappingsNeeded bool) (map[uint16]*metadata.CheckpointsDoc, error) {
	checkpointsDocs := make(map[uint16]*metadata.CheckpointsDoc)
	catalogKey := ckpt_svc.getCheckpointCatalogKey(replicationId)
	ckpt_entries, err := ckpt_svc.metadata_svc.GetAllMetadataFromCatalog(catalogKey)
	if err != nil {
		return nil, err
	}

	var shaToBrokenMapping metadata.ShaToCollectionNamespaceMap
	var CollectionNsMappingsDoc *metadata.CollectionNsMappingsDoc
	var refCounterRecorder service_def.IncrementerFunc
	var ckptSvcCntIsPopulated bool
	if brokenMappingsNeeded {
		shaToBrokenMapping, CollectionNsMappingsDoc, refCounterRecorder, ckptSvcCntIsPopulated, err = ckpt_svc.LoadBrokenMappings(replicationId)
		if err != nil {
			return nil, err
		}
	}

	// This code path's initialization requires doing a level of ref counting
	// As a result, the init path needs to take place first before the rest can do read-only op
	mtx := ckpt_svc.getStopTheWorldMtx(replicationId)
	if !ckptSvcCntIsPopulated {
		mtx.Lock()
		defer mtx.Unlock()
	} else {
		mtx.RLock()
		defer mtx.RUnlock()
	}

	for _, ckpt_entry := range ckpt_entries {
		if ckpt_entry != nil {
			if ckpt_svc.isBrokenMappingDoc(ckpt_entry.Key) {
				continue
			}

			vbno, err := ckpt_svc.decodeVbnoFromCkptDocKey(ckpt_entry.Key)
			if err != nil {
				return nil, err
			}

			ckpt_doc, err := ckpt_svc.constructCheckpointDoc(ckpt_entry.Value, ckpt_entry.Rev, shaToBrokenMapping)
			if err != nil {
				if err == service_def.MetadataNotFoundErr {
					ckpt_svc.logger.Errorf("Unable to construct ckpt doc from metakv given replicationId: %v vbno: %v and key: %v",
						replicationId, vbno, ckpt_entry.Key)
					continue
				} else {
					return nil, err
				}
			} else {
				checkpointsDocs[vbno] = ckpt_doc
				if brokenMappingsNeeded && !ckptSvcCntIsPopulated {
					ckpt_svc.registerCkptDocBrokenMappings(ckpt_doc, refCounterRecorder)
				}
			}
		}
	}

	if brokenMappingsNeeded {
		err = ckpt_svc.GCDocUsingLatestCounterInfo(replicationId, CollectionNsMappingsDoc)
		if err != nil {
			ckpt_svc.logger.Errorf("Unable to GC brokenmapping cache - %v", err)
			return checkpointsDocs, err
		}

		shaToBrokenMapping, _, _, _, err = ckpt_svc.LoadBrokenMappings(replicationId)
		if err != nil {
			ckpt_svc.logger.Errorf("Unable to refresh brokenmapping cache - %v", err)
			return checkpointsDocs, err
		}

		err = ckpt_svc.InitCounterShaToActualMappings(replicationId, CollectionNsMappingsDoc.SpecInternalId, shaToBrokenMapping)
		if err == nil {
			ckpt_svc.logger.Infof("Loaded brokenMap: %v", shaToBrokenMapping)
		} else {
			ckpt_svc.logger.Warnf("Error %v trying init sha counter to mapping: %v", err, shaToBrokenMapping)
		}
	}
	return checkpointsDocs, nil
}

func (ckpt_svc *CheckpointsService) registerCkptDocBrokenMappings(ckpt_doc *metadata.CheckpointsDoc, recorder service_def.IncrementerFunc) {
	if ckpt_doc == nil || recorder == nil {
		return
	}
	for _, record := range ckpt_doc.Checkpoint_records {
		if record == nil || len(record.BrokenMappingSha256) == 0 {
			continue
		}
		recorder(record.BrokenMappingSha256, record.BrokenMappings())
	}
}

func (ckpt_svc *CheckpointsService) constructCheckpointDoc(content []byte, rev interface{}, shaToBrokenMapping metadata.ShaToCollectionNamespaceMap) (*metadata.CheckpointsDoc, error) {
	// The only time content is empty is when this is a fresh XDCR system and no checkpoints has been registered yet
	if len(content) > 0 {
		ckpt_doc := &metadata.CheckpointsDoc{}
		err := json.Unmarshal(content, ckpt_doc)
		if err != nil {
			return nil, err
		}
		if len(shaToBrokenMapping) > 0 {
			ckpt_svc.populateActualMapping(ckpt_doc, shaToBrokenMapping)
		}
		return ckpt_doc, nil
	} else {
		ckpt_svc.logger.Errorf("Unable to construct valid checkpoint due to empty checkpoint data")
		return nil, service_def.MetadataNotFoundErr
	}
}

func (ckpt_svc *CheckpointsService) populateActualMapping(doc *metadata.CheckpointsDoc, shaToActualMapping metadata.ShaToCollectionNamespaceMap) {
	if doc == nil {
		return
	}

	for _, record := range doc.Checkpoint_records {
		if record == nil {
			continue
		}
		mapping, brokenMapExists := shaToActualMapping[record.BrokenMappingSha256]
		if brokenMapExists {
			record.LoadBrokenMapping(*mapping)
		}
	}
}

// get vbnos of checkpoint docs for specified replicationId
func (ckpt_svc *CheckpointsService) GetVbnosFromCheckpointDocs(replicationId string) ([]uint16, error) {
	vbnos := make([]uint16, 0)
	catalogKey := ckpt_svc.getCheckpointCatalogKey(replicationId)
	ckpt_entries, err := ckpt_svc.metadata_svc.GetAllMetadataFromCatalog(catalogKey)
	if err != nil {
		return nil, err
	}

	for _, ckpt_entry := range ckpt_entries {
		if ckpt_entry != nil {
			vbno, err := ckpt_svc.decodeVbnoFromCkptDocKey(ckpt_entry.Key)
			if err != nil {
				return nil, err
			}
			vbnos = append(vbnos, vbno)
		}
	}
	return vbnos, nil
}

func (ckpt_svc *CheckpointsService) CollectionsManifestChangeCb(replId string, oldVal interface{}, newVal interface{}) error {
	oldManifests, ok := oldVal.(*metadata.CollectionsManifestPair)
	if !ok {
		ckpt_svc.logger.Errorf("%v expected collections manifest pair, got %v", reflect.TypeOf(oldVal))
		return base.ErrorInvalidInput
	}

	newManifests, ok := newVal.(*metadata.CollectionsManifestPair)
	if !ok {
		ckpt_svc.logger.Errorf("%v expected new collections manifest pair, got %v", reflect.TypeOf(newVal))
		return base.ErrorInvalidInput
	}

	if oldManifests.Source == nil && newManifests.Source != nil {
		// no-op
	} else if oldManifests.Source != nil && newManifests.Source == nil {
		// no-op, odd err
	} else {
		spec, err := ckpt_svc.getReplicationSpec(replId)
		if err != nil {
			ckpt_svc.logger.Warnf("Did not find spec %v from internal Cache", replId)
			return nil
		}
		err = ckpt_svc.handleManifestsChange(spec, oldManifests, newManifests)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ckpt_svc *CheckpointsService) ReplicationSpecChangeCallback(metadataId string, oldVal interface{}, newVal interface{}, wg *sync.WaitGroup) error {
	if wg != nil {
		defer wg.Done()
	}

	oldSpec, ok := oldVal.(*metadata.ReplicationSpecification)
	if !ok {
		return base.ErrorInvalidInput
	}
	newSpec, ok := newVal.(*metadata.ReplicationSpecification)
	if !ok {
		return base.ErrorInvalidInput
	}

	if oldSpec != nil && newSpec == nil {
		ckpt_svc.specsMtx.Lock()
		delete(ckpt_svc.cachedSpecs, oldSpec.Id)
		delete(ckpt_svc.stopTheWorldMtx, oldSpec.Id)
		ckpt_svc.specsMtx.Unlock()
	} else {
		ckpt_svc.specsMtx.Lock()
		ckpt_svc.cachedSpecs[newSpec.Id] = newSpec
		ckpt_svc.stopTheWorldMtx[newSpec.Id] = &sync.RWMutex{}
		ckpt_svc.specsMtx.Unlock()
	}
	return nil
}

func (ckpt_svc *CheckpointsService) BackfillReplicationSpecChangeCallback(id string, oldVal interface{}, newVal interface{}) error {
	oldSpec, ok := oldVal.(*metadata.BackfillReplicationSpec)
	if !ok {
		return base.ErrorInvalidInput
	}
	newSpec, ok := newVal.(*metadata.BackfillReplicationSpec)
	if !ok {
		return base.ErrorInvalidInput
	}

	if oldSpec != nil && newSpec == nil {
		ckpt_svc.backfillSpecsMtx.Lock()
		delete(ckpt_svc.cachedBackfillSpecs, oldSpec.Id)
		ckpt_svc.backfillSpecsMtx.Unlock()
	} else {
		ckpt_svc.backfillSpecsMtx.Lock()
		ckpt_svc.cachedBackfillSpecs[newSpec.Id] = newSpec
		ckpt_svc.backfillSpecsMtx.Unlock()
	}
	return nil
}

// If source namespaces are removed, then any mentioning of the source namespaces can be removed
// from the checkpoints
func (ckpt_svc *CheckpointsService) removeMappingFromCkptDocs(replicationId string, internalId string, sources metadata.ScopesMap) (mappingChanged bool, err error) {
	ckptDocs, _ := ckpt_svc.CheckpointsDocs(replicationId, true)
	for _, docs := range ckptDocs {
		if docs == nil {
			continue
		}
		for _, record := range docs.Checkpoint_records {
			if record == nil {
				continue
			}
			brokenMappings := record.BrokenMappings()
			if brokenMappings == nil {
				continue
			}
			toBeDelNamespace := make(metadata.CollectionNamespaceMapping)
			populateToDelNamespaces(brokenMappings, sources, toBeDelNamespace)
			if len(toBeDelNamespace) > 0 {
				incFunc, err := ckpt_svc.GetIncrementerFunc(replicationId)
				decFunc, err2 := ckpt_svc.GetDecrementerFunc(replicationId)
				if err != nil || err2 != nil {
					ckpt_svc.logger.Errorf("Unable to get increment or decrement func for %v when removing mapping", replicationId)
					continue
				}
				err = ckpt_svc.removeNamespacesFromCkpt(brokenMappings, toBeDelNamespace, record, incFunc, decFunc)
				if err == nil {
					mappingChanged = true
				}
			}
		}
	}
	if mappingChanged {
		err = ckpt_svc.UpsertMapping(replicationId, internalId)
	}
	return
}

func (ckpt_svc *CheckpointsService) removeNamespacesFromCkpt(incomingMapping *metadata.CollectionNamespaceMapping, toBeDelNamespace metadata.CollectionNamespaceMapping, record *metadata.CheckpointRecord, incFunc service_def.IncrementerFunc, decFunc service_def.DecrementerFunc) error {
	var err error
	origShaSlice, _ := incomingMapping.Sha256()
	origSha := fmt.Sprintf("%x", origShaSlice[:])
	if incomingMapping != nil && len(*incomingMapping) > 0 {
		newBrokenMapping := incomingMapping.Delete(toBeDelNamespace)
		newShaSlice, _ := newBrokenMapping.Sha256()
		newSha := fmt.Sprintf("%x", newShaSlice[:])
		err = record.LoadBrokenMapping(newBrokenMapping)
		if err != nil {
			ckpt_svc.logger.Warnf("when setting brokenMapping SHA: %v", err)
		} else {
			incFunc(newSha, &newBrokenMapping)
			if origSha != "" {
				decFunc(origSha)
			}
		}
	}
	return err
}

func populateToDelNamespaces(mappingToCheck *metadata.CollectionNamespaceMapping, sources metadata.ScopesMap, toBeDelNamespace metadata.CollectionNamespaceMapping) {
	for sourceNs, targetNsList := range *mappingToCheck {
		scope, scopeExists := sources[sourceNs.ScopeName]
		if !scopeExists {
			continue
		}
		_, collectionExists := scope.Collections[sourceNs.CollectionName]
		if !collectionExists {
			continue
		}
		// The source instance from "sources" exists in broken map and should be removed
		toBeDelNamespace.AddSingleMapping(sourceNs.CollectionNamespace, targetNsList[0])
	}
}

func (ckpt_svc *CheckpointsService) getReplicationSpec(id string) (*metadata.ReplicationSpecification, error) {
	ckpt_svc.specsMtx.RLock()
	defer ckpt_svc.specsMtx.RUnlock()
	origSpec, exists := ckpt_svc.cachedSpecs[id]
	if !exists {
		return nil, base.ErrorNotFound
	} else {
		return origSpec.Clone(), nil
	}
}

func (ckpt_svc *CheckpointsService) getBackfillReplSpec(id string) (*metadata.BackfillReplicationSpec, error) {
	ckpt_svc.backfillSpecsMtx.RLock()
	defer ckpt_svc.backfillSpecsMtx.RUnlock()
	origSpec, exists := ckpt_svc.cachedBackfillSpecs[id]
	if !exists {
		return nil, base.ErrorNotFound
	} else {
		return origSpec, nil
	}
}

func (ckpt_svc *CheckpointsService) handleManifestsChange(spec *metadata.ReplicationSpecification, oldManifests, newManifests *metadata.CollectionsManifestPair) error {
	if oldManifests.Source == nil || newManifests.Source == nil {
		return nil
	}

	// Only need to handle source changes
	_, _, removed, err := newManifests.Source.Diff(oldManifests.Source)
	if err != nil {
		ckpt_svc.logger.Errorf("Unable to diff between manifests %v and %v", oldManifests.Source.Uid(), newManifests.Source.Uid())
		return err
	}

	if len(removed) == 0 {
		return nil
	}

	// If there are source namespaces that are removed, then any of the mappings need to be removed as well
	changed, err := ckpt_svc.removeMappingFromCkptDocs(spec.Id, spec.InternalId, removed)
	if changed {
		ckpt_svc.logger.Infof("Replication %v checkpoints mapping changed due to collections manifest changes", spec.Id)
	}
	if err != nil {
		ckpt_svc.logger.Warnf("Unable to remove mappings %v from ckpt docs due to err %v", removed.String(), err)
	}

	// Do the same for any checkpoints related to backfill replication
	backfillSpec, err2 := ckpt_svc.getBackfillReplSpec(spec.Id)
	if err2 == nil && backfillSpec != nil {
		backfillSpecId := base.CompileBackfillPipelineSpecId(spec.Id)
		changed, err2 := ckpt_svc.removeMappingFromCkptDocs(backfillSpecId, backfillSpec.InternalId, removed)
		if changed {
			ckpt_svc.logger.Infof("Backfill Replication %v checkpoints mapping changed due to collections manifest changes", backfillSpecId)
		}
		if err2 != nil {
			ckpt_svc.logger.Warnf("Unable to remove mappings %v from backfill ckpt docs due to err %v", removed.String(), err2)
		}
	}
	if err != nil {
		return err
	} else if err2 != nil {
		return err2
	} else {
		return nil
	}
}

func (ckpt_svc *CheckpointsService) GetCkptsMappingsCleanupCallback(specId, specInternalId string, toBeRemoved metadata.ScopesMap) (base.StoppedPipelineCallback, base.StoppedPipelineErrCallback) {
	cb := func() error {
		_, err := ckpt_svc.removeMappingFromCkptDocs(specId, specInternalId, toBeRemoved)
		return err
	}
	errCb := func(err error) {
		// checkpoint cleanup err is fine, they'll get rolled over
		ckpt_svc.logger.Warnf("Unable to remove mappings %v from backfill ckpt docs due to err %v", toBeRemoved.String(), err)
	}
	return cb, errCb
}

func (ckpt_svc *CheckpointsService) initWithSpecs() error {
	specs, err := ckpt_svc.replicationSpecSvc.AllReplicationSpecs()
	if err != nil {
		return err
	}

	for specId, spec := range specs {
		ckpt_svc.cachedSpecs[specId] = spec
		ckpt_svc.stopTheWorldMtx[specId] = &sync.RWMutex{}
	}

	return nil
}

func (ckpt_svc *CheckpointsService) getStopTheWorldMtx(replId string) *sync.RWMutex {
	ckpt_svc.specsMtx.RLock()
	defer ckpt_svc.specsMtx.RUnlock()

	mtx, exists := ckpt_svc.stopTheWorldMtx[replId]
	if !exists {
		// potential race condition between spec deletion and access, return dummy
		return &sync.RWMutex{}
	} else {
		return mtx
	}
}

func (ckpt_svc *CheckpointsService) UpsertBrokenMappingsDoc(replicationId string, mappingDoc *metadata.CollectionNsMappingsDoc, ckptDoc map[uint16]*metadata.CheckpointsDoc, internalId string) error {
	err := ckpt_svc.validateSpecIsValid(replicationId, internalId)
	if err != nil {
		return err
	}

	mtx := ckpt_svc.getStopTheWorldMtx(replicationId)
	mtx.Lock()
	defer mtx.Unlock()

	return ckpt_svc.ReInitUsingMergedMappingDoc(replicationId, mappingDoc, ckptDoc, internalId)
}

func (ckpt_svc *CheckpointsService) validateSpecIsValid(fullReplId string, internalId string) error {
	replicationId, _ := common.DecomposeFullTopic(fullReplId)
	specCheck, err := ckpt_svc.replicationSpecSvc.ReplicationSpecReadOnly(replicationId)
	if err != nil {
		// Spec could have been deleted
		return fmt.Errorf("abort operation because spec %v is not found", replicationId)
	}
	if specCheck.InternalId != "" && internalId != "" && specCheck.InternalId != internalId {
		return fmt.Errorf("abort operation because spec %v internalId mismatch: expected %v actual %v",
			replicationId, internalId, specCheck.InternalId)
	}
	return nil
}
