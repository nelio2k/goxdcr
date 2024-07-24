package conflictlog

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/couchbase/cbauth"
	"github.com/couchbase/gocbcore/v9"
	"github.com/couchbase/gocbcore/v9/memd"
	"github.com/couchbase/goxdcr/log"
)

const ConflictWriterUserAgent = "xdcrConflictWriter"

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
}

func NewGocbConn(logger *log.CommonLogger, memdAddrGetter MemcachedAddrGetter, bucketName string) (conn *gocbCoreConn, err error) {
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
	}

	err = conn.setupAgent()
	if err != nil {
		conn = nil
	}

	return
}

func (conn *gocbCoreConn) setupAgent() (err error) {
	memdAddr, err := conn.memdAddrGetter.MyMemcachedAddr()
	if err != nil {
		return
	}

	auth := &MemcachedAuthProvider{
		logger: conn.logger,
	}

	config := &gocbcore.AgentConfig{
		MemdAddrs:      []string{memdAddr},
		Auth:           auth,
		BucketName:     conn.bucketName,
		UserAgent:      ConflictWriterUserAgent,
		UseCollections: true,
		UseTLS:         false,
		UseCompression: true,
		AuthMechanisms: []gocbcore.AuthMechanism{gocbcore.PlainAuthMechanism},

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

func (conn *gocbCoreConn) SetMeta(key string, body []byte, dataType uint8, target Target) (err error) {
	//conn.logger.Infof("writing id=%d key=%s bodyLen=%d", conn.id, key, len(body))

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
	logger *log.CommonLogger
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
	return true
}

func (auth *MemcachedAuthProvider) SupportsTLS() bool {
	return false
}

func (auth *MemcachedAuthProvider) Certificate(req gocbcore.AuthCertRequest) (*tls.Certificate, error) {
	// If the internal client certificate has been set, use it for client authentication.
	return nil, nil
}
