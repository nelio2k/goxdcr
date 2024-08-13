package conflictlog

import (
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
)

// Mapper evaluates and routes the conflict to the right target bucket
// It also stores the rules against which the mapping will happen
type Mapper interface {
	// Map evaluates and routes the conflict to the right target bucket
	Map(*Rules, Conflict) (base.ConflictLoggingTarget, error)
}

type FixedMapper struct {
	target base.ConflictLoggingTarget
	logger *log.CommonLogger
}

func NewFixedMapper(logger *log.CommonLogger, target base.ConflictLoggingTarget) *FixedMapper {
	return &FixedMapper{logger: logger, target: target}
}

func (m *FixedMapper) Map(_ *Rules, c Conflict) (base.ConflictLoggingTarget, error) {
	return m.target, nil
}
