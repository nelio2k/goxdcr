package conflictlog

import (
	"io"
	"testing"
	"time"

	"github.com/couchbase/goxdcr/log"
	"github.com/stretchr/testify/require"
)

type testConn struct {
	count *int
}

func (conn *testConn) Close() error {
	(*conn.count)++
	return nil
}

func TestPool(t *testing.T) {
	logger := log.NewLogger("testlogger", log.DefaultLoggerContext)

	buckets := map[string]*int{
		"B1": new(int),
		"B2": new(int),
	}

	newConnFn := func(bucketName string) (io.Closer, error) {
		count := buckets[bucketName]
		return &testConn{
			count: count,
		}, nil
	}

	pool := newConnPool(logger, newConnFn)
	pool.UpdateGCInterval(2 * time.Second)
	pool.UpdateReapInterval(5 * time.Second)

	connCount := 10
	bucket := "B1"
	connList := []io.Closer{}
	for i := 0; i < connCount; i++ {
		conn, err := pool.Get(bucket)
		require.Nil(t, err)
		connList = append(connList, conn)
	}

	for _, conn := range connList {
		pool.Put(bucket, conn)
	}

	time.Sleep(7 * time.Second)
	require.Equal(t, connCount, *(buckets[bucket]))
}
