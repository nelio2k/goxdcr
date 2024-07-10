package conflictlog

// Logger interface allows logging of conflicts in an abstracted manner
type Logger interface {
	// Log writes the conflict to the conflict buccket
	Log(c *ConflictRecord) (Handle, error)

	// UpdateWorkerCount changes the underlying log worker count
	UpdateWorkerCount(count int)

	// UpdateRules allow updates to the the rules which map
	// the conflict to the target conflict bucket
	UpdateRules(*Rules) error
	Close() error
}

type Handle interface {
	Wait(finch chan bool) error
}

// LoggerOptions defines optional args for a logger implementation
type LoggerOptions struct {
	rules       *Rules
	mapper      Mapper
	logQueueCap int
	workerCount int
}

type LoggerOpt func(o *LoggerOptions)
