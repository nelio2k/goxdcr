package conflictlog

import (
	"encoding/binary"
	"fmt"
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
	logger        *log.CommonLogger
	bucketName    string
	mccConn       *mcc.Client
	manifestCache *ManifestCache
	opaque        uint32
}

func NewMemcachedConn(logger *log.CommonLogger, utilsObj utils.UtilsIface, manCache *ManifestCache, bucketName string, addr string) (m *MemcachedConn, err error) {
	user, passwd, err := cbauth.GetMemcachedServiceAuth(addr)
	if err != nil {
		return
	}

	logger.Debugf("memcached user=%s passwd=%s", user, passwd)

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

	logger.Debugf("returned features: %s", retFeatures.String())

	m = &MemcachedConn{
		logger:        logger,
		mccConn:       mccConn,
		manifestCache: manCache,
	}

	return
}

func (m *MemcachedConn) fetchManifests() (man *metadata.CollectionsManifest, err error) {
	rsp, err := m.mccConn.GetCollectionsManifest()
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

func (m *MemcachedConn) getCollectionId(target Target, checkCache bool) (collId uint32, err error) {
	var ok bool

	if checkCache {
		collId, ok = m.manifestCache.GetCollectionId(target.Bucket, target.Scope, target.Collection)
		if ok {
			return
		}
	}

	m.logger.Infof("fetching manifests for checkCache=%v bucket=%s", checkCache, target.Bucket)

	man, err := m.fetchManifests()
	if err != nil {
		return 0, err
	}

	collId, err = man.GetCollectionId(target.Scope, target.Collection)
	if err != nil {
		if err == base.ErrorNotFound {
			err = fmt.Errorf("scope or collection not found. target=%s", target)
		}
		return 0, err
	}

	m.manifestCache.UpdateManifest(target.Bucket, man)
	return
}

func (m *MemcachedConn) setMeta(key string, body []byte, collId uint32, dataType uint8) (err error) {
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
	vbNo := getVBNo(key, 64)

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
	//options |= base.SKIP_CONFLICT_RESOLUTION_FLAG
	binary.BigEndian.PutUint32(req.Extras[0:4], 0)
	binary.BigEndian.PutUint64(req.Extras[8:16], 0)
	binary.BigEndian.PutUint64(req.Extras[16:24], cas)
	binary.BigEndian.PutUint32(req.Extras[24:28], options)

	//conn.logger.Debugf("bytes = %v", req.Bytes())

	//rsp, err := m.mccConn.Set(vbNo, key, 0, 0, val, &mcc.ClientContext{CollId: colId})
	rsp, err := m.mccConn.Send(req)
	if rsp != nil {
		if rsp.Opaque != opaque {
			err = fmt.Errorf("opaque value mismatch expected=%d,got=%d", opaque, rsp.Opaque)
			return
		}
		if rsp.Status == gomemcached.UNKNOWN_COLLECTION {
			err = ErrUnknownCollection
			return
		}
		m.logger.Infof("received rsp key=%s status=%d", rsp.Key, rsp.Status)
	}
	if err != nil {
		return
	}
	return
}

func (m *MemcachedConn) SetMeta(key string, body []byte, dataType uint8, target Target) (err error) {
	checkCache := true
	var collId uint32

	for i := 0; i < 2; i++ {
		collId, err = m.getCollectionId(target, checkCache)
		if err != nil {
			return err
		}

		err = m.setMeta(key, body, collId, dataType)
		if err == nil {
			return
		}

		if err == ErrUnknownCollection {
			m.logger.Infof("collection not found key=%s, target=%s", key, target.String())
			checkCache = false
			continue
		} else {
			return err
		}
	}

	return
}

func (m *MemcachedConn) Close() error {
	return nil
}
