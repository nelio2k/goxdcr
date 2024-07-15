package conflictlog

import (
	"fmt"
	"sync"
	"time"

	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
)

var _ Logger = (*fileLoggerImpl)(nil)

var defaultWaitTimeout time.Duration = 5 * time.Second

// fileLoggerImpl implements the Logger interface.
// Conflicts will be logged to a file.
// Used as a POC.
type fileLoggerImpl struct {
	// replId is the unique replication ID
	replId string

	// mapper maps the conflict to the target conflict bucket
	mapper Mapper

	// rules describe the mapping to the target conflict bucket
	rules     *Rules
	rulesLock sync.RWMutex

	logger *log.CommonLogger

	workerCount int
	loggerLock  sync.Mutex

	logCh      chan logRequest
	finch      chan bool
	shutdownCh chan bool
}

func NewFileLogger(logger *log.CommonLogger, replId string, opts ...LoggerOpt) (l *fileLoggerImpl) {
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
		// should not reach here in production.
		options.mapper = NewFixedMapper(logger, Target{Bucket: "B1"})
	}

	logger.Infof("creating new conflict logger replId=%s loggerOptions=%#v", replId, options)

	l = &fileLoggerImpl{
		logger:      logger,
		replId:      replId,
		rules:       options.rules,
		rulesLock:   sync.RWMutex{},
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

func (l *fileLoggerImpl) Log(c *ConflictRecord) (h base.ConflictLoggerHandle, err error) {
	ackCh := make(chan error)
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

	h = logReqHandleWithTimeout{
		ackCh: ackCh,
	}

	return
}

func (l *fileLoggerImpl) Close() (err error) {
	close(l.finch)
	return
}

func (l *fileLoggerImpl) UpdateWorkerCount(newCount int) {
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
}

func (l *fileLoggerImpl) UpdateRules(r *Rules) (err error) {
	if r == nil {
		return
	}

	err = r.Validate()
	if err != nil {
		return
	}

	l.rulesLock.Lock()
	l.rules = r
	l.rulesLock.Unlock()

	return
}

func (l *fileLoggerImpl) worker() {
	for {
		select {
		case <-l.finch:
			return
		case <-l.shutdownCh:
			l.logger.Infof("shutting down conflict log worker replId=%s", l.replId)
			return
		case req := <-l.logCh:
			err := l.processReq(req)

			// if no one is waiting, log the error in goxdcr.log and proceed to not block the worker.
			// Also close the ackCh so that if someone calls Wait later, they are not blocked too.
			select {
			case req.ackCh <- err:
			default:
				if err != nil {
					l.logger.Errorf("Error logging %s, err=%v. No one waiting yet.", req.conflictRec, err)
				}
			}
			close(req.ackCh)
		}
	}
}

func (l *fileLoggerImpl) processReq(req logRequest) (err error) {
	// copy the pointer for rules and from this point, the rules can can
	// be updated but the old rules will be used for this request
	var rules *Rules
	l.rulesLock.Lock()
	rules = l.rules
	l.rulesLock.Unlock()

	target, err := l.mapper.Map(rules, req.conflictRec)
	if err != nil {
		l.logger.Errorf("Error mapping %s with rules=%s", req.conflictRec.SmallString(), rules)
		return
	}

	req.conflictRec.PopulateData(l.replId)

	// Log to the goxdcr.log as POC, instead of logging/send over network to a conflict bucket
	l.logger.Infof("Logging conflict %s, to target %s, sourceDoc=%s, targetDoc=%s",
		req.conflictRec, target, req.conflictRec.Source.GetDocBody(), req.conflictRec.Target.GetDocBody())

	return
}

type logReqHandleWithTimeout struct {
	ackCh chan error
}

func (h logReqHandleWithTimeout) Wait(finch chan bool) (err error) {
	timer := time.NewTicker(defaultWaitTimeout)

	select {
	case <-finch:
		err = ErrLogWaitAborted
	case <-timer.C:
		err = fmt.Errorf("timed out")
	case err = <-h.ackCh:
	}

	return
}
