package conflictlog

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/couchbase/cbauth"
	"github.com/couchbase/gomemcached"
	mcc "github.com/couchbase/gomemcached/client"
	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/metadata"
	"github.com/couchbase/goxdcr/utils"
)

var _ Connection = (*MemcachedConn)(nil)

type MemcachedConn struct {
	id            int64
	addr          string
	logger        *log.CommonLogger
	bucketName    string
	connMap       map[string]mcc.ClientIface
	manifestCache *ManifestCache
	utilsObj      utils.UtilsIface
	bucketInfo    *BucketInfo
	opaque        uint32
	certs         *ClientCerts
	skipVerify    bool
}

func NewMemcachedConn(logger *log.CommonLogger, utilsObj utils.UtilsIface, manCache *ManifestCache, bucketName string, addr string, certs *ClientCerts, skipVerifiy bool) (m *MemcachedConn, err error) {
	connId := newConnId()

	user, passwd, err := cbauth.GetMemcachedServiceAuth(addr)
	if err != nil {
		return
	}

	conn, err := newMemcNodeConn(logger, connId, utilsObj, user, passwd, bucketName, addr, nil, skipVerifiy)
	if err != nil {
		return
	}

	m = &MemcachedConn{
		id:         connId,
		bucketName: bucketName,
		addr:       addr,
		logger:     logger,
		connMap: map[string]mcc.ClientIface{
			addr: conn,
		},
		utilsObj:      utilsObj,
		manifestCache: manCache,
		certs:         certs,
		skipVerify:    skipVerifiy,
	}

	return
}

// newTLSConn creates an SSL connection to a memcached node.
// Note: this is subtly different than base.NewTlsConn(). In this we use cbauth's user/passwd
// over SSL connection instead of client certificate's SAN
func newTLSConn(addr string, certs *ClientCerts, user, passwd, bucketName string, skipVerification bool) (conn mcc.ClientIface, err error) {
	// We load the certs everytime since the certs could have changed by ns_server
	// This will be actually not needed as these will be loaded only when the notification
	// from the ns_server. This will happen in security service.
	clientCert, clientKey, caCert, err := certs.Load()
	if err != nil {
		return nil, err
	}

	x509Cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipVerification,
		Certificates:       []tls.Certificate{x509Cert},
		RootCAs:            caCertPool,
	}

	tlsConn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return nil, err
	}

	conn, err = mcc.Wrap(tlsConn)
	if err != nil {
		tlsConn.Close()
		return nil, err
	}

	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	_, err = conn.Auth(user, passwd)
	if err != nil {
		return nil, err
	}

	_, err = conn.SelectBucket(bucketName)
	if err != nil {
		return nil, err
	}
	return
}

// newMemcNodeConn creates connection to a given memcached node
func newMemcNodeConn(logger *log.CommonLogger, id int64, utilsObj utils.UtilsIface, user, passwd, bucketName string, addr string, certs *ClientCerts, skipVerify bool) (conn mcc.ClientIface, err error) {
	logger.Infof("connecting to memcached id=%d user=%s addr=%s certsEnabled=%v tlsSkipVerify=%v", id, user, addr, certs != nil, skipVerify)

	if certs != nil {
		conn, err = newTLSConn(addr, certs, user, passwd, bucketName, skipVerify)
	} else {
		conn, err = base.NewConn(addr, user, passwd, bucketName, true, base.KeepAlivePeriod, logger)
	}

	if err != nil {
		return nil, err
	}

	var features utils.HELOFeatures
	features.Xattribute = true
	features.Xerror = true
	features.Collections = true
	features.DataType = true

	// For setMeta, negotiate compression, if it is set
	features.CompressionType = base.CompressionTypeSnappy
	userAgent := MemcachedConnUserAgent
	readTimeout := 30 * time.Second
	writeTimeout := 30 * time.Second

	retFeatures, err := utilsObj.SendHELOWithFeatures(conn, userAgent, readTimeout, writeTimeout, features, logger)
	if err != nil {
		return
	}

	logger.Debugf("returned features: %s", retFeatures.String())

	return
}

