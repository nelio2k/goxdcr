package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/couchbase/goxdcr/conflictlog"
	"github.com/couchbase/goxdcr/log"
	"github.com/couchbase/goxdcr/utils"

	"net/http"
	_ "net/http/pprof"
)

const ConnStr = "localhost:12000"
const Bucket = "B1"

var gKey = "abcd"

type MemAddrGetter struct {
}

func (m *MemAddrGetter) MyMemcachedAddr() (string, error) {
	return "127.0.0.1:12000", nil
}

func loadConfigFile(filepath string) (cfg Config, err error) {
	f, err := os.Open(filepath)
	if err != nil {
		return
	}

	d := json.NewDecoder(f)

	err = d.Decode(&cfg)
	if err != nil {
		return
	}

	return
}

func main() {

	//...

	if len(os.Args[1:]) < 1 {
		fmt.Println("missing config json file")
		os.Exit(1)
	}

	configFile := os.Args[1]
	cfg, err := loadConfigFile(configFile)
	if err != nil {
		fmt.Printf("failed to load config file err=%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("cfg=%v\n", cfg)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.DebugPort), nil)
		fmt.Printf("failed to launch debug prof server, err=%v\n", err)
	}()

	logLevel, err := log.LogLevelFromStr(cfg.LogLevel)
	if err != nil {
		fmt.Printf("error=%v\n", err)
		return
	}

	log.DefaultLoggerContext.SetLogLevel(logLevel)

	addrGetter := &MemAddrGetter{}
	utils := utils.NewUtilities()
	conflictlog.InitManager(log.DefaultLoggerContext, utils, addrGetter)

	switch cfg.Name {
	case "conflictLogLoadTest":
		err = conflictLogLoadTest(cfg)
		if err != nil {
			fmt.Printf("error=%v\n", err)
		}
	default:
		fmt.Println("error: unknown config name =", cfg.Name)
	}
}
