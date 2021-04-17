package main

import "sync/atomic"

type Nonce struct {
	nonces uint64
}

var gSequence uint64

func addNonceSequence() uint64 {
	return atomic.AddUint64(&gSequence, 1)
}

func (n *Nonce) next(nonce0 uint32, nonce1 uint32, reserveCount uint32, mask uint64) (bool, uint32, uint32) {
	mask &= 0x7FFFFFFFFFFFFFFF
	if reserveCount == 0 || mask < uint64(reserveCount)-1 {
		return false, 0, 0
	}

	counter := atomic.AddUint64(&n.nonces, uint64(reserveCount)) - uint64(reserveCount)
	for {
		if mask < counter {
			return false, 0, 0
		} else if mask-counter <= uint64(reserveCount)-1 {
			if mask-counter < uint64(reserveCount)-1 {
				return false, 0, 0
			}
		} else if uint32(0xFFFFFFFF)-(uint32)(counter) < reserveCount-1 {
			counter = atomic.AddUint64(&n.nonces, uint64(reserveCount)) - uint64(reserveCount)
			continue
		}
		ret_nonce0 := uint32((uint64(nonce0) & ^mask) | counter)
		ret_nonce1 := nonce1
		if mask > 0xFFFFFFFF {
			ret_nonce1 = uint32((uint64(nonce1) & (^mask >> 32)) | (counter >> 32))
		}
		return true, ret_nonce0, ret_nonce1
	}
}
