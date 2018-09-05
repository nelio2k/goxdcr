package utils

import (
	base "github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/log"
	"sync"
)

const numOfSizes = 15

type DataPoolIface interface {
	GetByteSlice(sizeRequested uint64) ([]byte, error)
	PutByteSlice(doneSlice []byte)
}

type DataPool struct {
	byteSlicePoolClasses [numOfSizes]uint64
	byteSlicePools       [numOfSizes]sync.Pool

	logger_utils *log.CommonLogger
}

func NewDataPool() *DataPool {
	newPool := &DataPool{
		logger_utils: log.NewLogger("DataPool", log.DefaultLoggerContext),
	}

	newPool.byteSlicePoolClasses = [numOfSizes]uint64{
		5 << 10,  // 5k
		10 << 10, // 10k
		20 << 10, // 20k... etc
		40 << 10,
		80 << 10,
		100 << 10,
		200 << 10,
		400 << 10,
		800 << 10,
		1 << 20, // 1MB
		2 << 20, // 2MB
		4 << 20,
		8 << 20,
		10 << 20,
		21 << 20, // max value size is 20MB
	}

	newPool.byteSlicePools = [numOfSizes]sync.Pool{
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[0]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[1]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[2]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[3]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[4]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[5]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[6]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[7]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[8]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[9]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[10]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[11]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[12]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[13]) }},
		{New: func() interface{} { return make([]byte, newPool.byteSlicePoolClasses[14]) }},
	}

	return newPool
}

func (p *DataPool) GetByteSlice(sizeRequested uint64) ([]byte, error) {
	for i := 0; i < numOfSizes; i++ {
		if sizeRequested <= p.byteSlicePoolClasses[i] {
			return p.byteSlicePools[i].Get().([]byte), nil
		}
	}
	return nil, base.ErrorSizeExceeded
}

func (p *DataPool) PutByteSlice(doneSlice []byte) {
	sliceCap := uint64(cap(doneSlice))
	for i, n := range p.byteSlicePoolClasses {
		if sliceCap == n {
			p.byteSlicePools[i].Put(doneSlice)
			return
		}
	}
	// Should we panic? Leaking memory?
	p.logger_utils.Warnf("Unable to return byte slice of capacity %v back to pool. Potential memory leak\n", sliceCap)
}
