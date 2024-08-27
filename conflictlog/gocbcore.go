package conflictlog

import (
	"crypto/tls"
	"crypto/x509"
	"strings"
	"time"

	"github.com/couchbase/cbauth"
	"github.com/couchbase/gocbcore/v9"
	"github.com/couchbase/gocbcore/v9/memd"
	"github.com/couchbase/goxdcr/v8/base"
	"github.com/couchbase/goxdcr/v8/log"
)

var _ Connection = (*gocbCoreConn)(nil)

type gocbCoreConn struct {
	id             int64
	MemcachedAddr  string
	bucketName     string
	memdAddrGetter MemcachedAddrGetter
	agent          *gocbcore.Agent
	logger         *log.CommonLogger
	timeout        time.Duration
	finch          chan bool
	certs          *ClientCerts
}

func NewGocbConn(logger *log.CommonLogger, memdAddrGetter MemcachedAddrGetter, bucketName string, certs *ClientCerts) (conn *gocbCoreConn, err error) {
	connId := newConnId()

	logger.Infof("creating new gocbcore connection id=%d", connId)
	conn = &gocbCoreConn{
		id:             connId,
		memdAddrGetter: memdAddrGetter,
		bucketName:     bucketName,
		logger:         logger,
		//sudeep todo: make it configurable
		timeout: 60 * time.Second,
		finch:   make(chan bool),
		certs:   certs,
	}

	err = conn.setupAgent()
	if err != nil {
		conn = nil
	}

	return
}

func (conn *gocbCoreConn) getCACertPool() (*x509.CertPool, error) {
	caCert, err := conn.certs.LoadCACert()
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return caCertPool, nil
}

func (conn *gocbCoreConn) setupAgent() (err error) {
	memdAddr, err := conn.memdAddrGetter.MyMemcachedAddr()
	if err != nil {
		return
	}

	auth := &MemcachedAuthProvider{
		logger:      conn.logger,
		clientCerts: conn.certs,
	}

	var caCertProvider func() *x509.CertPool
	if conn.certs != nil {
		caPool, err := conn.getCACertPool()
		if err != nil {
			return err
		}

		caCertProvider = func() *x509.CertPool {
			return caPool
		}
	}

	config := &gocbcore.AgentConfig{
		MemdAddrs:              []string{memdAddr},
		Auth:                   auth,
		BucketName:             conn.bucketName,
		UserAgent:              MemcachedConnUserAgent,
		UseCollections:         true,
		UseTLS:                 conn.certs != nil,
		UseCompression:         true,
		AuthMechanisms:         []gocbcore.AuthMechanism{gocbcore.PlainAuthMechanism},
		TLSRootCAProvider:      caCertProvider,
		InitialBootstrapNonTLS: true,

		// use KvPoolSize=1 to ensure only one connection is created by the agent
		KvPoolSize: 1,
	}

	conn.agent, err = gocbcore.CreateAgent(config)
	if err != nil {
		return
	}

	signal := make(chan error, 1)
	_, err = conn.agent.WaitUntilReady(time.Now().Add(5*time.Second), gocbcore.WaitUntilReadyOptions{}, func(wr *gocbcore.WaitUntilReadyResult, err error) {
		conn.logger.Debugf("agent WaitUntilReady err=%v", err)
		signal <- err
	})

	err = <-signal

	return
}

func (conn *gocbCoreConn) Id() int64 {
	return conn.id
}

func (conn *gocbCoreConn) Bucket() string {
	return conn.bucketName
}

func (conn *gocbCoreConn) SetMeta(key string, body []byte, dataType uint8, target base.ConflictLogTarget) (err error) {
	ch := make(chan error)

	opts := gocbcore.SetMetaOptions{
		Key:            []byte(key),
		Value:          body,
		Datatype:       dataType,
		ScopeName:      target.NS.ScopeName,
		CollectionName: target.NS.CollectionName,
		Options:        uint32(memd.SkipConflictResolution),
		Cas:            gocbcore.Cas(time.Now().UnixNano()),
	}

	cb := func(sr *gocbcore.SetMetaResult, err2 error) {
		conn.logger.Debugf("got setMeta callback sr=%v, err2=%v", sr, err2)
		ch <- err2
	}

	_, err = conn.agent.SetMeta(opts, cb)
	if err != nil {
		return
	}

	select {
	case <-conn.finch:
		err = ErrWriterClosed
	case err = <-ch:
	case <-time.After(conn.timeout):
		err = ErrWriterTimeout
	}

	return
}

func (conn *gocbCoreConn) Close() error {
	close(conn.finch)
	return conn.agent.Close()
}

type MemcachedAuthProvider struct {
	logger      *log.CommonLogger
	clientCerts *ClientCerts
}

func (auth *MemcachedAuthProvider) Credentials(req gocbcore.AuthCredsRequest) (
	[]gocbcore.UserPassPair, error) {
	endpoint := req.Endpoint

	// get rid of the http:// or https:// prefix from the endpoint
	endpoint = strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
	username, password, err := cbauth.GetMemcachedServiceAuth(endpoint)
	if err != nil {
		return []gocbcore.UserPassPair{{}}, err
	}

	return []gocbcore.UserPassPair{{
		Username: username,
		Password: password,
	}}, nil
}

func (auth *MemcachedAuthProvider) SupportsNonTLS() bool {
	return !auth.SupportsTLS()
}

func (auth *MemcachedAuthProvider) SupportsTLS() bool {
	return auth.clientCerts != nil
}

func (auth *MemcachedAuthProvider) Certificate(req gocbcore.AuthCertRequest) (*tls.Certificate, error) {
	if auth.clientCerts == nil {
		return nil, nil
	}

	auth.logger.Infof("loading client certificates")

	clientCert, clientKey, err := auth.clientCerts.LoadCert()
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(clientCert, clientKey)
	return &cert, nil
}
