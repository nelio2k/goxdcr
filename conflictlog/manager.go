package conflictlog

import (
	"io"

	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/utils"
)

var manager Manager

var _ Manager = (*managerImpl)(nil)

// Manager defines behaviour for conflict manager
type Manager interface {
	NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error)
	ConnPool() ConnPool
	// [TEMP]: SetConnType exists only for perf test
	SetConnType(connType string) error
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
		manifestCache:  newManifestCache(),
		connType:       "gocbcore",
	}

	logger.Info("creating conflict manager writer pool")
	impl.setConnPool()

	manager = impl
}

// managerImpl implements conflict manager
type managerImpl struct {
	logger         *log.CommonLogger
	memdAddrGetter MemcachedAddrGetter
	utils          utils.UtilsIface
	connPool       *connPool
	manifestCache  *ManifestCache
	// [TEMP] connType only exists for perf test
	connType string
}

func (m *managerImpl) NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error) {
	l, err = newLoggerImpl(logger, replId, m.utils, m.connPool, opts...)
	return
}

func (m *managerImpl) setConnPool() {
	m.logger.Infof("creating conflict manager connection pool type=%s", m.connType)

	fn := m.newGocbCoreConn
	if m.connType == "memcached" {
		fn = m.newMemcachedConn
	}

	m.connPool = newConnPool(m.logger, fn)
	return
}

func (m *managerImpl) SetConnType(connType string) error {
	if m.connType == connType {
		return nil
	}

	m.logger.Infof("closing conflict manager connection pool type=%s", m.connType)
	err := m.connPool.Close()
	if err != nil {
		return err
	}

	m.connType = connType
	m.setConnPool()

	return nil
}

func (m *managerImpl) ConnPool() ConnPool {
	return m.connPool
}

func (m *managerImpl) newGocbCoreConn(bucketName string) (w io.Closer, err error) {
	m.logger.Infof("creating new conflict writer bucket=%s", bucketName)
	return NewGocbConn(m.logger, m.memdAddrGetter, bucketName)
}

func (m *managerImpl) newMemcachedConn(bucketName string) (conn io.Closer, err error) {
	addr, err := m.memdAddrGetter.MyMemcachedAddr()
	if err != nil {
		return
	}

	conn, err = NewMemcachedConn(m.logger, m.utils, m.manifestCache, bucketName, addr)

	return
}
