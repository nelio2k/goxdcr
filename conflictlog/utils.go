package conflictlog

import (
	"fmt"

	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
)

type LoggerGetter func() Logger

// returns a logger only if non-null rules are parsed without any errors.
func LoggerForRules(conflictLoggingMap base.ConflictLoggingMappingInput, replId string, logger_ctx *log.LoggerContext) (Logger, error) {
	conflictLoggingEnabled := len(conflictLoggingMap) > 0 // {} is disabled.
	if !conflictLoggingEnabled {
		return nil, fmt.Errorf("conflict logging disabled with input %v", conflictLoggingMap)
	}

	var rules *Rules
	var err error
	rules, err = NewRules(conflictLoggingMap)
	if err != nil {
		return nil, fmt.Errorf("error converting %v to rules, err=%v", conflictLoggingMap, err)
	}

	if rules == nil {
		return nil, fmt.Errorf("%v maps to nil rules", conflictLoggingMap)
	}

	var clm Manager
	var conflictLogger Logger
	clm, err = GetManager()
	if err != nil {
		return nil, fmt.Errorf("error getting conflict logging manager. err=%v", err)
	}

	logger := log.NewLogger(ConflictLoggerName, logger_ctx)
	conflictLogger, err = clm.NewLogger(
		logger,
		replId,
		WithMapper(NewConflictMapper(logger)),
		WithCapacity(1000), // SUMUKH TODO - make the default size configurable.
		WithRules(rules),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting a new conflict logger for %v. err=%v", conflictLoggingMap, err)
	}

	return conflictLogger, nil
}
