// Copyright 2021-Present Couchbase, Inc.
//
// Use of this software is governed by the Business Source License included in
// the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
// file, in accordance with the Business Source License, use of this software
// will be governed by the Apache License, Version 2.0, included in the file
// licenses/APL2.txt.

package peerToPeer

import (
	"fmt"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/common"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/service_def"
	"io/ioutil"
	"sync"
	"time"
)

type VBMasterCheckHandler struct {
	*HandlerCommon
	finCh     chan bool
	receiveCh chan interface{}

	bucketTopologySvc service_def.BucketTopologySvc
	ckptSvc           service_def.CheckpointsService
	colManifestSvc    service_def.CollectionsManifestSvc
	backfillReplSvc   service_def.BackfillReplSvc
}

const VBMasterCheckSubscriberId = "VBMasterCheckHandler"

func NewVBMasterCheckHandler(reqCh chan interface{}, logger *log.CommonLogger, lifeCycleId string, cleanupInterval time.Duration,
	bucketTopologySvc service_def.BucketTopologySvc, ckptSvc service_def.CheckpointsService,
	collectionsManifestSvc service_def.CollectionsManifestSvc, backfillReplSvc service_def.BackfillReplSvc) *VBMasterCheckHandler {
	finCh := make(chan bool)
	handler := &VBMasterCheckHandler{
		HandlerCommon:     NewHandlerCommon(logger, lifeCycleId, finCh, cleanupInterval),
		finCh:             finCh,
		receiveCh:         reqCh,
		bucketTopologySvc: bucketTopologySvc,
		ckptSvc:           ckptSvc,
		colManifestSvc:    collectionsManifestSvc,
		backfillReplSvc:   backfillReplSvc,
	}
	return handler
}

func (v *VBMasterCheckHandler) Start() error {
	v.HandlerCommon.Start()
	go v.handler()
	return nil
}

func (v *VBMasterCheckHandler) Stop() error {
	close(v.finCh)
	return nil
}

func (h *VBMasterCheckHandler) handleRequest(req *VBMasterCheckReq) {
	// VBMasterCheck doesn't care about life cycle ID
	if req == nil {
		h.logger.Warnf("Received nil req")
		return
	}
	bucketVBsMap := req.GetBucketVBMap()
	h.logger.Infof("Received VB master check request from %v with specID %v for the following Bucket -> VBs %v", req.GetSender(), req.ReplicationId, bucketVBsMap)

	resp := req.GenerateResponse().(*VBMasterCheckResp)
	resp.Init()

	waitGrp := &sync.WaitGroup{}
	waitGrp.Add(1)
	go h.populateBucketVBMapsIntoResp(bucketVBsMap, resp, waitGrp)

	waitGrp.Add(1)
	var bgErr error
	var result map[uint16]*metadata.CheckpointsDoc
	go h.populatePipelineCkpts(common.ComposeFullTopic(req.ReplicationId, req.PipelineType), waitGrp, &bgErr, &result)

	cachedSrcManifests := make(metadata.ManifestsCache)
	cachedTgtManifests := make(metadata.ManifestsCache)
	var manifestErr error
	waitGrp.Add(1)
	go h.fetchAllManifests(req.ReplicationId, &cachedSrcManifests, &cachedTgtManifests, &manifestErr, waitGrp)

	var brokenMappingErr error
	var brokenMappingDoc metadata.CollectionNsMappingsDoc
	waitGrp.Add(1)
	go h.fetchBrokenMappingDoc(req.ReplicationId, &brokenMappingDoc, &brokenMappingErr, waitGrp)

	var backfillTasksErr error
	backfillTasks := metadata.NewVBTasksMap()
	waitGrp.Add(1)
	go h.fetchBackfillTasks(req.ReplicationId, backfillTasks, &backfillTasksErr, waitGrp)

	// Get all errors in order
	waitGrp.Wait()
	if bgErr != nil {
		h.logger.Errorf("%v", bgErr)
		resp.ErrorMsg = bgErr.Error()
		req.CallBack(resp)
		return
	}

	if manifestErr != nil {
		h.logger.Errorf("%v", manifestErr)
		resp.ErrorMsg = manifestErr.Error()
		req.CallBack(resp)
		return
	}

	if brokenMappingErr != nil {
		h.logger.Errorf("%v", brokenMappingErr)
		resp.ErrorMsg = brokenMappingErr.Error()
		req.CallBack(resp)
		return
	}

	if backfillTasksErr != nil {
		h.logger.Errorf("%v", backfillTasksErr)
		resp.ErrorMsg = backfillTasksErr.Error()
		req.CallBack(resp)
		return
	}

	err := resp.LoadPipelineCkpts(result, req.SourceBucketName)
	if err != nil {
		h.logger.Errorf("when loading pipeline ckpt into response, got %v", err)
		resp.ErrorMsg = err.Error()
		req.CallBack(resp)
		return
	}

	err = resp.LoadManifests(cachedSrcManifests, cachedTgtManifests, req.SourceBucketName)
	if err != nil {
		h.logger.Errorf("when loading manifests into response, got %v", err)
		resp.ErrorMsg = err.Error()
		req.CallBack(resp)
		return
	}

	err = resp.LoadBrokenMappingDoc(brokenMappingDoc, req.SourceBucketName)
	if err != nil {
		h.logger.Errorf("when loading brokenMappingDoc into response, got %v", err)
		resp.ErrorMsg = err.Error()
		req.CallBack(resp)
		return
	}

	err = resp.LoadBackfillTasks(backfillTasks, req.SourceBucketName)
	if err != nil {
		h.logger.Errorf("when loading brokenMappingDoc into response, got %v", err)
		resp.ErrorMsg = err.Error()
		req.CallBack(resp)
		return
	}
	toBeSent2, _ := resp.Serialize()
	ioutil.WriteFile(fmt.Sprintf("/tmp/toBeSent2_%v", counter), toBeSent2, 0644)

	// Final Callback
	fmt.Printf("NEIL DEBUG sending callback with opaque %v\n", resp.Opaque)
	handlerResult, err := req.CallBack(resp)
	if err != nil || handlerResult != nil && handlerResult.GetError() != nil {
		var handlerResultErr error
		if handlerResult != nil {
			handlerResultErr = handlerResult.GetError()
		}
		h.logger.Errorf("Unable to send resp %v to original req %v - %v %v", resp, req, err, handlerResultErr)
	}
	return
}

