package conflictlog

import (
	"fmt"

	"github.com/couchbase/goxdcr/base"
)

// Rules captures the logging rules for a replication
type Rules struct {
	// Target is the default or fallback target conflict bucket
	Target base.ConflictLoggingTarget

	// Mapping describes the how the conflicts from a source scope
	// & collection is logged to the Target
	// Empty map implies that all conflicts will be logged
	Mapping map[base.CollectionNamespace]base.ConflictLoggingTarget
}

// ParseRules parses map[string]interface{} object into rules.
// should be in sync with base.ValidateConflictLoggingMapValues
// j should not be empty or nil
func ParseRules(j base.ConflictLoggingMappingInput) (rules *Rules, err error) {
	if j == nil {
		err = ErrNilMapping
		return
	}

	fallbackTarget, err := base.ParseConflictLoggingTarget(j)
	if err != nil {
		return
	}

	if !fallbackTarget.IsComplete() {
		err = ErrIncompleteTarget
		return
	}

	rules = &Rules{
		Target:  fallbackTarget,
		Mapping: map[base.CollectionNamespace]base.ConflictLoggingTarget{},
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
				target, err = base.ParseConflictLoggingTarget(targetMap)
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

	var otherSource base.ConflictLoggingTarget
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
