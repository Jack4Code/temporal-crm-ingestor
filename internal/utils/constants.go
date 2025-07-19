package utils

import "time"

const (
	// Expected header key for auth token
	AuthHeaderKey = "X-Auth-Token"

	// Default server port
	DefaultPort = ":8080"
)

var (
	// Default timeouts and retry intervals for Temporal activities
	DefaultActivityTimeout = time.Minute
	InitialRetryInterval   = time.Second * 2
	MaxRetryInterval       = time.Second * 30
)
