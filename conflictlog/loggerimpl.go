package conflictlog

import (
	"fmt"
	"sync"
	"time"

	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/utils"
)

const (
	ConflictLoggerName string = "conflictLogger"
)

const DefaultLogCapacity = 5
const DefaultLoggerWorkerCount = 3
const DefaultNetworkRetryCount = 6
const DefaultNetworkRetryInterval = 10 * time.Second

var _ Logger = (*loggerImpl)(nil)

// loggerImpl implements the Logger interface
type loggerImpl struct {
	logger *log.CommonLogger

	// id uniquely identifies a logger instance
	id int64

	// utils object for misc utilities
	utils utils.UtilsIface

	// replId is the unique replication ID
	replId string

	// rules describe the mapping to the target conflict bucket
	rulesLock sync.RWMutex

	// connPool is used to get the connection to the cluster with conflict bucket
	connPool ConnPool

	// opts are the logger options
	opts LoggerOptions

	// mu is the logger level lock
	mu sync.Mutex

	// logReqCh is the work queue for all logging requests
	logReqCh chan logRequest

	// finch is intended to close all log workers
	finch chan bool

	// shutdownCh is intended to shut a subset of workers
	shutdownCh chan bool

	// wg is used to wait for outstanding log requests to be finished
	wg sync.WaitGroup
}

type logRequest struct {
	conflictRec *ConflictRecord
	//ackCh is the channel on which the result of the logging is reported back
	ackCh chan error
}

func newLoggerImpl(logger *log.CommonLogger, replId string, utils utils.UtilsIface, connPool ConnPool, opts ...LoggerOpt) (l *loggerImpl, err error) {
	// set the defaults
	options := LoggerOptions{
		rules:                nil,
		logQueueCap:          DefaultLogCapacity,
		workerCount:          DefaultLoggerWorkerCount,
		mapper:               NewConflictMapper(logger),
		networkRetryCount:    DefaultNetworkRetryCount,
		networkRetryInterval: DefaultNetworkRetryInterval,
	}

	// override the defaults
	for _, opt := range opts {
		opt(&options)
	}

	logger.Infof("creating new conflict logger replId=%s loggerOptions=%#v", replId, options)

	l = &loggerImpl{
		logger:    logger,
		id:        newLoggerId(),
		utils:     utils,
		replId:    replId,
		rulesLock: sync.RWMutex{},
		connPool:  connPool,
		opts:      options,
		mu:        sync.Mutex{},
		logReqCh:  make(chan logRequest, options.logQueueCap),
		finch:     make(chan bool, 1),
		// the value 10 is arbitrary. It basically means max 10 workers can be shutdown in parallel
		// The assumption is that >10 log workers would be very rare.
		shutdownCh: make(chan bool, 10),
		wg:         sync.WaitGroup{},
	}

	logger.Infof("spawning conflict logger workers replId=%s count=%d", l.replId, l.opts.workerCount)
	for i := 0; i < l.opts.workerCount; i++ {
		l.startWorker()
	}

	return
}

func (l *loggerImpl) Id() int64 {
	return l.id
}

func (l *loggerImpl) log(c *ConflictRecord) (ackCh chan error, err error) {
	ackCh = make(chan error, 1)
	req := logRequest{
		conflictRec: c,
		ackCh:       ackCh,
	}

	select {
	case <-l.finch:
		err = ErrLoggerClosed
	case l.logReqCh <- req:
	default:
		err = ErrQueueFull
	}

	return
}

func (l *loggerImpl) Log(c *ConflictRecord) (h base.ConflictLoggerHandle, err error) {
	l.logger.Debugf("logging conflict record replId=%s sourceKey=%s", l.replId, c.Source.Id)
	if l.isClosed() {
		err = ErrLoggerClosed
		return
	}

	ackCh, err := l.log(c)
	if err != nil {
		return
	}

	h = logReqHandle{
		ackCh: ackCh,
	}

	return
}

func (l *loggerImpl) isClosed() bool {
	select {
	case <-l.finch:
		return true
	default:
		return false
	}
}

func (l *loggerImpl) Close() (err error) {
	// check for closed channel as multiple threads could have attempted it
	if l.isClosed() {
		return nil
	}

	l.mu.Lock()

	// check for closed channel as multiple threads could have attempted it
	if l.isClosed() {
		return nil
	}

	close(l.finch)
	close(l.logReqCh)

	defer l.mu.Unlock()

	l.wg.Wait()

	return
}

func (l *loggerImpl) startWorker() {
	l.wg.Add(1)
	go l.worker()
}

func (l *loggerImpl) UpdateWorkerCount(newCount int) {
	l.logger.Infof("changing conflict logger worker count replId=%s old=%d new=%d", l.replId, l.opts.workerCount, newCount)

	l.mu.Lock()
	defer l.mu.Unlock()

	if newCount <= 0 || newCount == l.opts.workerCount {
		return
	}

	if newCount > l.opts.workerCount {
		for i := 0; i < (newCount - l.opts.workerCount); i++ {
			l.startWorker()
		}
	} else {
		for i := 0; i < (l.opts.workerCount - newCount); i++ {
			l.shutdownCh <- true
		}
	}

	l.opts.workerCount = newCount
}

