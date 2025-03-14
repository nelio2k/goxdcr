/*
Copyright 2020-Present Couchbase, Inc.

Use of this software is governed by the Business Source License included in
the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
file, in accordance with the Business Source License, use of this software will
be governed by the Apache License, Version 2.0, included in the file
licenses/APL2.txt.
*/

package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"testing"

	mcc "github.com/couchbase/gomemcached/client"
	"github.com/couchbase/goxdcr/v8/base"
	"github.com/stretchr/testify/assert"
)

func TestCheckpointDocMarshaller(t *testing.T) {
	fmt.Println("============== Test case start: TestCheckpointDocMarshaller =================")
	defer fmt.Println("============== Test case end: TestCheckpointDocMarshaller =================")
	assert := assert.New(t)

	vbUuidAndTimestamp := &TargetVBUuidAndTimestamp{
		Target_vb_uuid: "abc",
		Startup_time:   "012",
	}
	newCkptRecord := CheckpointRecord{
		SourceVBTimestamp: SourceVBTimestamp{
			Failover_uuid:                0,
			Seqno:                        1,
			Dcp_snapshot_seqno:           2,
			Dcp_snapshot_end_seqno:       3,
			SourceManifestForDCP:         7,
			SourceManifestForBackfillMgr: 8,
		},
		TargetVBTimestamp: TargetVBTimestamp{
			Target_vb_opaque:    vbUuidAndTimestamp,
			Target_Seqno:        4,
			TargetManifest:      9,
			BrokenMappingSha256: "",
			brokenMappings:      nil,
		},
		SourceVBCounters: SourceVBCounters{
			FilteredItemsCnt:  5,
			FilteredFailedCnt: 6,
			CasPoisonCnt:      1,
		},
		TargetPerVBCounters: TargetPerVBCounters{
			GuardrailResidentRatioCnt: 100,
			GuardrailDataSizeCnt:      200,
			GuardrailDiskSpaceCnt:     300,
		},
	}

	brokenMap := make(CollectionNamespaceMapping)
	ns1, err := base.NewCollectionNamespaceFromString("s1.col1")
	assert.Nil(err)
	ns2, err := base.NewCollectionNamespaceFromString("s1.col2")
	assert.Nil(err)

	brokenMap.AddSingleMapping(&ns1, &ns2)
	ckptRecord2 := CheckpointRecord{
		SourceVBTimestamp: SourceVBTimestamp{
			Failover_uuid:                0,
			Seqno:                        1,
			Dcp_snapshot_seqno:           2,
			Dcp_snapshot_end_seqno:       3,
			SourceManifestForDCP:         7,
			SourceManifestForBackfillMgr: 8,
		},
		TargetVBTimestamp: TargetVBTimestamp{
			Target_vb_opaque:    vbUuidAndTimestamp,
			Target_Seqno:        4,
			TargetManifest:      9,
			BrokenMappingSha256: "",
			brokenMappings:      brokenMap,
		},
		SourceVBCounters: SourceVBCounters{
			FilteredItemsCnt:  5,
			FilteredFailedCnt: 6,
			CasPoisonCnt:      2,
		},
		TargetPerVBCounters: TargetPerVBCounters{
			GuardrailResidentRatioCnt: 50,
			GuardrailDataSizeCnt:      100,
			GuardrailDiskSpaceCnt:     200,
		},
	}
	assert.Nil(ckptRecord2.PopulateBrokenMappingSha())

	ckpt_doc := NewCheckpointsDoc("testInternalId")
	added, _ := ckpt_doc.AddRecord(&newCkptRecord)
	assert.True(added)
	// Adding newCkptRecord at this point should not have brokenMappings
	brokenMappings := ckpt_doc.Checkpoint_records[0].BrokenMappings()
	assert.Len(brokenMappings, 0)
	assert.True(ckpt_doc.Checkpoint_records[0].SameAs(&newCkptRecord))

	added, _ = ckpt_doc.AddRecord(&ckptRecord2)
	assert.True(added)
	assert.True(ckpt_doc.Checkpoint_records[0].SameAs(&ckptRecord2))
	assert.True(ckpt_doc.Checkpoint_records[1].SameAs(&newCkptRecord))

	brokenMappings = ckpt_doc.Checkpoint_records[0].BrokenMappings()
	assert.Len(brokenMappings, 1)
	brokenMappings = ckpt_doc.Checkpoint_records[1].BrokenMappings()
	assert.Len(brokenMappings, 0)

	marshalledData, err := json.Marshal(ckpt_doc)
	assert.Nil(err)

	ckptDocCompressed, shaMapCompressed, err := ckpt_doc.SnappyCompress()
	assert.Nil(err)
	// Make sure shaMapCompressed contain only one sha, which is record 0's
	assert.Len(shaMapCompressed, 1)
	var key string
	for k, _ := range shaMapCompressed {
		key = k
	}
	assert.Equal(key, ckptRecord2.BrokenMappingSha256)

	var checkDoc CheckpointsDoc
	err = json.Unmarshal(marshalledData, &checkDoc)
	assert.Nil(err)
	shaToBrokenMap := make(ShaToCollectionNamespaceMap)
	brokenmapSha, err := brokenMap.Sha256()
	assert.Nil(err)
	shaToBrokenMap[fmt.Sprintf("%s", brokenmapSha)] = &brokenMap

	assert.Equal(5, len(checkDoc.Checkpoint_records))
	assert.NotNil(checkDoc.Checkpoint_records[1])
	assert.True(checkDoc.Checkpoint_records[1].SameAs(&newCkptRecord))
	assert.Equal(newCkptRecord.SourceManifestForBackfillMgr, checkDoc.Checkpoint_records[1].SourceManifestForBackfillMgr)
	assert.NotNil(checkDoc.Checkpoint_records[0])
	assert.True(checkDoc.Checkpoint_records[0].SameAs(&ckptRecord2))

	var decompressCheck CheckpointsDoc
	brokenMap1, globalInfoShaMap, err := SnappyDecompressShaMap(shaMapCompressed)
	assert.Nil(err)
	assert.Nil(decompressCheck.SnappyDecompress(ckptDocCompressed, brokenMap1, globalInfoShaMap))
	assert.Equal(5, len(decompressCheck.Checkpoint_records))
	assert.NotNil(decompressCheck.Checkpoint_records[1])
	assert.True(decompressCheck.Checkpoint_records[1].SameAs(&newCkptRecord))
	assert.Equal(newCkptRecord.SourceManifestForBackfillMgr, decompressCheck.Checkpoint_records[1].SourceManifestForBackfillMgr)
	assert.NotNil(decompressCheck.Checkpoint_records[0])
	assert.True(decompressCheck.Checkpoint_records[0].SameAs(&ckptRecord2))
}

