package conflictlog

import "errors"

var (
	ErrManagerNotInitialized      error = errors.New("conflict manager not initialized")
	ErrWriterTimeout              error = errors.New("conflict writer timed out")
	ErrWriterClosed               error = errors.New("conflict writer closed")
	ErrQueueFull                  error = errors.New("conflict log is full")
	ErrLoggerClosed               error = errors.New("conflict logger is closed")
	ErrLogWaitAborted             error = errors.New("conflict log handle received abort")
	ErrEmptyRules                 error = errors.New("empty conflict rules")
	ErrUnknownCollection          error = errors.New("unknown collection")
	ErrClosedConnPool             error = errors.New("use of closed connection pool")
	ErrNotMyBucket                error = errors.New("not my bucket")
	ErrEmptyTarget                error = errors.New("empty target")
	ErrEmptyScope                 error = errors.New("conflict logging scope should not be empty")
	ErrIncompleteTarget           error = errors.New("one or more of target bucket,scope or collection is empty")
	ErrInvalidLoggingRulesType    error = errors.New("invalid logging rules type")
	ErrInvalidTargetType          error = errors.New("invalid target type")
	ErrInvalidCollectionValueType error = errors.New("collection value is not string")
)