func (l *loggerImpl) UpdateRules(r *Rules) (err error) {
	defer l.logger.Infof("Logger got the updated rules %v", r)

	if r != nil {
		// r is nil, meaning conflict logging is off.
		err = r.Validate()
		if err != nil {
			return
		}
	}

	l.rulesLock.Lock()
	l.opts.rules = r
	l.rulesLock.Unlock()

	return
}

func (l *loggerImpl) worker() {
	defer l.wg.Done()

	for {
		select {
		case <-l.shutdownCh:
			l.logger.Infof("shutting down conflict log worker replId=%s", l.replId)
			return
		case req := <-l.logReqCh:
			// nil implies that logReqCh might be closed
			if req.conflictRec == nil {
				return
			}

			err := l.processReq(req)
			req.ackCh <- err
		}
	}
}

func (l *loggerImpl) getTarget(rec *ConflictRecord) (t Target, err error) {
	l.rulesLock.RLock()
	defer l.rulesLock.RUnlock()

	t, err = l.opts.mapper.Map(l.opts.rules, rec)
	if err != nil {
		return
	}

	return
}

func (l *loggerImpl) getFromPool(bucketName string) (conn Connection, err error) {
	obj, err := l.connPool.Get(bucketName)
	if err != nil {
		return
	}

	conn, ok := obj.(Connection)
	if !ok {
		err = fmt.Errorf("pool object is of invalid type got=%T", obj)
		return
	}

	return
}

func (l *loggerImpl) writeDocs(req logRequest, target Target) (err error) {

	// Write source document.
	err = l.writeDocRetry(target.Bucket, func(conn Connection) error {
		err := conn.SetMeta(req.conflictRec.Source.Id, req.conflictRec.Source.Body, req.conflictRec.Source.Datatype, target)
		return err
	})
	if err != nil {
		return fmt.Errorf("error writing source doc, err=%v", err)
	}

	// Write target document.
	err = l.writeDocRetry(target.Bucket, func(conn Connection) error {
		err = conn.SetMeta(req.conflictRec.Target.Id, req.conflictRec.Target.Body, req.conflictRec.Target.Datatype, target)
		return err
	})
	if err != nil {
		return fmt.Errorf("error writing target doc, err=%v", err)
	}

	// Write conflict record.
	err = l.writeDocRetry(target.Bucket, func(conn Connection) error {
		err = conn.SetMeta(req.conflictRec.Id, req.conflictRec.body, req.conflictRec.datatype, target)
		return err
	})
	if err != nil {
		return fmt.Errorf("error writing conflict record, err=%v", err)
	}

	return
}

// writeDocRetry arranges for a connection from pool which the supplied function can use. The function wraps the
// supplied function to check for network errors and appropriately releases the connection back to the pool
func (l *loggerImpl) writeDocRetry(bucketName string, fn func(conn Connection) error) (err error) {
	var conn Connection

	for i := 0; i < l.opts.networkRetryCount; i++ {
		conn, err = l.getFromPool(bucketName)
		if err != nil {
			// This is to account for nw errors while connecting
			if !l.utils.IsSeriousNetError(err) {
				return
			}
			time.Sleep(l.opts.networkRetryInterval)
			continue
		}

		l.logger.Debugf("got connection from pool, id=%d", conn.Id())

		// The call to connPool.Put is not advised to be done in a defer here.
		// This is because calling defer in a loop accumulates multiple defers
		// across the iterations and calls all of them when function exits which
		// will be error prone in this case
		err = fn(conn)
		if err == nil {
			l.logger.Debugf("releasing connection to pool after success, id=%d damaged=%v", conn.Id(), false)
			l.connPool.Put(bucketName, conn, false)
			break
		}

		l.logger.Errorf("error in writing doc to conflict bucket err=%v", err)
		nwError := l.utils.IsSeriousNetError(err)
		l.logger.Debugf("releasing connection to pool after failure, id=%d damaged=%v", conn.Id(), nwError)
		l.connPool.Put(bucketName, conn, nwError)
		if !nwError {
			break
		}
		time.Sleep(l.opts.networkRetryInterval)
	}

	return err
}

func (l *loggerImpl) processReq(req logRequest) error {
	var err error

	err = req.conflictRec.PopulateData(l.replId)
	if err != nil {
		return err
	}

	target, err := l.getTarget(req.conflictRec)
	if err != nil {
		return err
	}

	err = l.writeDocs(req, target)
	if err != nil {
		return err
	}

	return nil
}

type logReqHandle struct {
	ackCh chan error
}

func (h logReqHandle) Wait(finch chan bool) (err error) {
	select {
	case <-finch:
		err = ErrLogWaitAborted
	case err = <-h.ackCh:
	}

	return
}