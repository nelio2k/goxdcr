package conflictlog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/couchbase/goxdcr/base"
)

const (
	SourcePrefix    string = "src"
	TargetPrefix    string = "tgt"
	CRDPrefix       string = "crd"
	TimestampFormat string = "2006-01-02T15:04:05.000Z07:00" // YYYY-MM-DDThh:mm:ss:SSSZ format
	// max increase in document body size after adding `"_xdcr_conflict":true` xattr.
	MaxBodyIncrease int = 4 /* entire xattr section size */ +
		4 /* _xdcr_conflict xattr size */ +
		len(base.ConflictLoggingXattrKey) + 2 /* quotes for the xattr key string */ +
		len(base.ConflictLoggingXattrVal) + 2 /* null terminators one each after key and value */
)

// Conflict is an abstraction over conflict record
type Conflict interface {
	// Scope is source bucket's scope
	Scope() string

	// Collection is source bucket's scope
	Collection() string
}

// DocInfo is the subset of the information about a doc needed for
// conflict logging
type DocInfo struct {
	Id          string `json:"id"`
	NodeId      string `json:"nodeId"`
	BucketUUID  string `json:"bucketUUID"`
	ClusterUUID string `json:"clusterUUID"`
	IsDeleted   bool   `json:"isDeleted"`
	Collection  string `json:"collection"`
	Scope       string `json:"scope"`
	Expiry      uint32 `json:"expiry"`
	Flags       uint32 `json:"flags"`
	Cas         uint64 `json:"cas"`
	RevSeqno    uint64 `json:"revSeqno"`
	Datatype    uint8  `json:"datatype"`
	Xattrs      Xattrs `json:"xattrs"`

	// Note: The following will not be serialized to json
	body   []byte // document body.
	vbno   uint16
	seqno  uint64
	vbUUID uint64
}

func (d *DocInfo) String() string {
	if d == nil {
		return "{}"
	}

	return fmt.Sprintf("{id=%s,node=%s,bucket=%s,cluster=%s,deleted=%v,collection=%s,scope=%s,expiry=%v,flags=%v,cas=%v,revId=%v,datatype=%v,vb=%v,seqno=%v,vbuuid=%v,xattrs=%s}",
		d.Id, d.NodeId, d.BucketUUID, d.ClusterUUID, d.IsDeleted, d.Collection, d.Scope, d.Expiry, d.Flags,
		d.Cas, d.RevSeqno, d.Datatype, d.vbno, d.seqno, d.vbUUID, &d.Xattrs)
}

// xattrs to be logged.
// These are the system xattrs used as metadata for conflict resolution.
type Xattrs struct {
	Hlv  string `json:"_vv"`
	Sync string `json:"_sync"`
	Mou  string `json:"_mou"`
}

func (x *Xattrs) String() string {
	if x == nil {
		return "{}"
	}

	return fmt.Sprintf("{Hlv=%s,Sync=%s,Mou=%s}", x.Hlv, x.Sync, x.Mou)
}

// ConflictRecord has the all the details of the detected conflict
// which are needed to be persisted
type ConflictRecord struct {
	Timestamp     string `json:"timestamp"`
	Id            string `json:"id"`
	DocId         string `json:"docId"`
	ReplicationId string `json:"replId"`
	Source        DocInfo
	Target        DocInfo

	// Note: The following will not be serialized to json
	datatype uint8
	body     []byte
}

