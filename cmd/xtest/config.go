package main

import "github.com/couchbase/goxdcr/conflictlog"

type Config struct {
	// Name of the config to use
	Name      string `json:"name"`
	LogLevel  string `json:"logLevel"`
	DebugPort int    `json:"debugPort"`

	ConflictLogPertest *ConflictLogLoadTest `json:"conflictLogLoadTest"`
}

// ConflictLogLoadTest is the config for load testing of the
// conflict logging
type ConflictLogLoadTest struct {
	// Target is the target conflict bucket details
	Target conflictlog.Target `json:"target"`

	// DocSizeRange is the min and max size in bytes of the source & target documents
	// in a conflict
	DocSizeRange [2]int `json:"docSizeRange"`

	// ConnType is the underlying type of connection to use
	// Values: "gocbcore", "memcached"
	ConnType string `json:"connType"`

	// DocLoadCount is the count of conflicts to log
	// Note: the actual count of docs written will be much larger
	DocLoadCount int `json:"docLoadCount"`

	// Batch count is the count of conflicts are logged before
	// waiting for all of them
	BatchCount int `json:"batchCount"`

	// XMemCount simulates the XMemNozzle count which calls the Log()
	XMemCount int `json:"xmemCount"`

	// ----- Logger configuration -----
	// LogQueue is the logger's internal channel capacity
	LogQueue int `json:"logQueue"`

	// WorkerCount is the logger's spawned worker count
	WorkerCount int `json:"workerCount"`
}
