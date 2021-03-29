package main

const (
	kMaxBlobSize = 408
)

type Job struct {
	algorithm  *Algorithm
	nicehash   bool
	seed       []byte
	size       int
	clientId   string
	extraNonce string
	id         string
	poolWallet string
	backend    uint
	diff       uint64
	height     uint64
	target     uint64
	blob       [kMaxBlobSize]byte
	index      int
}

func (j *Job) setBlob(blob string) bool {
	// TODO
	return true
}

func (j *Job) setTarget(target string) bool {
	// TODO
	return true
}

func (j *Job) setSeedHash(hash string) bool {
	// TODO
	return true
}