// Id and ReplicationId are derived/generated during conflict logging.
func NewConflictRecord(key string, srcScopeName, tgtScopeName, srcCollectionName, tgtCollectionName, srcBucketUUID, tgtBucketUUID, srcClusterUUID, tgtClusterUUID, srcHostname, tgtHostname string,
	srcExpiry, tgtExpiry, srcFlags, tgtFlags uint32, srcCas, tgtCas, srcRevSeq, tgtRevSeq uint64, srcDatatype, tgtDatatype uint8, srcDeletion, tgtDeletion bool,
	srcHlv, tgtHlv, srcSync, tgtSync, srcMou, tgtMou string, srcVbno, tgtVbno uint16, srcVbuuid, tgtVbuuid, srcSeqno, tgtSeqno uint64, srcDocBody, tgtDocBody []byte) ConflictRecord {
	crd := ConflictRecord{
		Timestamp: time.Now().Format(TimestampFormat),
		DocId:     key,
		Source: DocInfo{
			Scope:       srcScopeName,
			Collection:  srcCollectionName,
			BucketUUID:  srcBucketUUID,
			ClusterUUID: srcClusterUUID,
			NodeId:      srcHostname,
			Expiry:      srcExpiry,
			Flags:       srcFlags,
			Datatype:    srcDatatype,
			Cas:         srcCas,
			RevSeqno:    srcRevSeq,
			IsDeleted:   srcDeletion,
			Xattrs: Xattrs{
				Hlv:  srcHlv,
				Sync: srcSync,
				Mou:  srcMou,
			},
			body:   srcDocBody,
			vbno:   srcVbno,
			vbUUID: srcVbuuid,
			seqno:  srcSeqno,
		},
		Target: DocInfo{
			Scope:       tgtScopeName,
			Collection:  tgtCollectionName,
			BucketUUID:  tgtBucketUUID,
			ClusterUUID: tgtClusterUUID,
			NodeId:      tgtHostname,
			Expiry:      tgtExpiry,
			Flags:       tgtFlags,
			Datatype:    tgtDatatype,
			Cas:         tgtCas,
			RevSeqno:    tgtRevSeq,
			IsDeleted:   tgtDeletion,
			Xattrs: Xattrs{
				Hlv:  tgtHlv,
				Sync: tgtSync,
				Mou:  tgtMou,
			},
			body:   tgtDocBody,
			vbno:   tgtVbno,
			vbUUID: tgtVbuuid,
			seqno:  tgtSeqno,
		},
		datatype: base.JSONDataType,
	}

	return crd
}

func (r *ConflictRecord) Scope() string {
	return r.Source.Scope
}

func (r *ConflictRecord) Collection() string {
	return r.Source.Collection
}

func (r *ConflictRecord) String() string {
	if r == nil {
		return "{}"
	}

	return fmt.Sprintf("{id=%s,docId=%s,replId=%s,source=%s,target=%s}",
		r.Id, r.DocId, r.ReplicationId, &r.Source, &r.Target)
}

func (r *ConflictRecord) SmallString() string {
	if r == nil {
		return "{}"
	}

	return fmt.Sprintf("{source=%s.%s,target=%s.%s}",
		r.Source.Scope, r.Source.Collection, r.Target.Scope, r.Target.Collection)
}

// populates the following:
// 1. document Ids.
// 2. CRD body.
// 3. xattr for all three conflict records.
func (r *ConflictRecord) PopulateData(replicationId string) error {
	if r == nil {
		return fmt.Errorf("nil CRD")
	}

	now := time.Now().UnixNano()

	r.ReplicationId = replicationId

	r.PopulateSourceDocId(now)
	r.PopulateTargetDocId(now)
	r.PopulateCRDocId(now)

	body, err := json.Marshal(r)
	if err != nil {
		return err
	}
	r.body = body

	err = r.InsertConflictXattr()
	if err != nil {
		return err
	}

	return nil
}

// A pair of mutations can be uniquely identified by:
// source bucket UUID, source vbno, source seqno, source vb failover UUID,
// target bucket UUID, target vbno, target seqno, target vb failover UUID.
// By using SHA256 on these properties in strict sequence, we generate a hex string <SHA>.
// Thus, the CRD documents would be:
// crd_<timestamp>_<SHA> - the CRD document.
// src_<timestamp>_<SHA> - the source document that caused the conflict.
// tgt_<timestamp>_<SHA> - the target document that caused the conflict.
func (r *ConflictRecord) GenerateUniqHash() string {
	uniqKey := []byte(
		fmt.Sprintf("%s_%v_%v_%v_%s_%v_%v_%v",
			r.Source.BucketUUID, r.Source.vbno, r.Source.seqno, r.Source.vbUUID,
			r.Target.BucketUUID, r.Target.vbno, r.Target.seqno, r.Target.vbUUID,
		))

	sha256Hash := sha256.Sum256(uniqKey)
	sha256HashHex := hex.EncodeToString(sha256Hash[:])
	return sha256HashHex
}

