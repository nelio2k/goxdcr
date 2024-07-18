package conflictlog

import (
	"time"

	"github.com/couchbase/goxdcr/base"
)

// Logger interface allows logging of conflicts in an abstracted manner
type Logger interface {
	// Log writes the conflict to the conflict buccket
	Log(c *ConflictRecord) (base.ConflictLoggerHandle, error)

	// UpdateWorkerCount changes the underlying log worker count
	UpdateWorkerCount(count int)

	// UpdateRules allow updates to the the rules which map
	// the conflict to the target conflict bucket
	UpdateRules(*Rules) error

	// Closes the logger. Hence forth the logger will error out
	Close() error
}

// LoggerOptions defines optional args for a logger implementation
type LoggerOptions struct {
	rules                *Rules
	mapper               Mapper
	logQueueCap          int
	workerCount          int
	networkRetryCount    int
	networkRetryInterval time.Duration
}

func WithRules(r *Rules) LoggerOpt {
	return func(o *LoggerOptions) {
		o.rules = r
	}
}

func WithMapper(m Mapper) LoggerOpt {
	return func(o *LoggerOptions) {
		o.mapper = m
	}
}

func WithCapacity(cap int) LoggerOpt {
	return func(o *LoggerOptions) {
		o.logQueueCap = cap
	}
}

func WithWorkerCount(val int) LoggerOpt {
	return func(o *LoggerOptions) {
		o.workerCount = val
	}
}

func WithNetworkRetryCount(val int) LoggerOpt {
	return func(o *LoggerOptions) {
		o.networkRetryCount = val
	}
}

func WithNetworkRetryInterval(val time.Duration) LoggerOpt {
	return func(o *LoggerOptions) {
		o.networkRetryInterval = val
	}
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
