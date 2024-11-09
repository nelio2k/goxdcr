package conflictlog

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/couchbase/goxdcr/v8/common"
)

// DocInfo is the subset of the information about a doc needed for
// conflict logging
type DocInfo struct {
	// Generated doc id
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
	Xattrs      Xattrs `json:"xattrs"`

	// Note: The following will not be serialized to json
	Body     []byte `json:"-"`
	Datatype uint8  `json:"-"` // datatype of the Body above.
	VBNo     uint16 `json:"-"`
	Seqno    uint64 `json:"-"`
	VBUUID   uint64 `json:"-"`
}

func (d *DocInfo) String() string {
	if d == nil {
		return "{}"
	}

	return fmt.Sprintf("{id=%s,node=%s,bucket=%s,cluster=%s,deleted=%v,collection=%s,scope=%s,expiry=%v,flags=%v,cas=%v,revId=%v,datatype=%v,vb=%v,seqno=%v,vbuuid=%v,xattrs=%s}",
		d.Id, d.NodeId, d.BucketUUID, d.ClusterUUID, d.IsDeleted, d.Collection, d.Scope, d.Expiry, d.Flags,
		d.Cas, d.RevSeqno, d.Datatype, d.VBNo, d.Seqno, d.VBUUID, &d.Xattrs)
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
// which are needed to be persisted.
// Body, Id and ReplicationId are derived or generated during conflict logging.
type ConflictRecord struct {
	Timestamp string `json:"timestamp"`
	// generated by conflict log
	Id string `json:"id"`
	// original doc key
	DocId         string  `json:"docId"`
	ReplicationId string  `json:"replId"`
	Source        DocInfo `json:"srcDoc"`
	Target        DocInfo `json:"tgtDoc"`

	// Note: The following will not be serialized to json
	Body                []byte              `json:"-"`
	Datatype            uint8               `json:"-"`
	StartTime           time.Time           `json:"-"`
	OriginatingPipeline common.PipelineType `json:"-"`
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
func (r *ConflictRecord) PopulateData(replicationId string, pipelineType common.PipelineType) error {
	if r == nil {
		return fmt.Errorf("nil CRD")
	}

	r.ReplicationId = common.ComposeFullTopic(replicationId, pipelineType)
	now := time.Now().UnixNano()

	uniqKey := r.GenerateUniqKey()
	r.PopulateDocIds(uniqKey, now)

	body, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("error marshaling CRD, err=%v", err)
	}
	r.Body = body

	// source and target document have been formatted with conflict logging xattrs as needed for SetWMeta
	// Do it now for CRD.
	newBody, newDatatype, err := InsertConflictXattrToBody(r.Body, r.Datatype)
	if err != nil {
		return fmt.Errorf("error inserting xattr to CRD, err=%v", err)
	}
	r.Body = newBody
	r.Datatype = newDatatype

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
// The following function generates and returns <SHA> part of the above doc keys.
func (r *ConflictRecord) GenerateUniqKey() string {
	uniqKey := []byte(
		fmt.Sprintf("%s_%v_%v_%v_%s_%v_%v_%v",
			r.Source.BucketUUID, r.Source.VBNo, r.Source.Seqno, r.Source.VBUUID,
			r.Target.BucketUUID, r.Target.VBNo, r.Target.Seqno, r.Target.VBUUID,
		))

	sha256Hash := sha256.Sum256(uniqKey)
	sha256HashHex := hex.EncodeToString(sha256Hash[:])
	return sha256HashHex
}

// uniqKey is the output of GenerateUniqKey
func (r *ConflictRecord) PopulateDocIds(uniqKey string, now int64) {
	r.Source.Id = fmt.Sprintf("%s_%v_%s", SourcePrefix, now, uniqKey)
	r.Target.Id = fmt.Sprintf("%s_%v_%s", TargetPrefix, now, uniqKey)
	r.Id = fmt.Sprintf("%s_%v_%s", CRDPrefix, now, uniqKey)
}

func (r *ConflictRecord) ResetStartTime() {
	r.StartTime = time.Now()
}
