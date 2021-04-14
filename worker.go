package main

import (
	"encoding/binary"
	"github.com/esrrhs/go-engine/src/crypto"
	"github.com/esrrhs/go-engine/src/loggo"
	"sync"
	"time"
)

const (
	kReserveCount = 32768
)

type Worker struct {
	wj     *WorkerJob
	count  uint64
	lock   sync.Locker
	result chan<- *JobResult
}

func (w *Worker) start() {
	for {
		if w.wj == nil {
			time.Sleep(time.Millisecond * 5)
			continue
		}

		for w.wj.seq == gSequence {
			job := w.wj.currentJob()
			currentJobNonces := w.wj.nonce()

			algo := job.algorithm.supportAlgoName()
			hash := crypto.Sum(w.wj.blob()[0:job.size], algo, job.height)

			if !w.nextRound() {
				break
			}

			value := binary.LittleEndian.Uint64(hash[24:])
			if value < job.target {
				w.submit(job, currentJobNonces, hash)
			}

			w.count++
		}

		if w.wj.seq == gSequence {
			w.lock.Lock()
			if w.wj.seq == gSequence {
				w.wj = nil
				loggo.Debug("remove job %v", gSequence)
			}
			w.lock.Unlock()
		}
	}
}

func (w *Worker) nextRound() bool {
	if !w.wj.nextRound(kReserveCount, 1) {
		w.done(w.wj.currentJob())
		return false
	}
	return true
}

func (w *Worker) done(job *Job) {
	loggo.Debug("job done %v", job.id)
}

func (w *Worker) submit(job *Job, nonces uint32, hash []byte) {
	loggo.Debug("job submit %v %v", job.id, nonces)

	jr := &JobResult{job, nonces, hash}
	w.result <- jr
}

func (w *Worker) setJob(j *Job, sequence uint64, non *Nonce) {
	wj := &WorkerJob{}
	wj.non = non
	wj.add(j, sequence, kReserveCount)
	w.lock.Lock()
	w.wj = wj
	loggo.Debug("add done %v", sequence)
	w.lock.Unlock()
}
