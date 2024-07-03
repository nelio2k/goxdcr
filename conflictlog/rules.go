package conflictlog

import "fmt"

// Rules captures the logging rules for a replication
type Rules struct {
	// Target is the default or fallback target conflict bucket
	Target Target

	// Mapping describes the how the conflicts from a source scope
	// & collection is logged to the Target
	// Empty map implies that all conflicts will be logged
	Mapping map[Mapping]Target
}

func (r *Rules) Validate() (err error) {
	err = r.Target.Validate()
	if err != nil {
		return
	}

	for m, t := range r.Mapping {
		err = m.Validate()
		if err != nil {
			return
		}

		err = t.Validate()
		if err != nil {
			return
		}
	}

	return
}

// Mapping is the source scope and collection
type Mapping struct {
	// Scope is source bucket's scope
	Scope string
	// Collection is source bucket's collection
	Collection string
}

func (m Mapping) Validate() (err error) {
	if m.Scope == "" {
		err = fmt.Errorf("empty scope in conflict rules mapping")
		return
	}
	return
}

// Target describes the target bucket, scope and collection where
// the conflicts will be logged
type Target struct {
	// Bucket is the conflict bucket
	Bucket string
	// Scope is the conflict bucket's scope
	Scope string
	// Collection is the conflict bucket's collection
	Collection string
}

func (t Target) Validate() (err error) {
	if t.Bucket == "" {
		err = fmt.Errorf("empty conflict bucket in conflict rules mapping")
		return
	}
	return
}