func (h *VBMasterCheckHandler) populateBucketVBMapsIntoResp(bucketVBsMap BucketVBMapType, resp *VBMasterCheckResp, waitGrp *sync.WaitGroup) {
	defer waitGrp.Done()
	for bucketName, vbsList := range bucketVBsMap {
		resp.InitBucket(bucketName)

		// localFeed only cares about source bucket name
		tempRef, err := metadata.NewReplicationSpecification(bucketName, "", "", "", "")
		if err != nil {
			errMsg := fmt.Sprintf("Unable to get vbsList for bucket %v - %v", bucketName, err)
			h.logger.Warnf(errMsg)
			(*resp.responsePayload)[bucketName].OverallPayloadErr = errMsg
			continue
		}
		srcNotificationCh, err := h.bucketTopologySvc.SubscribeToLocalBucketFeed(tempRef, VBMasterCheckSubscriberId)
		if err != nil {
			errMsg := fmt.Sprintf("Unable to get srcNotificationCh for bucket %v - %v", bucketName, err)
			h.logger.Warnf(errMsg)
			(*resp.responsePayload)[bucketName].OverallPayloadErr = errMsg
			continue
		}

		unsubsFunc := func() {
			err = h.bucketTopologySvc.UnSubscribeLocalBucketFeed(tempRef, VBMasterCheckSubscriberId)
			if err != nil {
				h.logger.Warnf("Unable to unsubscribe srcNotificationCh for bucket %v - %v", bucketName, err)
				// Not an error remote side cares about
			}
		}

		var latestInfo service_def.SourceNotification
		latestInfo = <-srcNotificationCh
		//select {
		// Discard the first info as it may not be the most up-to-date
		//case <-srcNotificationCh:
		//	select {
		//	case latestInfo = <-srcNotificationCh:
		// done
		//}
		//}

		// SourceVBMapRO should only contain one node
		myVBMap := latestInfo.GetSourceVBMapRO()
		var oneKey string
		for key, _ := range myVBMap {
			oneKey = key
		}
		myVbsList := myVBMap[oneKey]
		//fmt.Printf("NEIL DEBUG myVBList: %v\n", myVbsList)
		_, _, vbsIntersect := base.ComputeDeltaOfUint16Lists(myVbsList, vbsList, true)
		// Given my list and another list of VBs that I should not own,
		// if there is any intersection, then that's an issue
		if len(vbsIntersect) > 0 {
			errMsg := fmt.Sprintf("Bucket %v has VBs intersect of %v", bucketName, vbsIntersect)
			h.logger.Errorf(errMsg)
			// TODO - not setting up VBs correctly
			(*resp.responsePayload)[bucketName].RegisterVbsIntersect(vbsIntersect)
			removed, _, _ := base.ComputeDeltaOfUint16Lists(myVbsList, vbsIntersect, true)
			// Whatever are not intersected are OK
			(*resp.responsePayload)[bucketName].RegisterNotMyVBs(removed)
		} else {
			// Everything is not my VBs
			(*resp.responsePayload)[bucketName].RegisterNotMyVBs(vbsList)
		}
		unsubsFunc()
	}
}

