package metadata

import (
	"fmt"
)

const (
	//the maximum number of checkpoints ketp in the file
	MaxCheckpointsKept int = 100

	FailOverUUID        string = "failover_uuid"
	Seqno               string = "seqno"
	DcpSnapshotSeqno    string = "dcp_snapshot_seqno"
	DcpSnapshotEndSeqno string = "dcp_snapshot_end_seqno"
	TargetVbUUID        string = "target_vb_uuid"
	TargetSeqno        string = "target_seqno"
)

type CheckpointRecord struct {
	//source vbucket failover uuid
	Failover_uuid uint64 `json:"failover_uuid"`
	//source vbucket high sequence number
	Seqno uint64 `json:"seqno"`
	//source snapshot start sequence number
	Dcp_snapshot_seqno uint64 `json:"dcp_snapshot_seqno"`
	//source snapshot end sequence number
	Dcp_snapshot_end_seqno uint64 `json:"dcp_snapshot_end_seqno"`
	//target vb uuid
	Target_vb_uuid uint64 `json:"target_vb_uuid"`
	//target vb high sequence number
	Target_Seqno uint64 `json:"target_seqno"`
}

func (ckpt_record *CheckpointRecord) String() string{
	return fmt.Sprintf("{Failover_uuid=%v; Seqno=%v; Dcp_snapshot_seqno=%v; Dcp_snapshot_end_seqno=%v; Target_vb_uuid=%v; Commitopaque=%v}", 
	ckpt_record.Failover_uuid, ckpt_record.Seqno, ckpt_record.Dcp_snapshot_seqno, ckpt_record.Dcp_snapshot_end_seqno, ckpt_record.Target_vb_uuid,
	ckpt_record.Target_Seqno)
}

type CheckpointsDoc struct {
	//keep 100 checkpoint record
	Checkpoint_records []*CheckpointRecord `json:"checkpoints"`

	//revision number
	Revision interface{}
}

func (ckpt *CheckpointRecord) ToMap() map[string]interface{} {
	ckpt_record_map := make(map[string]interface{})
	ckpt_record_map[FailOverUUID] = ckpt.Failover_uuid
	ckpt_record_map[Seqno] = ckpt.Seqno
	ckpt_record_map[DcpSnapshotSeqno] = ckpt.Dcp_snapshot_seqno
	ckpt_record_map[DcpSnapshotEndSeqno] = ckpt.Dcp_snapshot_end_seqno
	ckpt_record_map[TargetVbUUID] = ckpt.Target_vb_uuid
	ckpt_record_map[TargetSeqno] = ckpt.Target_Seqno
	return ckpt_record_map
}

func NewCheckpointsDoc() *CheckpointsDoc {
	ckpt_doc := &CheckpointsDoc{Checkpoint_records: []*CheckpointRecord {},
		Revision: nil}

	for i := 0; i < MaxCheckpointsKept; i++ {
		ckpt_doc.Checkpoint_records = append(ckpt_doc.Checkpoint_records, nil)
	}

	return ckpt_doc
}

//Not currentcy safe. It should be used by one goroutine only
func (ckptsDoc *CheckpointsDoc) AddRecord(record *CheckpointRecord) {
	if len(ckptsDoc.Checkpoint_records) > 0 {
		for i := len(ckptsDoc.Checkpoint_records) - 1; i >= 0; i-- {
			if i+1 < MaxCheckpointsKept {
				ckptsDoc.Checkpoint_records[i+1] = ckptsDoc.Checkpoint_records[i]
			}
		}
		ckptsDoc.Checkpoint_records[0] = record
	} else {
		ckptsDoc.Checkpoint_records = append(ckptsDoc.Checkpoint_records, record)
	}
}
