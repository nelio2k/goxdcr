package conflictlog

import (
	"errors"
	"fmt"

	"github.com/couchbase/goxdcr/base"
)

var (
	ErrEmptyConflictLoggingMap      error = errors.New("nil or empty conflict logging map not allowed")
	ErrInvalidBucket                error = errors.New("conflict logging bucket not provided or wrong format")
	ErrInvalidCollection            error = errors.New("conflict logging collection not provided or wrong format")
	ErrInvalidCollectionFormat      error = errors.New("conflict logging collection should be in format [scope].[collection]")
	ErrInvalidLoggingRulesFormat    error = errors.New("conflict logging loggingRules should be a json object")
	ErrInvalidLoggingRulesValFormat error = errors.New("conflict logging loggingRules' values should be a json object")
	ErrEmptyBucketEmpty             error = errors.New("conflict logging bucket should not be empty")
	ErrEmptyCollectionEmpty         error = errors.New("conflict logging collection should not be empty")
)

// Rules captures the logging rules for a replication
type Rules struct {
	// Target is the default or fallback target conflict bucket
	Target Target

	// Mapping describes the how the conflicts from a source scope
	// & collection is logged to the Target
	// Empty map implies that all conflicts will be logged
	Mapping map[base.CollectionNamespace]Target
}

func parseString(o interface{}) (ok bool, val string) {
	if o == nil {
		return
	}
	val, ok = o.(string)
	return
}

func parseTarget(m map[string]interface{}) (t Target, err error) {
	if m == nil {
		return
	}

	bucketObj, ok := m[base.CLBucketKey]
	if ok {
		ok, s := parseString(bucketObj)
		if ok {
			t.Bucket = s
		}
	}

	collectionObj, ok := m[base.CLCollectionKey]
	if ok {
		ok, s := parseString(collectionObj)
		if ok {
			t.NS, err = base.NewCollectionNamespaceFromString(s)
		} else {
			err = ErrInvalidCollectionValueType
			return
		}
	}

	return
}

// ParseRules parses map[string]interface{} object into rules.
func ParseRules(j base.ConflictLoggingMappingInput) (rules *Rules, err error) {
	fallbackTarget, err := parseTarget(j)
	if err != nil {
		return
	}

	if !fallbackTarget.IsComplete() {
		err = ErrIncompleteTarget
		return
	}

	rules = &Rules{
		Target:  fallbackTarget,
		Mapping: map[base.CollectionNamespace]Target{},
	}

	loggingRulesObj, ok := j[base.CLLoggingRulesKey]
	if !ok || loggingRulesObj == nil {
		return
	}

	loggingRulesMap, ok := loggingRulesObj.(map[string]interface{})
	if !ok {
		rules = nil
		err = ErrInvalidLoggingRulesType
		return
	}

	for collectionStr, targetObj := range loggingRulesMap {
		if collectionStr == "" {
			rules = nil
			err = ErrInvalidCollection
			return
		}

		source, err := base.NewOptionalCollectionNamespaceFromString(collectionStr)
		if err != nil {
			return nil, err
		}

		target := fallbackTarget

		if targetObj != nil {
			targetMap, ok := targetObj.(map[string]interface{})
			if !ok {
				rules = nil
				err = ErrInvalidTargetType
				return nil, err
			}

			if len(targetMap) > 0 {
				target, err = parseTarget(targetMap)
				if err != nil {
					return nil, err
				}

				if !target.IsComplete() {
					return nil, ErrIncompleteTarget
				}
			}
		}

		rules.Mapping[source] = target
	}

	return
}

func (r *Rules) Validate() (err error) {
	if !r.Target.IsComplete() {
		err = ErrIncompleteTarget
		return
	}

	for m, t := range r.Mapping {
		if m.ScopeName == "" {
			err = ErrEmptyBucketEmpty
			return
		}

		if !t.IsComplete() {
			err = ErrIncompleteTarget
			return
		}
	}

	return
}

func (r *Rules) SameAs(other *Rules) (same bool) {
	if r == nil || other == nil {
		return r == nil && other == nil
	}

	same = r.Target.SameAs(other.Target)
	if !same {
		return
	}

	if r.Mapping == nil || other.Mapping == nil {
		return r.Mapping == nil && other.Mapping == nil
	}

	same = len(r.Mapping) == len(other.Mapping)
	if !same {
		return
	}

	var otherSource Target
	for mapping, target := range r.Mapping {
		otherSource, same = other.Mapping[mapping]
		if !same {
			return
		}

		same = otherSource.SameAs(target)
		if !same {
			return
		}
	}

	same = true
	return
}

func (r *Rules) String() string {
	if r == nil {
		return "<nil>"
	}

	loggingRules := ""
	for source, target := range r.Mapping {
		loggingRules += fmt.Sprintf("%s -> %s,", source, target)
	}

	loggingRules = "{" + loggingRules + "}"

	return fmt.Sprintf("Target: %s, loggingRules: %s", r.Target, loggingRules)
}

// Mapping is the source scope and collection
type Mapping struct {
	// Scope is source bucket's scope
	Scope string
	// Collection is source bucket's collection
	Collection string
}

func (m Mapping) String() string {
	return fmt.Sprintf("%v.%v", m.Scope, m.Collection)
}

func (m Mapping) Validate() (err error) {
	if m.Scope == "" {
		err = ErrEmptyScope
		return
	}
	// SUMUKH TODO - can be empty, deal with it
	if m.Collection == "" {
		err = ErrEmptyCollectionEmpty
		return
	}
	return
}

func (m *Mapping) SameAs(other *Mapping) bool {
	if m == nil || other == nil {
		return m == nil && other == nil
	}

	return m.Scope == other.Scope &&
		m.Collection == other.Collection
}

// Target describes the target bucket, scope and collection where
// the conflicts will be logged. There are few terms:
// Complete => The triplet (Bucket, Scope, Collection) are populated
// Empty => All components of the triplet are empty
type Target struct {
	// Bucket is the conflict bucket
	Bucket string `json:"bucket"`

	// NS is namespace which defines scope and collection
	NS base.CollectionNamespace `json:"ns"`
}

func (t Target) String() string {
	return fmt.Sprintf("%v.%v.%v", t.Bucket, t.NS.ScopeName, t.NS.CollectionName)
}

func NewTarget(bucket, scope, collection string) Target {
	return Target{
		Bucket: bucket,
		NS: base.CollectionNamespace{
			ScopeName:      scope,
			CollectionName: collection,
		},
	}
}

func (t Target) IsEmpty() bool {
	return t.Bucket == "" && t.NS.IsEmpty()
}

// IsComplete implies that all components of the triplet (bucket, scope & collection)
// are populated
func (t Target) IsComplete() bool {
	return t.Bucket != "" && t.NS.ScopeName != "" && t.NS.CollectionName != ""
}

func (t Target) SameAs(other Target) bool {
	return t.Bucket == other.Bucket && t.NS.IsSameAs(other.NS)
}
