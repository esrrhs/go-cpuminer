package main

import "time"

type JobResult struct {
	job    *Job
	nonce  uint32
	hash   [32]byte
	submit time.Time
}
