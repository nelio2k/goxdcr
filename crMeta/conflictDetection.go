/*
Copyright 2024-Present Couchbase, Inc.

Use of this software is governed by the Business Source License included in
the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
file, in accordance with the Business Source License, use of this software will
be governed by the Apache License, Version 2.0, included in the file
licenses/APL2.txt.
*/

package crMeta

import (
	"fmt"

	mc "github.com/couchbase/gomemcached"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/hlv"
	"github.com/couchbase/goxdcr/log"
)

// Return values from conflict detection.
type ConflictDetectionResult uint32

const (
	CDNone     ConflictDetectionResult = iota // There was no conflict detection done. Eg: If conflict logging is not enabled or if hlv is not used.
	CDWin      ConflictDetectionResult = iota // No conflict, source wins.
	CDLose     ConflictDetectionResult = iota // No conflict, target wins.
	CDConflict ConflictDetectionResult = iota // Conflict detected between source and target doc.
	CDEqual    ConflictDetectionResult = iota // Source and target docs are equal.
	CDError    ConflictDetectionResult = iota // There was a error in conflict detection. This should ideally never happen.
)

func (cdr ConflictDetectionResult) String() string {
	switch cdr {
	case CDNone:
		return "None"
	case CDWin:
		return "Win"
	case CDLose:
		return "Lose"
	case CDConflict:
		return "Conflict"
	case CDEqual:
		return "Equal"
	}
	return "Unknown"
}

// Detect conflict by using source and target HLVs.
func DetectConflict(source *CRMetadata, target *CRMetadata) (ConflictDetectionResult, error) {
	sourceHLV := source.GetHLV()
	targetHLV := target.GetHLV()

	if sourceHLV == nil || targetHLV == nil {
		return CDError, fmt.Errorf("cannot detect conflict with nil HLVs, sourceHLV=%v, targetHLV=%v", sourceHLV, targetHLV)
	}

	c1 := sourceHLV.Contains(targetHLV)
	c2 := sourceHLV.Contains(targetHLV)
	if c1 && c2 {
		return CDEqual, nil
	} else if c1 {
		return CDWin, nil
	} else if c2 {
		return CDLose, nil
	}
	return CDConflict, nil
}

