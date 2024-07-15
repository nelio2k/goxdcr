package conflictlog

import (
	"sync"

	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
)

const DefaultLogCapacity = 5

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
	ackCh = make(chan error)
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
	close(l.finch)
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
	err = r.Validate()
	if err != nil {
		return
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

func (l *loggerImpl) processReq(req logRequest) (err error) {
	target, err := l.getTarget(req.conflictRec)
	if err != nil {
		return
	}

	w, err := l.writerPool.get(target.Bucket)
	if err != nil {
		return
	}

	defer func() {
		l.writerPool.release(w)
	}()

	err = w.SetMetaObj(req.conflictRec.Id, req.conflictRec)
	return
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
