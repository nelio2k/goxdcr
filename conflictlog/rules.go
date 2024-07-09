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
	ErrEmptyScopeEmpty              error = errors.New("conflict logging scope should not be empty")
	ErrEmptyCollectionEmpty         error = errors.New("conflict logging collection should not be empty")
)

// Rules captures the logging rules for a replication
type Rules struct {
	// Target is the default or fallback target conflict bucket
	Target Target

	// Mapping describes the how the conflicts from a source scope
	// & collection is logged to the Target
	// Empty map implies that all conflicts will be logged
	Mapping map[Mapping]Target
}

// given user input of conflict logging mapping, the function creates a new and validated "Rules" object.
// ValidateAndConvertJsonMapToConflictLoggingMapping already would have done type validations.
// base.ParseOneConflictLoggingRule would have already been called in ValidateAndConvertJsonMapToConflictLoggingMapping too.
// So at this point, we shouldnt not get any error.
// jsonMapping should not be nil or {}. Such values needs to be handled explicitly by the caller.
func NewRules(jsonMapping base.ConflictLoggingMappingInput) (rules Rules, err error) {
	if len(jsonMapping) == 0 {
		// Nil is not a valid input
		// {} is a valid input, but represent conflict logging is turned off - no rules can be generated
		err = ErrEmptyConflictLoggingMap
		return
	}

	// parse common target
	targetBucketName, targetScopeName, targetCollectionName, err := base.ParseOneConflictLoggingRule(jsonMapping)
	if err != nil {
		return
	}

	rules.Target = NewTarget(
		targetBucketName,
		targetScopeName,
		targetCollectionName,
	)
	rules.Mapping = make(map[Mapping]Target)

	// parse special "logging rules" if present
	loggingRulesIn, ok := jsonMapping[base.CLLoggingRulesKey]
	if ok {
		loggingRules, ok := loggingRulesIn.(map[string]interface{})
		if !ok {
			err = ErrInvalidLoggingRulesFormat
			return
		}

		// SUMUKH TODO - complex rules.
		// Right now only [scope].[collection] -> {bucket: ..., collection: ...} is parsed
		for source, target := range loggingRules {
			srcScopeName, srcCollectionName := base.SeparateScopeCollection(source)

			targetMap, ok := target.(map[string]interface{})
			if !ok {
				return rules, ErrInvalidLoggingRulesValFormat
			}

			targetBucketName, targetScopeName, targetCollectionName, err = base.ParseOneConflictLoggingRule(targetMap)
			if err != nil {
				return
			}

			mapping := Mapping{
				Scope:      srcScopeName,
				Collection: srcCollectionName,
			}

			rules.Mapping[mapping] = NewTarget(
				targetBucketName,
				targetScopeName,
				targetCollectionName,
			)
		}
	}

	// validate the rules
	err = rules.Validate()
	if err != nil {
		return rules, err
	}

	return rules, nil
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

func (r *Rules) SameAs(other *Rules) (same bool) {
	if r == nil || other == nil {
		return r == nil && other == nil
	}

	same = r.Target.SameAs(&other.Target)
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

		same = otherSource.SameAs(&target)
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
		err = ErrEmptyScopeEmpty
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
// the conflicts will be logged
type Target struct {
	// Bucket is the conflict bucket
	Bucket string
	// Scope is the conflict bucket's scope
	Scope string
	// Collection is the conflict bucket's collection
	Collection string
}

func (t Target) String() string {
	return fmt.Sprintf("%v.%v.%v", t.Bucket, t.Scope, t.Collection)
}

func NewTarget(bucket, scope, collection string) Target {
	return Target{
		Bucket:     bucket,
		Scope:      scope,
		Collection: collection,
	}
}

func (t Target) Validate() (err error) {
	if t.Bucket == "" {
		err = ErrEmptyBucketEmpty
		return
	}
	if t.Collection == "" {
		err = ErrEmptyCollectionEmpty
		return
	}
	if t.Scope == "" {
		err = ErrEmptyScopeEmpty
		return
	}
	return
}

func (t *Target) SameAs(other *Target) bool {
	if t == nil || other == nil {
		return t == nil && other == nil
	}

	return t.Bucket == other.Bucket &&
		t.Collection == other.Collection &&
		t.Scope == other.Scope
}
