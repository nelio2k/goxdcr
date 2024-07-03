package conflictlog

import "github.com/couchbase/goxdcr/log"

type Manager struct {
}

func (m *Manager) NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error) {
	l, err = newLoggerImpl(logger, replId, opts...)
	if err != nil {
		return
	}

	return
}
