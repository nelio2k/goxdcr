package conflictlog

import (
	"time"

	"github.com/couchbase/goxdcr/base"
)

const (
	ConflictManagerLoggerName = "conflictMgr"

	ConflictLoggerName string = "conflictLogger"

	MemcachedConnUserAgent = "conflictLog"

	DefaultLogCapacity          = 5
	DefaultLoggerWorkerCount    = 3
	DefaultNetworkRetryCount    = 6
	DefaultNetworkRetryInterval = 10 * time.Second
	LoggerShutdownChCap         = 10

	// DefaultPoolGCInterval is the GC frequency for connection pool
	DefaultPoolGCInterval = 60 * time.Second
	// DefaultPoolReapInterval is the last used threshold for reaping unused connections
	DefaultPoolReapInterval = 120 * time.Second

	SourcePrefix    string = "src"
	TargetPrefix    string = "tgt"
	CRDPrefix       string = "crd"
	TimestampFormat string = "2006-01-02T15:04:05.000Z07:00" // YYYY-MM-DDThh:mm:ss:SSSZ format
	// max increase in document body size after adding `"_xdcr_conflict":true` xattr.
	MaxBodyIncrease int = 4 /* entire xattr section size */ +
		4 /* _xdcr_conflict xattr size */ +
		len(base.ConflictLoggingXattrKey) +
		len(base.ConflictLoggingXattrVal) + 2 /* null terminators one each after key and value */
)
