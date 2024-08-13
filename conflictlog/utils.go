package conflictlog

import (
	"fmt"

	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
)

type LoggerGetter func() Logger

// returns a logger only if non-null rules are parsed without any errors.
func LoggerForRules(conflictLoggingMap base.ConflictLoggingMappingInput, replId string, logger_ctx *log.LoggerContext, logger *log.CommonLogger) (Logger, error) {
	conflictLoggingEnabled := len(conflictLoggingMap) > 0 // {} is disabled.
	if !conflictLoggingEnabled {
		logger.Infof("Conflict logger will be off for pipeline=%s, with input=%v", replId, conflictLoggingMap)
		return nil, fmt.Errorf("conflict logging disabled with input %v", conflictLoggingMap)
	}

	var rules *Rules
	var err error
	rules, err = ParseRules(conflictLoggingMap)
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

	fileLogger := log.NewLogger(ConflictLoggerName, logger_ctx)
	conflictLogger, err = clm.NewLogger(
		fileLogger,
		replId,
		WithMapper(NewConflictMapper(logger)),
		WithCapacity(1000), // SUMUKH TODO - make the default size configurable.
		WithRules(rules),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting a new conflict logger for %v. err=%v", conflictLoggingMap, err)
	}

	logger.Infof("Conflict logger will be on for pipeline=%s, with rules=%s for input=%v", replId, rules, conflictLoggingMap)

	return conflictLogger, nil
}