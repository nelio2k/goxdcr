package conflictlog

import (
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/service_def"
)

const (
	ConflictManagerLoggerName = "conflictMgr"
)

var manager Manager

var _ Manager = (*managerImpl)(nil)

// Manager defines behaviour for conflict manager
type Manager interface {
	NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error)
}

// GetManager returns the global conflict manager
func GetManager() (Manager, error) {
	if manager == nil {
		return nil, ErrManagerNotInitialized
	}

	return manager, nil
}

// InitManager intializes global conflict manager
func InitManager(loggerCtx *log.LoggerContext, xdcrTopoSvc service_def.XDCRCompTopologySvc) Manager {
	manager = &managerImpl{
		logger:      log.NewLogger(ConflictManagerLoggerName, loggerCtx),
		xdcrTopoSvc: xdcrTopoSvc,
	}
	return manager
}

// managerImpl implements conflict manager
type managerImpl struct {
	logger      *log.CommonLogger
	xdcrTopoSvc service_def.XDCRCompTopologySvc
}

func (m *managerImpl) NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error) {
	l, err = newLoggerImpl(logger, replId, opts...)
	if err != nil {
		return
	}

	return
}
