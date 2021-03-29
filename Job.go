package main

import (
	"encoding/binary"
)

const (
	kMaxBlobSize = 408
	kMaxSeedSize = 32
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

	if blob == "" {
		return false
	}

	j.size = len(blob)
	if j.size%2 != 0 {
		return false
	}

	j.size /= 2

	minSize := j.nonceOffset() + j.nonceSize()
	if j.size < minSize || j.size >= kMaxBlobSize {
		return false
	}

	if !fromHexWithBuffer(j.blob[:], blob) {
		return false
	}

	if j.nonce() != 0 && !j.nicehash {
		j.nicehash = true
	}

	return true
}

func (j *Job) setTarget(target string) bool {

	if target == "" {
		return false
	}

	ok, raw := fromHex(target)
	if !ok {
		return false
	}
	size := len(raw)

	if size == 4 {
		j.target = 0xFFFFFFFFFFFFFFFF / (0xFFFFFFFF / uint64(binary.LittleEndian.Uint32(raw)))
	} else if size == 8 {
		j.target = binary.LittleEndian.Uint64(raw)
	}

	j.diff = toDiff(j.target)

	return true
}

func (j *Job) setSeedHash(hash string) bool {

	if hash == "" || len(hash) != kMaxSeedSize*2 {
		return false
	}

	ok, seed := fromHex(hash)
	if !ok {
		return false
	}
	j.seed = seed

	return true
}

func (j *Job) nonceOffset() int {
	if j.algorithm.family() == KAWPOW {
		return 32
	}
	return 39
}

func (j *Job) nonceSize() int {
	if j.algorithm.family() == KAWPOW {
		return 8
	}
	return 4
}

func (j *Job) nonce() uint32 {
	return binary.LittleEndian.Uint32(j.blob[j.nonceOffset():])
}