func TestCheckpointSortBySeqno(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestCheckpointSortBySeqno =================")
	defer fmt.Println("============== Test case end: TestCheckpointSortBySeqno =================")

	var unsortedList CheckpointRecordsList

	validFailoverLog := uint64(1234)
	validFailoverLog2 := uint64(12345)
	failoverLog := &mcc.FailoverLog{[2]uint64{validFailoverLog, 0}, [2]uint64{validFailoverLog2, 0}}
	//var invalidFailoverLog uint64 = "2345"

	earlySeqno := uint64(100)
	laterSeqno := uint64(200)
	latestSeqno := uint64(300)

	record := &CheckpointRecord{
		SourceVBTimestamp: SourceVBTimestamp{
			Failover_uuid: validFailoverLog,
			Seqno:         earlySeqno,
		},
		TargetVBTimestamp: TargetVBTimestamp{
			Target_vb_opaque: nil,
			Target_Seqno:     0,
		},
	}
	record2 := &CheckpointRecord{
		SourceVBTimestamp: SourceVBTimestamp{
			Failover_uuid: validFailoverLog,
			Seqno:         laterSeqno,
		},
	}
	record3 := &CheckpointRecord{
		SourceVBTimestamp: SourceVBTimestamp{
			Failover_uuid: validFailoverLog,
			Seqno:         latestSeqno,
		},
	}

	unsortedList = append(unsortedList, record)
	unsortedList = append(unsortedList, record2)
	unsortedList = append(unsortedList, record3)

	toSortList := unsortedList.PrepareSortStructure(failoverLog, nil)
	sort.Sort(toSortList)
	outputList := toSortList.ToRegularList()
	for i := 0; i < len(outputList)-1; i++ {
		assert.True(outputList[i].Seqno > outputList[i+1].Seqno)
	}
}

