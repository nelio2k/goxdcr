package conflictlog

import (
	"context"
	"sync"

	"github.com/couchbase/goxdcr/log"
)

var _ Logger = (*loggerImpl)(nil)

// loggerImpl implements the Logger interface
type loggerImpl struct {
	// replId is the unique replication ID
	replId string

	// mapper maps the conflict to the target conflict bucket
	mapper Mapper

	// rules describe the mapping to the target conflict bucket
	rules     Rules
	rulesLock sync.RWMutex

	logger *log.CommonLogger
}

func WithRules(r Rules) LoggerOpt {
	return func(o *LoggerOptions) {
		o.rules = r
	}
}

func WithMapper(m Mapper) LoggerOpt {
	return func(o *LoggerOptions) {
		o.mapper = m
	}
}

func newLoggerImpl(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (m *loggerImpl, err error) {
	options := &LoggerOptions{}
	for _, opt := range opts {
		opt(options)
	}

	m = &loggerImpl{
		logger:    logger,
		replId:    replId,
		rules:     options.rules,
		rulesLock: sync.RWMutex{},
		mapper:    options.mapper,
	}

	return
}

func (l *loggerImpl) Log(ctx context.Context, c *ConflictRecord) (err error) {
	return
}

func (l *loggerImpl) UpdateRules(r Rules) (err error) {
	err = r.Validate()
	if err != nil {
		return
	}

	l.rulesLock.Lock()
	l.rules = r
	l.rulesLock.Unlock()

	return
}