// Given source mutation and target doc response,
// Peform conflict detection if:
// 1. Conflict resolution type is CCR.
// 2. Conflict logging feature is turned on.
// Returns conflict detection result if performed,
// also returns the source and target document metadata for subsequent conflict resolution.
func DetectConflictIfNeeded(req *base.WrappedMCRequest, resp *mc.MCResponse, specs []base.SubdocLookupPathSpec, sourceId, targetId hlv.DocumentSourceId,
	xattrEnabled bool, uncompressFunc base.UncompressFunc, logger *log.CommonLogger) (ConflictDetectionResult, base.DocumentMetadata, base.DocumentMetadata, error) {

	// source dcp mutation as "source document"
	var sourceDoc *SourceDocument
	// target getMeta or subdocOp response as "target document"
	var targetDoc *TargetDocument
	// source and target document metadata like CAS, RevSeqno etc, except xattrs (hlv).
	var sourceDocMeta, targetDocMeta base.DocumentMetadata
	// source and target entire metadata including document metadata like CAS, RevSeqno etc, and xattrs like hlv.
	var targetMeta, sourceMeta *CRMetadata
	var err error

	if resp.Opcode == base.GET_WITH_META {
		// non-hlv based replications
		// No conflict detection required
		sourceDocMeta = base.DecodeSetMetaReq(req.Req)
		targetDocMeta, err = base.DecodeGetMetaResp(req.Req.Key, resp, xattrEnabled)
		if err != nil {
			err = fmt.Errorf("error decoding GET_META response for key=%v%s%v, respBody=%v%v%v, xattrEnabled=%v",
				base.UdTagBegin, req.Req.Key, base.UdTagEnd,
				base.UdTagBegin, resp.Body, base.UdTagEnd,
				xattrEnabled)
			return CDError, sourceDocMeta, targetDocMeta, err
		}

		return CDNone, sourceDocMeta, targetDocMeta, nil
	}

	if resp.Opcode != mc.SUBDOC_MULTI_LOOKUP {
		err = fmt.Errorf("unknown response %v for CR, for key=%v%s%v, req=%v%s%v, reqBody=%v%v%v", resp.Opcode,
			base.UdTagBegin, req.Req.Key, base.UdTagEnd,
			base.UdTagBegin, req.Req, base.UdTagEnd,
			base.UdTagBegin, req.Req.Body, base.UdTagEnd)
		return CDError, sourceDocMeta, targetDocMeta, err
	}

	// source document metadata
	sourceDoc = NewSourceDocument(req, sourceId)
	sourceMeta, err = sourceDoc.GetMetadata(uncompressFunc)
	if err != nil {
		err = fmt.Errorf("error decoding source mutation for key=%v%s%v, req=%v%s%v, reqBody=%v%v%v",
			base.UdTagBegin, req.Req.Key, base.UdTagEnd,
			base.UdTagBegin, req.Req, base.UdTagEnd,
			base.UdTagBegin, req.Req.Body, base.UdTagEnd)
		return CDError, sourceDocMeta, targetDocMeta, err
	}
	sourceDocMeta = *sourceMeta.docMeta

	// target document metadata
	targetDoc, err = NewTargetDocument(req.Req.Key, resp, specs, targetId, xattrEnabled, true)
	if err == base.ErrorDocumentNotFound {
		return CDNone, sourceDocMeta, targetDocMeta, err
	} else if err != nil {
		err = fmt.Errorf("error creating target document for key=%v%s%v, respBody=%v, resp=%s, specs=%v, xattrEnabled=%v",
			base.UdTagBegin, req.Req.Key, base.UdTagEnd,
			resp.Body, resp.Status, specs, xattrEnabled)
		return CDError, sourceDocMeta, targetDocMeta, err
	}
	targetMeta, err = targetDoc.GetMetadata()
	if err == base.ErrorDocumentNotFound {
		return CDNone, sourceDocMeta, targetDocMeta, err
	} else if err != nil {
		err = fmt.Errorf("error decoding target SUBDOC_MULTI_LOOKUP response for key=%v%s%v, respBody=%v, resp=%s, specs=%v, xattrEnabled=%v",
			base.UdTagBegin, req.Req.Key, base.UdTagEnd,
			resp.Body, resp.Status, specs, xattrEnabled)
		return CDError, sourceDocMeta, targetDocMeta, err
	}
	targetDocMeta = *targetMeta.docMeta

	// decide if we need to perform subdoc op to avoid target CAS rollback.
	setSubdocOpIfNeeded(sourceMeta, targetMeta, req)

	// perform conflict detection if needed
	var cdResult ConflictDetectionResult
	if true /* TODO: Only detect conflict if (a) conflict logging is on. (b) CCR conflict mode. */ {
		cdResult, err = DetectConflict(sourceMeta, targetMeta)
		if err != nil {
			err = fmt.Errorf("error detecting conflict key=%v%s%v, sourceMeta=%s, targetMeta=%s",
				base.UdTagBegin, req.Req.Key, base.UdTagEnd,
				sourceMeta, targetMeta)
			return CDError, sourceDocMeta, targetDocMeta, err
		}
	}

	if logger.GetLogLevel() >= log.LogLevelDebug {
		logger.Debugf("req=%v%s%v,reqBody=%v%v%v,resp=%s,respBody=%v,sourceMeta=%s,targetMeta=%s,cdRes=%v",
			base.UdTagBegin, req.Req, base.UdTagEnd,
			base.UdTagBegin, req.Req.Body, base.UdTagEnd,
			resp.Status, resp.Body, sourceMeta, targetMeta, cdResult,
		)
	}

	return cdResult, sourceDocMeta, targetDocMeta, nil
}
