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

	// Closes the logger. Hence forth the logger will error out
	Close() error
}

// Handle is returned for every logging request.
// The handle allows it caller to wait on the logging to complete (or error out)
type Handle interface {
	// Wait allows caller to complete the logging
	// The finch is the caller's finch. If the caller
	// wants to exit early then the Wait will unblock as well
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
