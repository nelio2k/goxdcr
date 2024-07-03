package conflictlog

import (
	"context"
)

// Logger interface allows logging of conflicts in an abstracted manner
type Logger interface {
	// Log writes the conflict to the conflict buccket
	Log(ctx context.Context, c *ConflictRecord) error

	// UpdateRules allow updates to the the rules which map
	// the conflict to the target conflict bucket
	UpdateRules(Rules) error
}

// LoggerOptions defines optional args for a logger implementation
type LoggerOptions struct {
	rules  Rules
	mapper Mapper
}

type LoggerOpt func(o *LoggerOptions)
