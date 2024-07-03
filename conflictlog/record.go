package conflictlog

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
	Id         string
	BucketUUID string
	IsDeleted  bool
	Collection string
	Scope      string
}

// ConflictRecord has the all the details of the detected conflict
// which are needed to be persisted
type ConflictRecord struct {
	ReplicationId string
	Source        DocInfo
	Target        DocInfo
}

func (r *ConflictRecord) Scope() string {
	return r.Source.Scope
}

func (r *ConflictRecord) Collection() string {
	return r.Source.Collection
}
