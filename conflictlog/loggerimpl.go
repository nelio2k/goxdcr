package conflictlog

import (
	"fmt"
	"sync"

	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
)

const (
	ConflictLoggerName string = "conflictLogger"
)

const DefaultLogCapacity = 5
const DefaultRetryCntOnWriteFailure = 5

var _ Logger = (*loggerImpl)(nil)

// loggerImpl implements the Logger interface
type loggerImpl struct {
	// replId is the unique replication ID
	replId string

	// mapper maps the conflict to the target conflict bucket
	mapper Mapper

	// rules describe the mapping to the target conflict bucket
	rules     *Rules
	rulesLock sync.RWMutex

	writerPool *writerPool
	logger     *log.CommonLogger

	workerCount int
	loggerLock  sync.Mutex

	logCh      chan logRequest
	finch      chan bool
	shutdownCh chan bool

	// Logger can be shared between different nozzles,
	// hence it should be closed by only one of them when the pipeline is stopping.
	closed bool
}

type logRequest struct {
	conflictRec *ConflictRecord
	ackCh       chan error
}

func WithRules(r *Rules) LoggerOpt {
	return func(o *LoggerOptions) {
		o.rules = r
	}
}

func WithMapper(m Mapper) LoggerOpt {
	return func(o *LoggerOptions) {
		o.mapper = m
	}
}

func WithCapacity(cap int) LoggerOpt {
	return func(o *LoggerOptions) {
		o.logQueueCap = cap
	}
}

func newLoggerImpl(logger *log.CommonLogger, replId string, writerPool *writerPool, opts ...LoggerOpt) (l *loggerImpl, err error) {
	options := &LoggerOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.logQueueCap <= 0 {
		options.logQueueCap = DefaultLogCapacity
	}

	if options.workerCount <= 0 {
		options.workerCount = 3
	}

	if options.mapper == nil {
		options.mapper = NewFixedMapper(logger, Target{Bucket: "B1"})
	}

	logger.Infof("creating new conflict logger replId=%s loggerOptions=%#v", replId, options)

	l = &loggerImpl{
		logger:      logger,
		replId:      replId,
		rules:       options.rules,
		rulesLock:   sync.RWMutex{},
		writerPool:  writerPool,
		mapper:      options.mapper,
		workerCount: options.workerCount,
		loggerLock:  sync.Mutex{},
		logCh:       make(chan logRequest, options.logQueueCap),
		finch:       make(chan bool),
		shutdownCh:  make(chan bool, 10),
	}

	logger.Infof("spawning conflict logger workers replId=%s count=%d", l.replId, l.workerCount)
	for i := 0; i < l.workerCount; i++ {
		go l.worker()
	}

	return
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
	case l.logCh <- req:
	default:
		err = ErrQueueFull
	}

	if err != nil {
		return
	}

	return
}

func (l *loggerImpl) Log(c *ConflictRecord) (h base.ConflictLoggerHandle, err error) {
	l.logger.Infof("logging conflict record replId=%s sourceKey=%s", l.replId, c.Source.Id)

	ackCh, err := l.log(c)

	h = logReqHandle{
		ackCh: ackCh,
	}

	return
}

func (l *loggerImpl) Close() (err error) {
	l.loggerLock.Lock()
	defer l.loggerLock.Unlock()

	if !l.closed {
		l.closed = true
		close(l.finch)
	}

	return
}

func (l *loggerImpl) UpdateWorkerCount(newCount int) {
	l.logger.Infof("changing conflict logger worker count replId=%s old=%d new=%d", l.replId, l.workerCount, newCount)

	l.loggerLock.Lock()
	defer l.loggerLock.Unlock()

	if newCount <= 0 || newCount == l.workerCount {
		return
	}

	if newCount > l.workerCount {
		for i := 0; i < (newCount - l.workerCount); i++ {
			go l.worker()
		}
	} else {
		for i := 0; i < (l.workerCount - newCount); i++ {
			l.shutdownCh <- true
		}
	}

	l.workerCount = newCount
}

func (l *loggerImpl) UpdateRules(r *Rules) (err error) {
	defer l.logger.Infof("Logger got the updated rules %s", r)

	if r != nil {
		// r is nil, meaning conflict logging is off.
		err = r.Validate()
		if err != nil {
			return
		}
	}

	l.rulesLock.Lock()
	l.rules = r
	l.rulesLock.Unlock()

	return
}

func (l *loggerImpl) worker() {
	for {
		select {
		case <-l.finch:
			return
		case <-l.shutdownCh:
			l.logger.Infof("shutting down conflict log worker replId=%s", l.replId)
			return
		case req := <-l.logCh:
			err := l.processReq(req)
			req.ackCh <- err
		}
	}
}

func (l *loggerImpl) getTarget(rec *ConflictRecord) (t Target, err error) {
	l.rulesLock.RLock()
	defer l.rulesLock.RUnlock()

	t, err = l.mapper.Map(l.rules, rec)
	if err != nil {
		return
	}

	return
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

	w, err := l.writerPool.get(target.Bucket)
	if err != nil {
		return err
	}

	defer func() {
		l.writerPool.release(w)
	}()

	// Write source document.
	for i := 0; i < DefaultRetryCntOnWriteFailure; i++ {
		err = w.SetMeta(req.conflictRec.Source.Id, req.conflictRec.Source.body, req.conflictRec.Source.Datatype, target)
		if err != nil {
			continue
		}
	}
	if err != nil {
		return fmt.Errorf("error writing source doc, err=%v", err)
	}

	// Write target document.
	for i := 0; i < DefaultRetryCntOnWriteFailure; i++ {
		err = w.SetMeta(req.conflictRec.Target.Id, req.conflictRec.Target.body, req.conflictRec.Target.Datatype, target)
		if err != nil {
			continue
		}
	}
	if err != nil {
		return fmt.Errorf("error writing target doc, err=%v", err)
	}

	// Write conflict record.
	// Write target document.
	for i := 0; i < DefaultRetryCntOnWriteFailure; i++ {
		err = w.SetMeta(req.conflictRec.Id, req.conflictRec.body, req.conflictRec.datatype, target)
		if err != nil {
			continue
		}
	}
	if err != nil {
		return fmt.Errorf("error writing conflict record, err=%v", err)
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
