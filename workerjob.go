package main

import "encoding/binary"

type WorkerJob struct {
	non        *Nonce
	seq        uint64
	nonce_mask uint64
	rounds     uint32
	job        *Job
	blobs      [kMaxBlobSize]byte
}

func (wj *WorkerJob) currentJob() *Job {
	return wj.job
}

func (wj *WorkerJob) blob() []byte {
	return wj.blobs[:]
}

func (wj *WorkerJob) sequence() uint64 {
	return wj.seq
}

func (wj *WorkerJob) nonce0() uint32 {
	return binary.LittleEndian.Uint32(wj.blob()[wj.nonceOffset():])
}

func (wj *WorkerJob) nonce1() uint32 {
	return binary.LittleEndian.Uint32(wj.blob()[wj.nonceOffset()+4:])
}

func (wj *WorkerJob) setNonce0(n uint32) {
	binary.LittleEndian.PutUint32(wj.blob()[wj.nonceOffset():], n)
}

func (wj *WorkerJob) setNonce1(n uint32) {
	binary.LittleEndian.PutUint32(wj.blob()[wj.nonceOffset()+4:], n)
}

func (wj *WorkerJob) nonceOffset() int {
	return wj.currentJob().nonceOffset()
}

func (wj *WorkerJob) nonceSize() int {
	return wj.currentJob().nonceSize()
}

func (wj *WorkerJob) nonceMask() uint64 {
	return wj.nonce_mask
}

func (wj *WorkerJob) add(job *Job, sequence uint64, reserveCount uint32) {
	wj.seq = sequence
	size := job.size
	wj.job = job
	wj.rounds = 0
	wj.nonce_mask = job.nonceMask()
	copy(wj.blobs[:size], job.blob[:size])
	_, n0, n1 := wj.non.next(wj.nonce0(), wj.nonce1(), reserveCount, wj.nonceMask())
	wj.setNonce0(n0)
	wj.setNonce1(n1)
}

func (wj *WorkerJob) nextRound(rounds uint32, roundSize uint32) bool {
	wj.rounds++
	if (wj.rounds & (rounds - 1)) == 0 {
		b, n0, n1 := wj.non.next(wj.nonce0(), wj.nonce1(), rounds*roundSize, wj.nonceMask())
		if !b {
			return false
		}
		wj.setNonce0(n0)
		wj.setNonce1(n1)
	} else {
		n := wj.nonce0() + roundSize
		wj.setNonce0(n)
	}
	return true
}
