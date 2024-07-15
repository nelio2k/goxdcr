package conflictlog

import "github.com/couchbase/goxdcr/base"

// Logger interface allows logging of conflicts in an abstracted manner
type Logger interface {
	// Log writes the conflict to the conflict buccket
	Log(c *ConflictRecord) (base.ConflictLoggerHandle, error)

	// UpdateWorkerCount changes the underlying log worker count
	UpdateWorkerCount(count int)

	// UpdateRules allow updates to the the rules which map
	// the conflict to the target conflict bucket
	UpdateRules(*Rules) error

	Close() error
}

// LoggerOptions defines optional args for a logger implementation
type LoggerOptions struct {
	rules       *Rules
	mapper      Mapper
	logQueueCap int
	workerCount int
}

func (o *LoggerOptions) SetRules(rules *Rules) {
	o.rules = rules
}

func (o *LoggerOptions) SetMapper(mapper Mapper) {
	o.mapper = mapper
}

func (o *LoggerOptions) SetLogQueueCap(cap int) {
	o.logQueueCap = cap
}

type LoggerOpt func(o *LoggerOptions)
