package metadata

import ()

type BackfillReplicationSpec struct {
	//id of the replication
	Id_ string `json:"id"`

	// Stored backfill request
	BackfillTasks *BackfillPersistInfo `json:"BackfillTasks"`

	// Soft link to corresponding actual ReplicationSpec - not stored as they should be pulled dynamically
	// Should not be nil
	ReplicationSpec *ReplicationSpecification
}

func NewBackfillReplicationSpec(id string,
	tasks *BackfillPersistInfo) *BackfillReplicationSpec {
	return &BackfillReplicationSpec{
		Id_:           id,
		BackfillTasks: tasks,
	}
}

func (s *BackfillReplicationSpec) Id() string {
	return s.Id_
}
func (s *BackfillReplicationSpec) InternalId() string {
	return s.ReplicationSpec.InternalId()
}
func (s *BackfillReplicationSpec) SourceBucketName() string {
	return s.ReplicationSpec.SourceBucketName()
}
func (s *BackfillReplicationSpec) SourceBucketUUID() string {
	return s.ReplicationSpec.SourceBucketUUID()
}
func (s *BackfillReplicationSpec) TargetClusterUUID() string {
	return s.ReplicationSpec.TargetClusterUUID()
}
func (s *BackfillReplicationSpec) TargetBucketName() string {
	return s.ReplicationSpec.TargetBucketName()
}
func (s *BackfillReplicationSpec) TargetBucketUUID() string {
	return s.ReplicationSpec.TargetBucketUUID()
}
func (s *BackfillReplicationSpec) Settings() *ReplicationSettings {
	return s.ReplicationSpec.Settings()
}

func (b *BackfillReplicationSpec) SameSpec(other *BackfillReplicationSpec) bool {
	if b == nil && other != nil || b != nil && other == nil {
		return false
	} else if b == nil && other == nil {
		return true
	}

	return b.Id() == other.Id() && b.BackfillTasks.Same(other.BackfillTasks) &&
		b.ReplicationSpec.SameSpec(other.ReplicationSpec)
}

func (b *BackfillReplicationSpec) Clone() *BackfillReplicationSpec {
	if b == nil {
		return nil
	}
	spec := &BackfillReplicationSpec{
		Id_:           b.Id(),
		BackfillTasks: b.BackfillTasks.Clone(),
	}
	if b.ReplicationSpec != nil {
		spec.ReplicationSpec = b.ReplicationSpec.Clone()
	}
	return spec
}

func (b *BackfillReplicationSpec) Redact() *BackfillReplicationSpec {
	if b != nil && b.ReplicationSpec != nil {
		b.ReplicationSpec.Settings_.Redact()
	}
	return b
}

func (b *BackfillReplicationSpec) SameSpecGeneric(other GenericSpecification) bool {
	return b.SameSpec(other.(*BackfillReplicationSpec))
}

func (b *BackfillReplicationSpec) CloneGeneric() GenericSpecification {
	return b.Clone()
}

func (b *BackfillReplicationSpec) RedactGeneric() GenericSpecification {
	return b.Redact()
}
