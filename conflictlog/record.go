package conflictlog

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	SourcePrefix = "src"
	TargetPrefix = "tgt"
	CRDPrefix    = "crd"
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
	Body   []byte
	Vbno   uint16
	Seqno  uint64
	VbUUID uint64
	// IsSource bool // If true, Body above will contain xattrs.
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

func (d *DocInfo) String() string {
	if d == nil {
		return "{}"
	}

	return fmt.Sprintf("{id=%s,node=%s,bucket=%s,cluster=%s,deleted=%v,collection=%s,scope=%s,expiry=%v,flags=%v,cas=%v,revId=%v,datatype=%v,vb=%v,seqno=%v,vbuuid=%v,xattrs=%s}",
		d.Id, d.NodeId, d.BucketUUID, d.ClusterUUID, d.IsDeleted, d.Collection, d.Scope, d.Expiry, d.Flags,
		d.Cas, d.RevSeqno, d.Datatype, d.Vbno, d.Seqno, d.VbUUID, &d.Xattrs)
}

// ConflictRecord has the all the details of the detected conflict
// which are needed to be persisted
type ConflictRecord struct {
	Id            string  `json:"id"`
	DocId         string  `json:"docId"`
	ReplicationId string  `json:"replId"`
	Source        DocInfo `json:"source"`
	Target        DocInfo `json:"target"`
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

// populates some derived data.
func (r *ConflictRecord) PopulateData(replicationId string) {
	if r == nil {
		return
	}

	now := time.Now().UnixNano()

	r.ReplicationId = replicationId
	r.PopulateSourceDocId(now)
	r.PopulateTargetDocId(now)
	r.PopulateCRDocId(now)
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
			r.Source.BucketUUID, r.Source.Vbno, r.Source.Seqno, r.Source.VbUUID,
			r.Target.BucketUUID, r.Target.Vbno, r.Target.Seqno, r.Target.VbUUID,
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
