package conflictlog

import (
	"github.com/couchbase/goxdcr/log"
)

// ConflictMapper impements Mapper interface.
type conflictMapper struct {
	logger *log.CommonLogger
}

func NewConflictMapper(logger *log.CommonLogger) *conflictMapper {
	return &conflictMapper{logger: logger}
}

// returns the "target" to which the conflict record needs to be routed.
func (m *conflictMapper) Map(rules *Rules, c Conflict) (target Target, err error) {
	if rules == nil {
		err = ErrEmptyRules
		return
	}

	// If there are no special "loggingRules", rules.Target is the return target.
	target = rules.Target

	// Check for special "loggingRules" if any
	if rules.Mapping == nil {
		// consider as no loggingRules.
		return
	}

	mapping := Mapping{
		Scope:      c.Scope(),
		Collection: c.Collection(),
	}

	// SUMUKH TODO - complex rules mapping.
	targetOverride, ok := rules.Mapping[mapping]
	if ok {
		target = targetOverride
	}

	return
}
