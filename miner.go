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
	workers []*Worker
	jobs    chan *Job
	result  chan *JobResult

	stat *Stat
}

func NewMiner(server string, algo *Algorithm, usrname string, password string, thread int) (*Miner, error) {
	m := &Miner{}

	m.jobs = make(chan *Job, 16)
	m.result = make(chan *JobResult, 1024)
	m.stat = &Stat{}

	p, err := NewStratum(server, algo, usrname, password, m.jobs, m.stat)
	if err != nil {
		return nil, err
	}
	m.pool = p

	if thread <= 0 {
		thread = 1
	}
	m.workers = make([]*Worker, thread)
	for i, _ := range m.workers {
		w := NewWorker(m.result, m.stat)
		m.workers[i] = w
		go func() {
			defer common.CrashLog()
			w.start()
		}()
	}

	go func() {
		defer common.CrashLog()
		m.dispatch()
	}()

	go func() {
		defer common.CrashLog()
		m.commit()
	}()

	return m, nil
}

func (m *Miner) Stop() {
	m.exit = true
}

func (m *Miner) Run() {
	for !m.exit {
		time.Sleep(time.Second * 5)
		m.pool.hb()
		loggo.Info("Hash=%v, Job=%v, JobSubmit=%v, JobAccept=%v, JobFail=%v", m.stat.hash, m.stat.job,
			m.stat.submitJob, m.stat.submitJobOK, m.stat.submitJobFail)
		m.stat.clear()
	}
}

func (m *Miner) commit() {
	for {
		select {
		case data := <-m.result:
			m.pool.submit(data)
		}
	}
}

func (m *Miner) dispatch() {
	for {
		select {
		case j := <-m.jobs:
			seq := addNonceSequence()
			non := &Nonce{}
			for _, w := range m.workers {
				w.setJob(j, seq, non)
			}
			loggo.Info("Miner setJob ok id=%v algo=%v height=%v target=%v diff=%v", j.id, j.algorithm.name(), j.height, j.target, j.diff)
		}
	}
}
