package main

import "sync/atomic"

type Nonce struct {
	nonces [2]uint64
}

var gSequence uint64

func addNonceSequence() uint64 {
	return atomic.AddUint64(&gSequence, 1)
}

func (n *Nonce) next(nonce uint32, count uint32, mask uint64) (bool, uint32) {
	return false, 0
}
