package main

import (
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"time"
)

type Miner struct {
	validShares   uint64
	staleShares   uint64
	invalidShares uint64

	exit bool

	pool    *Stratum
	workers []Worker
}

func NewMiner(server string, algo *Algorithm, usrname string, password string, thread int) (*Miner, error) {
	m := &Miner{}

	p, err := NewStratum(server, algo, usrname, password)
	if err != nil {
		return nil, err
	}
	m.pool = p

	if thread <= 0 {
		thread = 1
	}
	m.workers = make([]Worker, thread)
	for _, w := range m.workers {
		worker := w
		worker.result = make(chan *JobResult, 1024)
		go func() {
			defer common.CrashLog()
			worker.start()
		}()
	}

	return m, nil
}

func (m *Miner) Stop() {
	m.exit = true
}

func (m *Miner) Run() {
	for !m.exit {
		needSleep := true
		if m.pool.job != nil {
			j := m.pool.job
			m.pool.job = nil
			seq := addNonceSequence()
			non := &Nonce{}
			for _, w := range m.workers {
				w.setJob(j, seq, non)
			}
			loggo.Info("Miner setJob ok id=%v algo=%v height=%v target=%v diff=%v", j.id, j.algorithm.name(), j.height, j.target, j.diff)
			needSleep = false
		}

		if needSleep {
			time.Sleep(time.Millisecond * 5)
		}
	}
}
