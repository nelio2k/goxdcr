package conflictlog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/couchbase/cbauth"
	"github.com/couchbase/gocbcore/v9"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
)

const ConflictWriterUserAgent = "xdcrConflictWriter"

var _ Connection = (*gocbCoreConn)(nil)

type gocbCoreConn struct {
	MemcachedAddr  string
	bucketName     string
	memdAddrGetter MemcachedAddrGetter
	agent          *gocbcore.Agent
	logger         *log.CommonLogger
	timeout        time.Duration
	finch          chan bool
}

func newGocbConn(logger *log.CommonLogger, memdAddrGetter MemcachedAddrGetter, bucketName string) (conn *gocbCoreConn, err error) {
	conn = &gocbCoreConn{
		memdAddrGetter: memdAddrGetter,
		bucketName:     bucketName,
		logger:         logger,
		//sudeep todo: make it configurable
		timeout: 5 * time.Second,
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

	user, passwd, err := cbauth.GetMemcachedServiceAuth(memdAddr)
	if err != nil {
		fmt.Println("err=", err)
		return
	}

	auth := gocbcore.PasswordAuthProvider{
		Username: user,
		Password: passwd,
	}

	config := &gocbcore.AgentConfig{
		MemdAddrs:  []string{memdAddr},
		Auth:       auth,
		BucketName: conn.bucketName,
		UserAgent:  ConflictWriterUserAgent,

		// use KvPoolSize=1 to ensure only one connection is created by the agent
		KvPoolSize: 1,
	}

	conn.agent, err = gocbcore.CreateAgent(config)
	if err != nil {
		return
	}

	return
}

func (conn *gocbCoreConn) Bucket() string {
	return conn.bucketName
}

func (conn *gocbCoreConn) SetMeta(key string, body []byte, dataType uint8) (err error) {
	ch := make(chan error)

	opts := gocbcore.SetOptions{
		Key:      []byte(key),
		Value:    body,
		Datatype: dataType,
	}

	cb := func(sr *gocbcore.StoreResult, err2 error) {
		conn.logger.Infof("got set callback sr=%v, err2=%v", sr, err2)
		ch <- err2
	}

	_, err = conn.agent.Set(opts, cb)
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

func (conn *gocbCoreConn) SetMetaObj(key string, val interface{}) (err error) {
	conn.logger.Infof("set: key=%s", key)

	docBody, isDocBody := val.([]byte)

	var body []byte
	if isDocBody {
		body = docBody
	} else {
		body, err = json.Marshal(val)
		if err != nil {
			return
		}
	}

	return conn.SetMeta(key, body, base.JSONDataType)
}

func (conn *gocbCoreConn) Close() error {
	close(conn.finch)
	return conn.agent.Close()
}
