package conflictlog

import (
	"context"
	"io"

	"github.com/couchbase/cbauth"
	"github.com/couchbase/gocb/v2"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/service_def"
)

type Writer interface {
	io.Closer
	Upsert(ctx context.Context, c *ConflictRecord, target Target) (err error)
}

type writerImpl struct {
	logger      *log.CommonLogger
	xdcrTopoSvc service_def.XDCRCompTopologySvc
	cluster     *gocb.Cluster
}

func newWriterImpl(logger *log.CommonLogger, topo service_def.XDCRCompTopologySvc) (w *writerImpl, err error) {
	w = &writerImpl{
		logger:      logger,
		xdcrTopoSvc: topo,
	}
	err = w.connect()
	if err != nil {
		w = nil
	}

	return
}

func (w *writerImpl) Upsert(ctx context.Context, c *ConflictRecord, target Target) (err error) {
	err = target.Validate()
	if err != nil {
		return
	}

	bucket := w.cluster.Bucket(target.Bucket)
	var coll *gocb.Collection
	if target.Scope == "" {
		coll = bucket.DefaultCollection()
	} else {
		coll = bucket.Scope(target.Scope).Collection(target.Collection)
	}

	_, err = coll.Upsert(c.Id, c, nil)
	if err != nil {
		return
	}

	return
}

func (w *writerImpl) Close() (err error) {
	return w.cluster.Close(nil)
}

func (w *writerImpl) connect() (err error) {
	addr, err := w.xdcrTopoSvc.MyMemcachedAddr()
	if err != nil {
		return
	}

	user, passwd, err := cbauth.GetMemcachedServiceAuth(addr)
	if err != nil {
		return
	}

	c, err := gocb.Connect("couchbase://"+addr, gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: user,
			Password: passwd,
		},
	})

	if err != nil {
		return
	}

	w.cluster = c

	return
}