func (h *VBMasterCheckHandler) handler() {
	for {
		select {
		case <-h.finCh:
			return
		case req := <-h.receiveCh:
			// Can be either req or response
			vbMasterReq, isReq := req.(*VBMasterCheckReq)
			vbMasterResp, isResp := req.(*VBMasterCheckResp)
			if isReq {
				h.handleRequest(vbMasterReq)
			} else if isResp {
				h.handleResponse(vbMasterResp)
			}

		}
	}
}

func (v *VBMasterCheckHandler) handleResponse(resp *VBMasterCheckResp) {
	req, retCh, found := v.GetReqAndClearOpaque(resp.GetOpaque())
	if !found {
		v.logger.Warnf("Unable to find opaque %v", resp.GetOpaque())
	}
	if retCh != nil {
		retPair := ReqRespPair{
			ReqPtr:  req,
			RespPtr: resp,
		}
		go v.sendBackSynchronously(retCh, retPair)
	}
}

func (v *VBMasterCheckHandler) populatePipelineCkpts(replSpecId string, waitGrp *sync.WaitGroup, err *error, result *map[uint16]*metadata.CheckpointsDoc) {
	defer waitGrp.Done()

	ckptDocs, opErr := v.ckptSvc.CheckpointsDocs(replSpecId, true)
	v.logger.Infof("Handler for %v retrieving CheckpointsDocs request found %v docs", replSpecId, len(ckptDocs))
	if opErr != nil {
		*err = opErr
		return
	}

	*result = ckptDocs

	for _, ckptDoc := range ckptDocs {
		records := ckptDoc.Checkpoint_records
		for _, record := range records {
			if record == nil {
				continue
			}
		}
	}

	return
}

// We depend on CollectionsManifestSvc to always keep only the minimal manifests needed
// So whatever it returns, most likely it is needed, and no need to go filter through them
func (v *VBMasterCheckHandler) fetchAllManifests(replId string, srcManifests *metadata.ManifestsCache, tgtManifests *metadata.ManifestsCache, errPtr *error, waitGrp *sync.WaitGroup) {
	defer waitGrp.Done()
	nameOnlySpec := &metadata.ReplicationSpecification{}
	nameOnlySpec.Id = replId
	src, tgt, err := v.colManifestSvc.GetAllCachedManifests(nameOnlySpec)
	if err != nil {
		*errPtr = err
		return
	}
	*srcManifests = src
	*tgtManifests = tgt
}

func (v *VBMasterCheckHandler) fetchBrokenMappingDoc(replId string, mappingDoc *metadata.CollectionNsMappingsDoc, errPtr *error, waitGrp *sync.WaitGroup) {
	defer waitGrp.Done()

	_, loadedDoc, _, _, err := v.ckptSvc.LoadBrokenMappings(replId)
	if err != nil {
		*errPtr = err
		return
	}
	if loadedDoc == nil {
		*errPtr = fmt.Errorf("Nil doc when loading brokenMapping")
		return
	}

	*mappingDoc = *loadedDoc
}

func (v *VBMasterCheckHandler) fetchBackfillTasks(replId string, backfillTasks *metadata.VBTasksMapType, backfillErr *error, waitGrp *sync.WaitGroup) {
	defer waitGrp.Done()

	backfillSpec, err := v.backfillReplSvc.BackfillReplSpec(replId)
	if err != nil {
		if err != base.ReplNotFoundErr {
			*backfillErr = err
		}
		return
	}

	if backfillSpec.VBTasksMap != nil {
		clonedTask := backfillSpec.VBTasksMap.Clone()
		*backfillTasks = *clonedTask
	}
}
