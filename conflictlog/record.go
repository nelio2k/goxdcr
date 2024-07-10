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
	Id         string `json:"id"`
	BucketUUID string `json:"bucketUUID"`
	IsDeleted  bool   `json:"isDeleted"`
	Collection string `json:"collection"`
	Scope      string `json:"scope"`

	// Body of the doc. Note: it will not be serialized to json
	Body []byte
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
