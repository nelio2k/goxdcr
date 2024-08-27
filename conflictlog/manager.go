package conflictlog

import (
	"fmt"
	"io"
	"os"

	"github.com/couchbase/goxdcr/v8/log"
	"github.com/couchbase/goxdcr/v8/utils"
)

var manager Manager

var _ Manager = (*managerImpl)(nil)

// Manager defines behaviour for conflict manager
type Manager interface {
	NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error)
	ConnPool() ConnPool
	// [TEMP]: SetConnType exists only for perf test
	SetConnType(connType string) error
	SetConnLimit(limit int)
	SetSkipTlsVerify(bool)
}

// ClientCerts are the certs created and owned by ns_server for node-2-node encrypted
// connection
type ClientCerts struct {
	ClientCertFile string
	ClientKeyFile  string
	ClusterCAFile  string
}

func (c *ClientCerts) LoadCert() (clientCert, clientKey []byte, err error) {
	clientCert, err = os.ReadFile(c.ClientCertFile)
	if err != nil {
		return
	}

	clientKey, err = os.ReadFile(c.ClientKeyFile)
	if err != nil {
		return
	}

	return
}

func (c *ClientCerts) LoadCACert() (caCert []byte, err error) {
	caCert, err = os.ReadFile(c.ClusterCAFile)
	return
}

func (c *ClientCerts) Load() (clientCert, clientKey, caCert []byte, err error) {
	clientCert, clientKey, err = c.LoadCert()
	if err != nil {
		return
	}

	caCert, err = os.ReadFile(c.ClusterCAFile)
	return
}

type MemcachedAddrGetter interface {
	MyMemcachedAddr() (string, error)
}

type EncryptionInfoGetter interface {
	IsMyClusterEncryptionLevelStrict() bool
}

// GetManager returns the global conflict manager
func GetManager() (Manager, error) {
	if manager == nil {
		return nil, ErrManagerNotInitialized
	}

	return manager, nil
}

// InitManager intializes global conflict manager
func InitManager(loggerCtx *log.LoggerContext, utils utils.UtilsIface, memdAddrGetter MemcachedAddrGetter, encInfoGetter EncryptionInfoGetter, certs *ClientCerts) {
	logger := log.NewLogger(ConflictManagerLoggerName, loggerCtx)

	logger.Info("intializing conflict manager")

	impl := &managerImpl{
		logger:         logger,
		memdAddrGetter: memdAddrGetter,
		encInfoGetter:  encInfoGetter,
		utils:          utils,
		manifestCache:  newManifestCache(),
		connLimit:      DefaultPoolConnLimit,
		connType:       "gocbcore",
		certs:          certs,
	}

	logger.Info("creating conflict manager writer pool")
	impl.setConnPool()

	manager = impl
}

// managerImpl implements conflict manager
type managerImpl struct {
	logger         *log.CommonLogger
	memdAddrGetter MemcachedAddrGetter
	encInfoGetter  EncryptionInfoGetter
	utils          utils.UtilsIface
	connPool       *connPool
	manifestCache  *ManifestCache

	// connLimit max number of connections
	connLimit int

	// [TEMP] connType only exists for perf test
	connType string

	certs         *ClientCerts
	skipTlsVerify bool
}

func (m *managerImpl) NewLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l Logger, err error) {
	opts = append(opts, WithSkipTlsVerify(m.skipTlsVerify))
	l, err = newLoggerImpl(logger, replId, m.utils, m.connPool, opts...)
	return
}

func (m *managerImpl) SetConnLimit(limit int) {
	m.logger.Infof("setting connection limit = %d", limit)
	m.connPool.SetLimit(limit)
}

func (m *managerImpl) SetSkipTlsVerify(v bool) {
	m.logger.Infof("setting tls skip verify=%v", v)
	m.skipTlsVerify = v
}

func (m *managerImpl) setConnPool() {
	m.logger.Infof("creating conflict manager connection pool type=%s", m.connType)

	fn := m.newGocbCoreConn
	if m.connType == "memcached" {
		fn = m.newMemcachedConn
	}

	m.connPool = newConnPool(m.logger, m.connLimit, fn)
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
	m.logger.Infof("creating new conflict writer bucket=%s certsEnabled=%v", bucketName, m.certs != nil)

	certs, err := m.getCerts()
	if err != nil {
		return nil, err
	}

	return NewGocbConn(m.logger, m.memdAddrGetter, bucketName, certs)
}

func (m *managerImpl) newMemcachedConn(bucketName string) (conn io.Closer, err error) {
	addr, err := m.memdAddrGetter.MyMemcachedAddr()
	if err != nil {
		return
	}

	certs, err := m.getCerts()
	if err != nil {
		return nil, err
	}

	conn, err = NewMemcachedConn(m.logger, m.utils, m.manifestCache, bucketName, addr, certs, m.skipTlsVerify)

	return
}

func (m *managerImpl) getCerts() (*ClientCerts, error) {
	isStrict := m.encInfoGetter.IsMyClusterEncryptionLevelStrict()
	m.logger.Infof("is cluster encryption strict = %v", isStrict)
	if !isStrict {
		return nil, nil
	}

	if m.encInfoGetter.IsMyClusterEncryptionLevelStrict() {
		if m.certs == nil {
			return nil, fmt.Errorf("Client certs not provided when encryption is mandatory")
		}
	}

	return m.certs, nil
}
