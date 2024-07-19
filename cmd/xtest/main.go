package main

import (
	"fmt"

	"github.com/couchbase/goxdcr/conflictlog"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/utils"
)

const ConnStr = "localhost:12000"
const Bucket = "B1"

var gKey = "abcd"

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type MemAddrGetter struct {
}

func (m *MemAddrGetter) MyMemcachedAddr() (string, error) {
	return "127.0.0.1:12000", nil
}

func conflictLogTest() (err error) {
	addrGetter := &MemAddrGetter{}
	utils := utils.NewUtilities()
	log.DefaultLoggerContext.SetLogLevel(log.LogLevelDebug)
	//logger := log.NewLogger("xtest", log.DefaultLoggerContext)

	conflictlog.InitManager(log.DefaultLoggerContext, utils, addrGetter)

	m, err := conflictlog.GetManager()
	if err != nil {
		return
	}

	m.SetConnType("memcached")

	/*
		target := conflictlog.NewTarget("B1", "S1", "C1")

		clog, err := m.NewLogger(logger, "1234", conflictlog.WithMapper(conflictlog.NewFixedMapper(logger, target)))
		if err != nil {
			return
		}

		srcDocJson, err := json.Marshal(&Person{
			Name: "Tony Stark",
			Age:  30,
		})
		if err != nil {
			return err
		}

		targetDocJson, err := json.Marshal(&Person{
			Name: "Bruce Wayne",
			Age:  32,
		})
		if err != nil {
			return err
		}

		crd := &conflictlog.ConflictRecord{
			Source: conflictlog.DocInfo{
				Id:       fmt.Sprintf("superhero-%d", 0),
				Body:     srcDocJson,
				Datatype: base.JSONDataType,
			},
			Target: conflictlog.DocInfo{
				Id:       fmt.Sprintf("superhero-%d", 0),
				Body:     targetDocJson,
				Datatype: base.JSONDataType,
			},
		}

		for i := 0; i < 10; i++ {
			h, err := clog.Log(crd)
			if err != nil {
				return err
			}

			err = h.Wait(make(chan bool, 1))
			if err != nil {
				return err
			}

			time.Sleep(1 * time.Second)
		}

	*/
	return nil
}

func main() {
	err := conflictLogTest()
	//err := gocbcoreTest()
	if err != nil {
		fmt.Printf("error=%v\n", err)
	}
}