func TestCheckpointSortByFailoverLog(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestCheckpointSortByFailoverLog =================")
	defer fmt.Println("============== Test case end: TestCheckpointSortByFailoverLog =================")

	var unsortedList CheckpointRecordsList

	validFailoverLog := uint64(1234)
	validFailoverLog2 := uint64(12345)
	failoverLog := &mcc.FailoverLog{[2]uint64{validFailoverLog, 0}, [2]uint64{validFailoverLog2, 0}}

	recordA := &CheckpointRecord{SourceVBTimestamp: SourceVBTimestamp{Failover_uuid: validFailoverLog2}}
	recordB := &CheckpointRecord{SourceVBTimestamp: SourceVBTimestamp{Failover_uuid: validFailoverLog}}

	unsortedList = append(unsortedList, recordA)
	unsortedList = append(unsortedList, recordB)

	toSortList := unsortedList.PrepareSortStructure(failoverLog, nil)
	sort.Sort(toSortList)
	outputList := toSortList.ToRegularList()
	// Failover log that shows up first is higher priority
	assert.Equal(validFailoverLog, outputList[0].Failover_uuid)
	assert.Equal(validFailoverLog2, outputList[1].Failover_uuid)
}

func TestCheckpointSortByFailoverLogRealData(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("============== Test case start: TestCheckpointSortByFailoverLogRealData =================")
	defer fmt.Println("============== Test case end: TestCheckpointSortByFailoverLogRealData =================")

	testDataDir := "./testData/sortCkptData"
	srcFailoverLogFile := testDataDir + "/srcFailoverLogs.json"
	tgtFailoverLogFile := testDataDir + "/tgtFailoverLogs.json"
	ckptsFile := testDataDir + "/currDocs.json"
	peersFile := testDataDir + "/peerDocs.json"

	srcFailoverLogBytes, err := ioutil.ReadFile(srcFailoverLogFile)
	assert.Nil(err)
	tgtFailoverLogBytes, err := ioutil.ReadFile(tgtFailoverLogFile)
	assert.Nil(err)
	ckptBytes, err := ioutil.ReadFile(ckptsFile)
	assert.Nil(err)
	peerBytes, err := ioutil.ReadFile(peersFile)
	assert.Nil(err)

	srcFailoverLogs := make(map[uint16]*mcc.FailoverLog)
	tgtFailoverLogs := make(map[uint16]*mcc.FailoverLog)
	ckpts := make(map[uint16]*CheckpointsDoc)
	peerCkpts := make(map[uint16]*CheckpointsDoc)

	assert.Nil(json.Unmarshal(srcFailoverLogBytes, &srcFailoverLogs))
	assert.Nil(json.Unmarshal(tgtFailoverLogBytes, &tgtFailoverLogs))
	assert.Nil(json.Unmarshal(ckptBytes, &ckpts))
	assert.Nil(json.Unmarshal(peerBytes, &peerCkpts))

	var nonNilDocCounter int
	var nonConformedVBList []uint16
	// The tgtFailoverLog captured above is nil, but should still work
	for vb, ckptDocs := range ckpts {
		srcFailoverLog := srcFailoverLogs[vb]
		assert.NotNil(ckptDocs)
		assert.NotNil(srcFailoverLog)

		if ckptDocs == nil || ckptDocs.Checkpoint_records == nil {
			continue
		}

		curRecordToSort := ckptDocs.Checkpoint_records.PrepareSortStructure(srcFailoverLog, nil)
		if curRecordToSort == nil {
			continue
		}

		var combinedRecords CheckpointSortRecordsList
		combinedRecords = append(combinedRecords, curRecordToSort...)
		assert.NotNil(combinedRecords[0])

		peerRecord, exists := peerCkpts[vb]
		if exists && peerRecord.Checkpoint_records != nil {
			peerRecordToSort := peerRecord.Checkpoint_records.PrepareSortStructure(srcFailoverLog, nil)
			combinedRecords = append(combinedRecords, peerRecordToSort...)
		}

		sort.Sort(combinedRecords)
		regList := combinedRecords.ToRegularList()
		ckpts[vb].Checkpoint_records = regList
		var prevSeqno uint64 = math.MaxUint64
		for _, ckptRecord := range regList {
			if prevSeqno == math.MaxUint64 {
				prevSeqno = ckptRecord.Seqno
			} else {
				if ckptRecord.Seqno > prevSeqno {
					nonConformedVBList = append(nonConformedVBList, vb)
				}
			}
		}

		nonNilDocCounter++
	}

	assert.NotEqual(0, nonNilDocCounter)
	assert.Len(nonConformedVBList, 0)
}

