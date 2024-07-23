package main

type Config struct {
	// Name of the config to use
	Name          string `json:"name"`
	LogLevel      string `json:"logLevel"`
	DebugPort     int    `json:"debugPort"`
	MemcachedAddr string `json:"memcachedAddr"`

	ConflictLogPertest *ConflictLogLoadTest `json:"conflictLogLoadTest"`
	GocbcoreTest       *GocbcoreTest        `json:"gocbcoreTest"`
}
