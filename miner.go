package main

import (
	"time"
)

type Miner struct {
	// The following variables must only be used atomically.
	validShares   uint64
	staleShares   uint64
	invalidShares uint64

	started uint32
	exit    bool

	pool *Stratum
}

func NewMiner(server, usrname, password, name string) (*Miner, error) {
	m := &Miner{}

	m.started = uint32(time.Now().Unix())

	s, err := StratumConn(server, usrname, password)
	if err != nil {
		return nil, err
	}
	m.pool = s

	return m, nil
}

func (m *Miner) Stop() {
	m.exit = true
}

func (m *Miner) Run() {
	for !m.exit {
		time.Sleep(time.Second)
	}
}
