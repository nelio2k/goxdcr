package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/couchbase/goxdcr/base"
	"github.com/couchbase/goxdcr/conflictlog"
	"github.com/couchbase/goxdcr/log"
)

type TestStruct struct {
	Id  string
	Str string `json:"str"`
}

func genRandomJson(id string, min, max int) ([]byte, error) {
	n := min + rand.Intn(max-min)
	t := TestStruct{
		Id:  id,
		Str: RandomString[:n],
	}

	buf, err := json.Marshal(&t)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func createCRDDoc(prefix string, i int, min, max int) (crd *conflictlog.ConflictRecord, err error) {
	sourceDocId := fmt.Sprintf("%s-%d", prefix, i)
	sourceBuf, err := genRandomJson(sourceDocId, min, max)
	if err != nil {
		return
	}

	targetBuf, err := genRandomJson(sourceDocId, min, max)
	if err != nil {
		return
	}

	crd = &conflictlog.ConflictRecord{
		Source: conflictlog.DocInfo{
			Id:       sourceDocId,
			Body:     sourceBuf,
			Datatype: base.JSONDataType,
		},
		Target: conflictlog.DocInfo{
			Id:       sourceDocId,
			Body:     targetBuf,
			Datatype: base.JSONDataType,
		},
	}

	return
}

func conflictLogLoadTest(cfg Config) (err error) {
	opts := cfg.ConflictLogPertest

	logger := log.NewLogger("xtest", log.DefaultLoggerContext)
	m, err := conflictlog.GetManager()
	if err != nil {
		return
	}

	m.SetConnType(opts.ConnType)

	target := conflictlog.NewTarget("B1", "S1", "C1")

	clog, err := m.NewLogger(logger, "1234",
		conflictlog.WithMapper(conflictlog.NewFixedMapper(logger, target)),
		conflictlog.WithCapacity(opts.LogQueue),
		conflictlog.WithWorkerCount(opts.WorkerCount),
	)
	if err != nil {
		return
	}

	start := time.Now()

	wg := &sync.WaitGroup{}
	docCountPerXmem := opts.DocLoadCount / opts.XMemCount
	finch := make(chan bool, 1)

	logger.Infof("xmemCount=%d, docCountPerXmem=%d", opts.XMemCount, docCountPerXmem)

	for i := 0; i < opts.XMemCount; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			docIdPrefix := fmt.Sprintf("xmem-doc-%d", n)
			logger.Infof("starting xmem with keyPrefix=%s", docIdPrefix)

			handles := []base.ConflictLoggerHandle{}
			for i := 0; i < docCountPerXmem; i++ {
				if i%opts.BatchCount == 0 {
					for _, h := range handles {
						err = h.Wait(finch)
						if err != nil {
							logger.Errorf("error in conflict log err=%v", err)
							return
						}
					}

					handles = []base.ConflictLoggerHandle{}
				}

				min, max := cfg.ConflictLogPertest.DocSizeRange[0], cfg.ConflictLogPertest.DocSizeRange[1]
				crd, err := createCRDDoc(docIdPrefix, i, min, max)
				if err != nil {
					logger.Errorf("error in creating crd doc err=%v", err)
					return
				}

				h, err := clog.Log(crd)
				if err != nil {
					logger.Errorf("error in sending conflict log err=%v", err)
					if err == conflictlog.ErrQueueFull {
						continue
					}
					return
				}

				handles = append(handles, h)
			}

			for _, h := range handles {
				err = h.Wait(finch)
				if err != nil {
					logger.Errorf("error in conflict log err=%v", err)
					return
				}
			}
		}(i)
	}

	wg.Wait()

	end := time.Now()

	fmt.Printf("Finished in %s\n", end.Sub(start))

	return nil
}
