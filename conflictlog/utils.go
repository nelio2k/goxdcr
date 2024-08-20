package conflictlog

import (
	"fmt"

	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
)

type LoggerGetter func() Logger

// returns a logger only if non-null rules are parsed without any errors.
func LoggerForRules(conflictLoggingMap base.ConflictLoggingMappingInput, replId string, logger_ctx *log.LoggerContext, logger *log.CommonLogger) (Logger, error) {
	if conflictLoggingMap == nil {
		return nil, fmt.Errorf("nil conflictLoggingMap")
	}

	conflictLoggingNotDisabled := !conflictLoggingMap.Disabled()
	if !conflictLoggingNotDisabled {
		logger.Infof("Conflict logger will be off for pipeline=%s, with input=%v", replId, conflictLoggingMap)
		return nil, nil
	}

	var rules *base.ConflictLogRules
	var err error
	rules, err = base.ParseConflictLogRules(conflictLoggingMap)
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

// Inserts "_xdcr_conflict": true to the input byte slice.
// If any error occurs, the original body is returned.
// Otherwise, returns new body and new datatype after xattr is successfully added.
func InsertConflictXattrToBody(body []byte, datatype uint8) ([]byte, uint8, error) {
	newbodyLen := len(body) + MaxBodyIncrease
	// TODO - Use datapool.
	newbody := make([]byte, newbodyLen)

	xattrComposer := base.NewXattrComposer(newbody)

	if base.HasXattr(datatype) {
		// insert the already existing xattrs
		it, err := base.NewXattrIterator(body)
		if err != nil {
			return body, datatype, err
		}

		for it.HasNext() {
			key, val, err := it.Next()
			if err != nil {
				return body, datatype, err
			}
			err = xattrComposer.WriteKV(key, val)
			if err != nil {
				return body, datatype, err
			}
		}
	}

	err := xattrComposer.WriteKV(base.ConflictLoggingXattrKeyBytes, base.ConflictLoggingXattrValBytes)
	if err != nil {
		return body, datatype, err
	}

	docWithoutXattr := base.FindDocBodyWithoutXattr(body, datatype)
	out, atLeastOneXattr := xattrComposer.FinishAndAppendDocValue(docWithoutXattr, nil, nil)

	if atLeastOneXattr {
		datatype = datatype | base.PROTOCOL_BINARY_DATATYPE_XATTR
	} else {
		// odd - shouldn't happen.
		datatype = datatype & ^(base.PROTOCOL_BINARY_DATATYPE_XATTR)
	}

	body = nil // no use of this body anymore, set to nil to help GC quicker.

	return out, datatype, err
}
