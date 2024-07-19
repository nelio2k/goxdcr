package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/couchbase/gocbcore/v9"
)

type GocbConn struct {
	MemcachedAddr string
	agent         *gocbcore.Agent
}

func newGocbConn(bucketName, user, passwd, memdAddr string) (conn *GocbConn, err error) {
	auth := gocbcore.PasswordAuthProvider{
		Username: user,
		Password: passwd,
	}

	config := &gocbcore.AgentConfig{
		MemdAddrs:      []string{memdAddr},
		Auth:           auth,
		BucketName:     bucketName,
		UserAgent:      "gocbcoreTest",
		UseCollections: true,
	}

	agent, err := gocbcore.CreateAgent(config)
	if err != nil {
		return
	}

	conn = &GocbConn{
		MemcachedAddr: memdAddr,
		agent:         agent,
	}

	return
}

func (conn *GocbConn) SetMeta(key string, val interface{}) (err error) {
	log.Printf("set: key=%s", key)

	body, err := json.Marshal(val)
	if err != nil {
		return
	}

	ch := make(chan error)

	opts := gocbcore.SetMetaOptions{
		Key:            []byte(key),
		Value:          body,
		Datatype:       1,
		ScopeName:      "S1",
		CollectionName: "C1",
		Cas:            gocbcore.Cas(time.Now().UnixNano()),
	}
	cb := func(sr *gocbcore.SetMetaResult, err2 error) {
		log.Printf("got set callback sr=%v, err2=%v", sr, err2)
		ch <- err2
	}

	_, err = conn.agent.SetMeta(opts, cb)
	if err != nil {
		return
	}

	err = <-ch

	return
}
