package conflictlog

import (
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
	connMap       map[string]*mcc.Client
	manifestCache *ManifestCache
	utilsObj      utils.UtilsIface
	bucketInfo    *BucketInfo
	opaque        uint32
}

func NewMemcachedConn(logger *log.CommonLogger, utilsObj utils.UtilsIface, manCache *ManifestCache, bucketName string, addr string) (m *MemcachedConn, err error) {
	user, passwd, err := cbauth.GetMemcachedServiceAuth(addr)
	if err != nil {
		return
	}

	logger.Debugf("memcached user=%s passwd=%s", user, passwd)

	connId := newConnId()
	conn, err := newMemcachedConn(logger, connId, utilsObj, bucketName, addr)
	if err != nil {
		return
	}

	m = &MemcachedConn{
		id:         connId,
		bucketName: bucketName,
		addr:       addr,
		logger:     logger,
		connMap: map[string]*mcc.Client{
			addr: conn,
		},
		utilsObj:      utilsObj,
		manifestCache: manCache,
	}

	return
}

func newMemcachedConn(logger *log.CommonLogger, id int64, utilsObj utils.UtilsIface, bucketName string, addr string) (conn *mcc.Client, err error) {
	user, passwd, err := cbauth.GetMemcachedServiceAuth(addr)
	if err != nil {
		return
	}

	logger.Infof("connecting to memcached id=%d user=%s addr=%s", id, user, addr)

	conn, err = base.NewConn(addr, user, passwd, bucketName, true, base.KeepAlivePeriod, logger)
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

	retFeatures, err := utilsObj.SendHELOWithFeatures(conn, userAgent, readTimeout, writeTimeout, features, logger)
	if err != nil {
		return
	}

	logger.Debugf("returned features: %s", retFeatures.String())

	return
}

func (m *MemcachedConn) fetchManifests(conn *mcc.Client) (man *metadata.CollectionsManifest, err error) {
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

func (m *MemcachedConn) getCollectionId(conn *mcc.Client, target Target, checkCache bool) (collId uint32, err error) {
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

func (m *MemcachedConn) setMeta(conn *mcc.Client, key string, vbNo uint16, body []byte, collId uint32, dataType uint8) (err error) {
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
		m.logger.Debugf("got UNKNOWN_COLLECTION id=%d, body=%s", m.id, string(rsp.Body))
		err = ErrUnknownCollection
		return
	}

	if rsp.Status == gomemcached.NOT_MY_VBUCKET {
		m.logger.Debugf("got NOT_MY_BUCKET id=%d, bucketName=%s", m.id, m.bucketName)
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

func (m *MemcachedConn) getConnByVB(vbno uint16, replicaNum int) (conn *mcc.Client, err error) {
	addr := m.addr
	if m.bucketInfo != nil {
		addr = m.bucketInfo.VBucketServerMap.GetAddrByVB(vbno, replicaNum)
	}

	m.logger.Debugf("selecting id=%d addr=%s for vb=%d", m.id, addr, vbno)
	conn, ok := m.connMap[addr]
	if ok {
		return
	}

	conn, err = newMemcachedConn(m.logger, m.id, m.utilsObj, m.bucketName, addr)
	if err != nil {
		return
	}
	m.connMap[addr] = conn
	return
}

func (m *MemcachedConn) SetMeta(key string, body []byte, dataType uint8, target Target) (err error) {
	checkCache := true
	var collId uint32
	vbNo := getVBNo(key, 1024)

	var conn *mcc.Client

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