func TestGlobalTimestamp_GetValue(t *testing.T) {
	tests := []struct {
		name string
		g    GlobalTimestamp
		want interface{}
	}{
		{
			name: "globalGetValueTest",
			g: GlobalTimestamp{
				0: &GlobalVBTimestamp{
					TargetVBTimestamp: TargetVBTimestamp{
						Target_vb_opaque: &TargetVBUuid{1},
						Target_Seqno:     2,
					},
				},
				1: &GlobalVBTimestamp{
					TargetVBTimestamp: TargetVBTimestamp{
						Target_vb_opaque: &TargetVBUuid{5},
						Target_Seqno:     3,
					},
				},
			},
			want: GlobalTimestamp{
				0: &GlobalVBTimestamp{
					TargetVBTimestamp: TargetVBTimestamp{
						Target_vb_opaque: &TargetVBUuid{1},
						Target_Seqno:     2,
					},
				},
				1: &GlobalVBTimestamp{
					TargetVBTimestamp: TargetVBTimestamp{
						Target_vb_opaque: &TargetVBUuid{5},
						Target_Seqno:     3,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantMap, ok := tt.want.(GlobalTimestamp)
			assert.True(t, ok)

			for k, v := range wantMap {
				gts, ok := tt.g.GetValue().(GlobalTimestamp)
				assert.True(t, ok)
				assert.True(t, gts[k].SameAs(v))
			}
		})
	}
}

func TestCheckpointDocMarshallerGlobalCkpt(t *testing.T) {
	fmt.Println("============== Test case start: TestCheckpointDocMarshallerGlobalCkpt =================")
	defer fmt.Println("============== Test case end: TestCheckpointDocMarshallerGlobalCkpt =================")
	assert := assert.New(t)

	vbUuidAndTimestamp := &TargetVBUuidAndTimestamp{
		Target_vb_uuid: "abc",
		Startup_time:   "012",
	}
	newCkptRecord := CheckpointRecord{
		SourceVBTimestamp: SourceVBTimestamp{
			Failover_uuid:                0,
			Seqno:                        1,
			Dcp_snapshot_seqno:           2,
			Dcp_snapshot_end_seqno:       3,
			SourceManifestForDCP:         7,
			SourceManifestForBackfillMgr: 8,
		},
		SourceVBCounters: SourceVBCounters{
			FilteredItemsCnt:  5,
			FilteredFailedCnt: 6,
			CasPoisonCnt:      2,
		},
		GlobalTimestamp: GlobalTimestamp{
			100: &GlobalVBTimestamp{
				TargetVBTimestamp{
					Target_vb_opaque: vbUuidAndTimestamp,
					Target_Seqno:     5,
					TargetManifest:   10,
				},
			},
		},
		GlobalCounters: GlobalTargetCounters{
			1: &TargetPerVBCounters{},
		},
	}
	assert.Nil(newCkptRecord.PopulateShasForGlobalInfo())

	ns1, err := base.NewCollectionNamespaceFromString("s1.col1")
	assert.Nil(err)
	ns2, err := base.NewCollectionNamespaceFromString("s1.col2")
	assert.Nil(err)

	brokenMap := make(CollectionNamespaceMapping)
	brokenMap.AddSingleMapping(&ns1, &ns2)
	brokenMap1Sha, err := brokenMap.Sha256()
	assert.Nil(err)

	brokenMap2 := make(CollectionNamespaceMapping)
	brokenMap2.AddSingleMapping(&ns2, &ns1)
	brokenMap2Sha, err := brokenMap.Sha256()
	assert.Nil(err)

	ckptRecord2 := CheckpointRecord{
		SourceVBTimestamp: SourceVBTimestamp{
			Failover_uuid:                0,
			Seqno:                        1,
			Dcp_snapshot_seqno:           2,
			Dcp_snapshot_end_seqno:       3,
			SourceManifestForDCP:         7,
			SourceManifestForBackfillMgr: 8,
		},
		SourceVBCounters: SourceVBCounters{
			FilteredItemsCnt:  5,
			FilteredFailedCnt: 6,
			CasPoisonCnt:      2,
		},
		GlobalTimestamp: GlobalTimestamp{
			100: &GlobalVBTimestamp{
				TargetVBTimestamp{
					Target_vb_opaque:    vbUuidAndTimestamp,
					Target_Seqno:        5,
					TargetManifest:      10,
					BrokenMappingSha256: fmt.Sprintf("%s", brokenMap1Sha[:]),
					brokenMappings:      brokenMap,
				},
			},
			200: &GlobalVBTimestamp{
				TargetVBTimestamp{
					Target_vb_opaque:    vbUuidAndTimestamp,
					Target_Seqno:        5,
					TargetManifest:      10,
					BrokenMappingSha256: fmt.Sprintf("%s", brokenMap2Sha[:]),
					brokenMappings:      brokenMap2,
				},
			},
		},
		GlobalCounters: GlobalTargetCounters{
			1: &TargetPerVBCounters{},
		},
	}
	assert.Nil(ckptRecord2.PopulateBrokenMappingSha())
	assert.Nil(ckptRecord2.PopulateShasForGlobalInfo())

	ckpt_doc := NewCheckpointsDoc("testInternalId")
	added, _ := ckpt_doc.AddRecord(&newCkptRecord)
	assert.True(added)
	added, _ = ckpt_doc.AddRecord(&ckptRecord2)
	assert.True(added)

	marshalledData, err := json.Marshal(ckpt_doc)
	assert.Nil(err)

	ckptDocCompressed, shaMapCompressed, err := ckpt_doc.SnappyCompress()
	assert.Nil(err)

	// Before checking, need to get sha for global timestamp
	gInfoShaMap := make(ShaToGlobalInfoMap)
	for _, record := range ckpt_doc.Checkpoint_records {
		if record == nil {
			continue
		}
		gts := record.GlobalTimestamp
		gInfoShaMap[record.GlobalTimestampSha256] = &gts
		gCounters := record.GlobalCounters
		gInfoShaMap[record.GlobalCountersSha256] = &gCounters

	}

	var checkDoc CheckpointsDoc
	err = json.Unmarshal(marshalledData, &checkDoc)
	assert.Nil(err)
	for _, record := range checkDoc.Checkpoint_records {
		if record == nil {
			continue
		}
		assert.Nil(record.LoadGlobalInfoMapping(gInfoShaMap))
	}

	assert.Equal(12, len(checkDoc.Checkpoint_records))
	assert.NotNil(checkDoc.Checkpoint_records[1])
	assert.True(checkDoc.Checkpoint_records[1].SameAs(&newCkptRecord))
	assert.Equal(newCkptRecord.SourceManifestForBackfillMgr, checkDoc.Checkpoint_records[1].SourceManifestForBackfillMgr)
	assert.NotNil(checkDoc.Checkpoint_records[0])
	assert.True(checkDoc.Checkpoint_records[0].SameAs(&ckptRecord2))

	var decompressCheck CheckpointsDoc
	brokenMap1, globalInfoShaMap, err := SnappyDecompressShaMap(shaMapCompressed)
	assert.Nil(err)
	assert.Nil(decompressCheck.SnappyDecompress(ckptDocCompressed, brokenMap1, globalInfoShaMap))
	assert.Equal(12, len(decompressCheck.Checkpoint_records))
	assert.NotNil(decompressCheck.Checkpoint_records[1])
	assert.True(decompressCheck.Checkpoint_records[1].SameAs(&newCkptRecord))
	assert.Equal(newCkptRecord.SourceManifestForBackfillMgr, decompressCheck.Checkpoint_records[1].SourceManifestForBackfillMgr)
	assert.NotNil(decompressCheck.Checkpoint_records[0])
	assert.True(decompressCheck.Checkpoint_records[0].SameAs(&ckptRecord2))

	// Test clone
	clonedDoc := ckpt_doc.Clone()
	for i, record := range clonedDoc.Checkpoint_records {
		assert.True(record.SameAs(ckpt_doc.Checkpoint_records[i]))
	}
}

func Test_compareFailoverLogPositionThenSeqnos(t *testing.T) {
	type args struct {
		aRecord *CheckpointSortRecord
		bRecord *CheckpointSortRecord
		source  bool
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 bool
	}{
		{
			name: "Traditional ckpt no failover, compare between seqno, aRecord newer than bRecord",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno: 20,
						},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno: 10,
						},
					},
				},
				source: true,
			},
			want:  true, // aRecord is more recent than bRecord
			want1: true, // validOp
		},
		{
			name: "Traditional ckpt no failover, compare between seqno, aRecord older than bRecord",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno: 10,
						},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno: 20,
						},
					},
				},
				source: true,
			},
			want:  false, // aRecord is not more recent than bRecord
			want1: true,  // valid op
		},
		{
			name: "Traditional ckpt failover, compare between failoverlogs even if seqno is lower, aRecord newer than bRecord",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         10,
							Failover_uuid: 3,
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         20,
							Failover_uuid: 2,
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				source: true,
			},
			want:  true, // aRecord is more recent than bRecord
			want1: true, // valid op
		},
		{
			name: "Traditional ckpt failover, compare between failoverlogs but aRecord does not exist in failover log. bRecord wins",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         10,
							Failover_uuid: 4,
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         20,
							Failover_uuid: 2,
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				source: true,
			},
			want:  false, // aRecord is more recent than bRecord
			want1: true,  // valid op
		},
		{
			name: "Traditional target ckpt no failover, compare between seqno, aRecord newer than bRecord",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						TargetVBTimestamp: TargetVBTimestamp{
							Target_Seqno:     20,
							Target_vb_opaque: &TargetVBUuid{0},
						},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						TargetVBTimestamp: TargetVBTimestamp{
							Target_Seqno:     10,
							Target_vb_opaque: &TargetVBUuid{0},
						},
					},
				},
				source: false,
			},
			want:  true, // aRecord is more recent than bRecord
			want1: true, // valid op
		},
		{
			name: "Traditional target ckpt no failover, compare between seqno, aRecord older than bRecord",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						TargetVBTimestamp: TargetVBTimestamp{
							Target_Seqno:     10,
							Target_vb_opaque: &TargetVBUuid{0},
						},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						TargetVBTimestamp: TargetVBTimestamp{
							Target_Seqno:     20,
							Target_vb_opaque: &TargetVBUuid{0},
						},
					},
				},
				source: false,
			},
			want:  false, // aRecord is not more recent than bRecord
			want1: true,  // valid op
		},
		{
			name: "Traditional target ckpt failover, compare between failoverlogs even if seqno is lower, aRecord newer than bRecord",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						TargetVBTimestamp: TargetVBTimestamp{
							Target_Seqno:     10,
							Target_vb_opaque: &TargetVBUuid{3},
						},
					},
					tgtFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						TargetVBTimestamp: TargetVBTimestamp{
							Target_Seqno:     20,
							Target_vb_opaque: &TargetVBUuid{2},
						},
					},
					tgtFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				source: false,
			},
			want:  true, // aRecord is more recent than bRecord
			want1: true, // valid op
		},
		{
			name: "Traditional target ckpt failover, compare between failoverlogs but aRecord does not exist in failover log. bRecord wins",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						TargetVBTimestamp: TargetVBTimestamp{
							Target_Seqno:     10,
							Target_vb_opaque: &TargetVBUuid{4},
						},
					},
					tgtFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						TargetVBTimestamp: TargetVBTimestamp{
							Target_Seqno:     20,
							Target_vb_opaque: &TargetVBUuid{2},
						},
					},
					tgtFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				source: false,
			},
			want:  false, // aRecord is more recent than bRecord
			want1: true,  // valid op
		},
		{
			name: "Global target ckpt no failover, source seqno same, target global vbucket map mismatch, invalid",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         10,
							Failover_uuid: 3,
						},
						GlobalTimestamp: GlobalTimestamp{
							0: &GlobalVBTimestamp{},
							1: &GlobalVBTimestamp{},
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         10,
							Failover_uuid: 3,
						},
						GlobalTimestamp: GlobalTimestamp{
							1: &GlobalVBTimestamp{},
							2: &GlobalVBTimestamp{},
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				source: false,
			},
			want:  false, // aRecord is more recent than bRecord
			want1: false, // valid op
		},
		{
			name: "Global target ckpt no failover, source seqno different, target global vbucket map mismatch, work off of source to ensure valid",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         20,
							Failover_uuid: 3,
						},
						GlobalTimestamp: GlobalTimestamp{
							0: &GlobalVBTimestamp{},
							1: &GlobalVBTimestamp{},
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         10,
							Failover_uuid: 3,
						},
						GlobalTimestamp: GlobalTimestamp{
							1: &GlobalVBTimestamp{},
							2: &GlobalVBTimestamp{},
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				source: true,
			},
			want:  true, // aRecord is more recent than bRecord
			want1: true, // valid op
		},
		{
			name: "Global target ckpt no failover, source seqno same, target global vbucket map mismatch, considered invalid",
			args: args{
				aRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         10,
							Failover_uuid: 3,
						},
						GlobalTimestamp: GlobalTimestamp{
							0: &GlobalVBTimestamp{},
							1: &GlobalVBTimestamp{},
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				bRecord: &CheckpointSortRecord{
					CheckpointRecord: &CheckpointRecord{
						SourceVBTimestamp: SourceVBTimestamp{
							Seqno:         10,
							Failover_uuid: 3,
						},
						GlobalTimestamp: GlobalTimestamp{
							1: &GlobalVBTimestamp{},
							2: &GlobalVBTimestamp{},
						},
					},
					srcFailoverLog: &mcc.FailoverLog{
						{3, 0},
						{2, 0},
					},
				},
				source: false,
			},
			want:  false, // aRecord is more recent than bRecord
			want1: false, // valid op
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := compareFailoverLogPositionThenSeqnos(tt.args.aRecord, tt.args.bRecord, tt.args.source)
			assert.Equalf(t, tt.want, got, "compareFailoverLogPositionThenSeqnos(%v, %v, %v)", tt.args.aRecord, tt.args.bRecord, tt.args.source)
			assert.Equalf(t, tt.want1, got1, "compareFailoverLogPositionThenSeqnos(%v, %v, %v)", tt.args.aRecord, tt.args.bRecord, tt.args.source)
		})
	}
}