func (m *MemcachedConn) fetchManifests(conn mcc.ClientIface) (man *metadata.CollectionsManifest, err error) {
	rsp, err := conn.GetCollectionsManifest()
	if err != nil {
		return
	}

	if rsp.Status != gomemcached.SUCCESS {
		err = fmt.Errorf("memcached request failed, req=GetCollectionManifest, status=%d, msg=%s", rsp.Status, string(rsp.Body))
		return
	}

	man = &metadata.CollectionsManifest{}
	err = man.UnmarshalJSON(rsp.Body)

	return
}

func (conn *MemcachedConn) Id() int64 {
	return conn.id
}

// getCollectionId first attempts to get the collectionId from the cache (if checkCache=true). If not found then
// it attempt to fetch it from the cluster using the same memcached connection. checkCache=false is generally used
// when we know that the value is cache is stale and a fresh one has to be fetched.
func (m *MemcachedConn) getCollectionId(conn mcc.ClientIface, target base.ConflictLogTarget, checkCache bool) (collId uint32, err error) {
	var ok bool

	if checkCache {
		collId, ok = m.manifestCache.GetCollectionId(target.Bucket, target.NS.ScopeName, target.NS.CollectionName)
		if ok {
			return
		}
	}

	m.logger.Infof("fetching manifests for checkCache=%v bucket=%s", checkCache, target.Bucket)

	man, err := m.fetchManifests(conn)
	if err != nil {
		return 0, err
	}

	collId, err = man.GetCollectionId(target.NS.ScopeName, target.NS.CollectionName)
	if err != nil {
		if err == base.ErrorNotFound {
			err = fmt.Errorf("scope or collection not found. target=%s", target)
		}
		return 0, err
	}

	m.manifestCache.UpdateManifest(target.Bucket, man)
	return
}

func (m *MemcachedConn) setMeta(conn mcc.ClientIface, key string, vbNo uint16, body []byte, collId uint32, dataType uint8) (err error) {
	bufGetter := func(sz uint64) ([]byte, error) {
		return make([]byte, sz), nil
	}

	encCid, encLen, err := base.NewUleb128(collId, bufGetter, true)
	if err != nil {
		return
	}

	totalLen := encLen + len(key)
	keybuf := make([]byte, totalLen)
	copy(keybuf[0:encLen], encCid[0:encLen])
	copy(keybuf[encLen:], []byte(key))

	m.logger.Debugf("vbNo=%d encCid: %v, len=%d, keybuf:%v", vbNo, encCid[0:encLen], totalLen, keybuf)

	cas := uint64(time.Now().UnixNano())
	opaque := atomic.AddUint32(&m.opaque, 1)

	req := &gomemcached.MCRequest{
		Opcode:   base.SET_WITH_META,
		VBucket:  vbNo,
		Key:      keybuf,
		Keylen:   totalLen,
		Body:     body,
		Opaque:   opaque,
		Cas:      0,
		Extras:   make([]byte, 30),
		DataType: dataType,
	}

	var options uint32
	options |= base.SKIP_CONFLICT_RESOLUTION_FLAG
	binary.BigEndian.PutUint32(req.Extras[0:4], 0)
	binary.BigEndian.PutUint64(req.Extras[8:16], 0)
	binary.BigEndian.PutUint64(req.Extras[16:24], cas)
	binary.BigEndian.PutUint32(req.Extras[24:28], options)

	rsp, err := conn.Send(req)
	err2 := m.handleResponse(rsp, opaque)
	if err2 != nil {
		return err2
	}

	if err != nil {
		return
	}
	return
}

