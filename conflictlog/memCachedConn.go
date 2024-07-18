package conflictlog

import (
	"time"

	"github.com/couchbase/cbauth"
	"github.com/couchbase/gomemcached"
	mcc "github.com/couchbase/gomemcached/client"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/utils"
)

//var _ Connection = (*MemcachedConn)(nil)

type MemcachedConn struct {
	logger     *log.CommonLogger
	bucketName string
	mccConn    *mcc.Client
}

func NewMemcachedConn(logger *log.CommonLogger, utilsObj utils.UtilsIface, bucketName string, addr string) (m *MemcachedConn, err error) {
	user, passwd, err := cbauth.GetMemcachedServiceAuth(addr)
	if err != nil {
		return
	}

	mccConn, err := base.NewConn(addr, user, passwd, bucketName, true, base.KeepAlivePeriod, logger)
	if err != nil {
		return
	}

	var features utils.HELOFeatures
	features.Xattribute = true
	features.Xerror = true
	features.Collections = true
	features.DataType = true

	// For setMeta, negotiate compression, if it is set
	features.CompressionType = base.CompressionTypeSnappy
	userAgent := "conflictLogConn"
	readTimeout := 30 * time.Second
	writeTimeout := 30 * time.Second

	retFeatures, err := utilsObj.SendHELOWithFeatures(mccConn, userAgent, readTimeout, writeTimeout, features, logger)
	if err != nil {
		return
	}

	logger.Debugf("returned features: %v", retFeatures)

	m = &MemcachedConn{
		logger:  logger,
		mccConn: mccConn,
	}

	return
}

func (m *MemcachedConn) SetMeta(key string, val []byte, dataType uint8) (err error) {
	req := gomemcached.MCRequest{
		Opcode:   base.SET_WITH_META,
		Key:      []byte(key),
		Keylen:   len(key),
		Body:     val,
		DataType: dataType,
	}

	rsp, err := m.mccConn.Send(&req)
	if err != nil {
		return
	}

	m.logger.Infof("received rsp key=%s status=%d", rsp.Key, rsp.Status)
	return
}

func (m *MemcachedConn) SetMetaObj(key string, obj interface{}) (err error) {
	return
}

func (m *MemcachedConn) Close() error {
	return nil
}
