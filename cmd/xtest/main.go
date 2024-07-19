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

func main() {
	addrGetter := &MemAddrGetter{}
	utils := utils.NewUtilities()
	//log.DefaultLoggerContext.SetLogLevel(log.LogLevelDebug)
	conflictlog.InitManager(log.DefaultLoggerContext, utils, addrGetter)

	err := conflictLogTest()
	if err != nil {
		fmt.Printf("error=%v\n", err)
	}
}