func (m *MemcachedConn) handleResponse(rsp *gomemcached.MCResponse, opaque uint32) (err error) {
	if rsp == nil {
		return
	}

	if rsp.Opaque != opaque {
		err = fmt.Errorf("opaque value mismatch expected=%d,got=%d", opaque, rsp.Opaque)
		return
	}

	if rsp.Status == gomemcached.UNKNOWN_COLLECTION {
		m.logger.Debugf("got unknown_collection id=%d, body=%s", m.id, string(rsp.Body))
		err = ErrUnknownCollection
		return
	}

	if rsp.Status == gomemcached.NOT_MY_VBUCKET {
		m.logger.Debugf("got not_my_vbucket id=%d, bucketName=%s", m.id, m.bucketName)
		m.bucketInfo, err = parseNotMyVbucketValue(m.logger, rsp.Body, m.addr)
		if err != nil {
			return
		}
		err = ErrNotMyBucket
		return
	}

	m.logger.Debugf("received rsp key=%s status=%d", rsp.Key, rsp.Status)
	return
}

// getConnByVB gets (or creates) connection to vbNo's memcached node
func (m *MemcachedConn) getConnByVB(vbno uint16, replicaNum int) (conn mcc.ClientIface, err error) {
	// The logic is as follows:
	//    We use non-tls addr if connecting to 'thisNode' (aka localhost).
	//    For everything else it depends if certs are enabled or not
	//    m.bucketInfo == nil implies that so far we have not received NOT_MY_VBUCKET error.

	addr2use := m.addr
	if m.bucketInfo != nil {
		_, hostname, port, sslPort, thisNode, err := m.bucketInfo.GetAddrByVB(vbno, replicaNum)
		if err != nil {
			return nil, err
		}

		addr2use = fmt.Sprintf("%s:%d", hostname, port)
		if !thisNode && m.certs != nil {
			addr2use = fmt.Sprintf("%s:%d", hostname, sslPort)
		}
	}

	user, passwd, err := cbauth.GetMemcachedServiceAuth(addr2use)
	if err != nil {
		return
	}

	m.logger.Debugf("selecting id=%d certs=%v addr2use=%s for vb=%d", m.id, m.certs != nil, addr2use, vbno)
	conn, ok := m.connMap[addr2use]
	if ok {
		return
	}

	conn, err = newMemcNodeConn(m.logger, m.id, m.utilsObj, user, passwd, m.bucketName, addr2use, m.certs, m.skipVerify)
	if err != nil {
		return
	}

	m.connMap[addr2use] = conn
	return
}

func (m *MemcachedConn) SetMeta(key string, body []byte, dataType uint8, target base.ConflictLogTarget) (err error) {
	checkCache := true
	var collId uint32
	vbNo := base.GetVBucketNo(key, base.NumberOfVbs)

	var conn mcc.ClientIface

	for i := 0; i < 2; i++ {
		conn, err = m.getConnByVB(vbNo, 0)
		if err != nil {
			return err
		}

		collId, err = m.getCollectionId(conn, target, checkCache)
		if err != nil {
			return err
		}

		err = m.setMeta(conn, key, vbNo, body, collId, dataType)
		if err == nil {
			return
		}

		switch err {
		case ErrUnknownCollection:
			m.logger.Infof("collection not found key=%s, target=%s", key, target.String())
			checkCache = false
		case ErrNotMyBucket:
		default:
			return err
		}
	}

	return
}

func (m *MemcachedConn) Close() error {
	m.logger.Infof("closing memcached conn id=%d", m.id)
	for _, conn := range m.connMap {
		conn.Close()
	}
	return nil
}

func parseNotMyVbucketValue(logger *log.CommonLogger, value []byte, sourceAddr string) (info *BucketInfo, err error) {
	logger.Debugf("parsing NOT_MY_BUCKET response")

	sourceHost := base.GetHostName(sourceAddr)
	// Try to parse the value as a bucket configuration
	info, err = parseConfig(value, sourceHost)
	return
}

func parseConfig(config []byte, srcHost string) (info *BucketInfo, err error) {
	configStr := strings.Replace(string(config), "$HOST", srcHost, -1)

	info = new(BucketInfo)
	err = json.Unmarshal([]byte(configStr), info)
	if err != nil {
		return nil, err
	}

	return info, nil
}
