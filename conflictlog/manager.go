package conflictlog

import (
	"io"

	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/utils"
)

const (
	ConflictManagerLoggerName = "conflictMgr"
)

var manager Manager

var _ Manager = (*managerImpl)(nil)

// Manager defines behaviour for conflict manager
type Manager interface {
	NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error)
	ConnPool() ConnPool
}

type MemcachedAddrGetter interface {
	MyMemcachedAddr() (string, error)
}

// GetManager returns the global conflict manager
func GetManager() (Manager, error) {
	if manager == nil {
		return nil, ErrManagerNotInitialized
	}

	return manager, nil
}

// InitManager intializes global conflict manager
func InitManager(loggerCtx *log.LoggerContext, utils utils.UtilsIface, memdAddrGetter MemcachedAddrGetter) {
	logger := log.NewLogger(ConflictManagerLoggerName, loggerCtx)

	logger.Info("intializing conflict manager")

	impl := &managerImpl{
		logger:         logger,
		memdAddrGetter: memdAddrGetter,
		utils:          utils,
	}

	logger.Info("creating conflict manager writer pool")
	impl.connPool = newConnPool(logger, impl.newConn)

	manager = impl
}

// managerImpl implements conflict manager
type managerImpl struct {
	logger         *log.CommonLogger
	memdAddrGetter MemcachedAddrGetter
	utils          utils.UtilsIface
	connPool       *connPool
}

func (m *managerImpl) NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error) {

	l, err = newLoggerImpl(logger, replId, m.utils, m.connPool, opts...)
	if err != nil {
		return
	}

	return
}

func (m *managerImpl) ConnPool() ConnPool {
	return m.connPool
}

func (m *managerImpl) newConn(bucketName string) (w io.Closer, err error) {
	m.logger.Infof("creating new conflict writer bucket=%s", bucketName)
	return newGocbConn(m.logger, m.memdAddrGetter, bucketName)
}