func (r *ConflictRecord) PopulateSourceDocId(now int64) {
	uniqKey := r.GenerateUniqHash()
	r.Source.Id = fmt.Sprintf("%s_%v_%s", SourcePrefix, now, uniqKey)
}

func (r *ConflictRecord) PopulateTargetDocId(now int64) {
	uniqKey := r.GenerateUniqHash()
	r.Target.Id = fmt.Sprintf("%s_%v_%s", TargetPrefix, now, uniqKey)
}

func (r *ConflictRecord) PopulateCRDocId(now int64) {
	uniqKey := r.GenerateUniqHash()
	r.Id = fmt.Sprintf("%s_%v_%s", CRDPrefix, now, uniqKey)
}

// inserts the "_xdcr_conflict": true xattr to avoid them replicating.
func (r *ConflictRecord) InsertConflictXattr() error {
	// add it to the source body
	newBody, newDatatype, err := insertConflictXattrToBody(r.Source.body, r.Source.Datatype)
	if err != nil {
		return fmt.Errorf("error inserting xattr to source body, err=%v", err)
	}
	r.Source.body = newBody
	r.Source.Datatype = newDatatype

	// add it to the target body
	newBody, newDatatype, err = insertConflictXattrToBody(r.Target.body, r.Target.Datatype)
	if err != nil {
		return fmt.Errorf("error inserting xattr to target body, err=%v", err)
	}
	r.Target.body = newBody
	r.Target.Datatype = newDatatype

	// add it to the CRD body
	newBody, newDatatype, err = insertConflictXattrToBody(r.body, r.datatype)
	if err != nil {
		return fmt.Errorf("error inserting xattr to CRD, err=%v", err)
	}
	r.body = newBody
	r.datatype = newDatatype

	return nil
}

// Inserts "_xdcr_conflict": true to the input byte slice.
// If any error occurs, the original body is returned.
// Otherwise, returns new body and new datatype after xattr is successfully added.
func insertConflictXattrToBody(body []byte, datatype uint8) ([]byte, uint8, error) {
	newbodyLen := len(body) + MaxBodyIncrease
	// TODO - Use datapool.
	newbody := make([]byte, newbodyLen)

	xattrComposer := base.NewXattrComposer(newbody)

	if base.HasXattr(datatype) {
		// insert the already existing xattrs
		it, err := base.NewXattrIterator(body)
		if err != nil {
			return body, datatype, err
		}

		for it.HasNext() {
			key, val, err := it.Next()
			if err != nil {
				return body, datatype, err
			}
			err = xattrComposer.WriteKV(key, val)
			if err != nil {
				return body, datatype, err
			}
		}
	}

	err := xattrComposer.WriteKV(base.ConflictLoggingXattrKeyBytes, base.ConflictLoggingXattrValBytes)
	if err != nil {
		return body, datatype, err
	}

	docWithoutXattr := base.FindDocBodyWithoutXattr(body, datatype)
	out, atLeastOneXattr := xattrComposer.FinishAndAppendDocValue(docWithoutXattr, nil, nil)

	if atLeastOneXattr {
		datatype = datatype | base.PROTOCOL_BINARY_DATATYPE_XATTR
	} else {
		// odd - shouldn't happen.
		datatype = datatype & ^(base.PROTOCOL_BINARY_DATATYPE_XATTR)
	}

	body = nil // no use of this body anymore, set to nil to help GC quicker.

	return out, datatype, err
}
