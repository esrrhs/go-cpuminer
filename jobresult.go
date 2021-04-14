package main

type JobResult struct {
	job    *Job
	nonces uint32
	hash   []byte
}
