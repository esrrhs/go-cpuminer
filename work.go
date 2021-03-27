package main

import (
	"math/big"
)

// These are the locations of various data inside Work.Data.
const (
	TimestampWord = 2
	Nonce0Word    = 3
	Nonce1Word    = 4
	Nonce2Word    = 5
)

// NewWork is the constructor for Work.
func NewWork(data [192]byte, target *big.Int, jobTime uint32, timeReceived uint32,
	isGetWork bool) *Work {
	return &Work{
		Data:         data,
		Target:       target,
		JobTime:      jobTime,
		TimeReceived: timeReceived,
		IsGetWork:    isGetWork,
	}
}

// Work holds the data returned from getwork and if needed some stratum related
// values.
type Work struct {
	Data         [192]byte
	Target       *big.Int
	JobTime      uint32
	TimeReceived uint32
	IsGetWork    bool
}
