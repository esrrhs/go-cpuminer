package main

import (
	"encoding/binary"
	"github.com/esrrhs/gohome/crypto"
	"github.com/esrrhs/gohome/loggo"
	"sync"
	"sync/atomic"
	"time"
)

const (
	kReserveCount = 32768
)

type Worker struct {
	wj     *WorkerJob
	lock   sync.Mutex
	result chan *JobResult
	stat   *Stat
}

func NewWorker(result chan *JobResult, stat *Stat) *Worker {
	w := &Worker{}
	w.result = result
	w.stat = stat
	return w
}

func (w *Worker) start() {

	cy := crypto.NewCrypto("")

	for {
		if w.wj == nil {
			time.Sleep(time.Millisecond * 5)
			continue
		}

		for w.wj.seq == gSequence {
			job := w.wj.currentJob()
			currentJobNonces := w.wj.nonce0()

			algo := job.algorithm.supportAlgoName()
			hash := cy.Sum(w.wj.blob()[0:job.size], algo, job.height)

			if !w.nextRound() {
				break
			}

			value := binary.LittleEndian.Uint64(hash[24:])
			if value < job.target {
				w.submit(job, currentJobNonces, hash)
			}

			atomic.AddUint32(&w.stat.hash, 1)
		}

		if w.wj.seq == gSequence {
			w.lock.Lock()
			if w.wj.seq == gSequence {
				w.wj = nil
				loggo.Debug("worker remove job %v", gSequence)
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
	loggo.Debug("worker job done %v", job.id)
}

func (w *Worker) submit(job *Job, nonces uint32, hash []byte) {
	loggo.Debug("worker job submit %v %v", job.id, nonces)
	jr := &JobResult{}
	jr.job = job
	jr.nonce = nonces
	copy(jr.hash[:], hash)
	w.result <- jr
}

func (w *Worker) setJob(j *Job, sequence uint64, non *Nonce) {
	wj := &WorkerJob{}
	wj.non = non
	wj.add(j, sequence, kReserveCount)
	w.lock.Lock()
	w.wj = wj
	w.lock.Unlock()
	loggo.Debug("worker add done %v", sequence)
}
