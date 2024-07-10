package conflictlog

import (
	"github.com/couchbase/goxdcr/log"
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
func InitManager(loggerCtx *log.LoggerContext, memdAddrGetter MemcachedAddrGetter) Manager {
	logger := log.NewLogger(ConflictManagerLoggerName, loggerCtx)

	logger.Info("intializing conflict manager")

	impl := &managerImpl{
		logger:         logger,
		memdAddrGetter: memdAddrGetter,
	}

	logger.Info("creating conflict manager writer pool")
	impl.writerPool = newWriterPool(impl.newWriter)

	return impl
}

// managerImpl implements conflict manager
type managerImpl struct {
	logger         *log.CommonLogger
	memdAddrGetter MemcachedAddrGetter
	writerPool     *writerPool
}

func (m *managerImpl) NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error) {

	l, err = newLoggerImpl(logger, replId, m.writerPool, opts...)
	if err != nil {
		return
	}

	return
}

func (m *managerImpl) newWriter(bucketName string) (w Writer, err error) {
	m.logger.Infof("creating new conflict writer bucket=%s", bucketName)
	return newGocbConn(m.logger, m.memdAddrGetter, bucketName)
}
