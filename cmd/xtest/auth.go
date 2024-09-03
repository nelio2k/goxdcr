package main

import (
	"fmt"

	"github.com/couchbase/cbauth"
)

type CBAuthTest struct {
	AddrList []string `json:"addrList"`
}

func cbauthTest(cfg Config) (err error) {
	for _, addr := range cfg.CBAuthTest.AddrList {
		user, passwd, err := cbauth.GetMemcachedServiceAuth(addr)
		if err != nil {
			return err
		}

		fmt.Printf("addr=%s user=[%s] password=[%s]\n", addr, user, passwd)
	}

	return
}
