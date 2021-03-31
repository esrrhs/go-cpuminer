package main

import (
	"github.com/pkg/errors"
	"time"
)

type Miner struct {
	validShares   uint64
	staleShares   uint64
	invalidShares uint64

	exit bool

	pool *Stratum
}

func NewMiner(server, algo, usrname, password, name string) (*Miner, error) {
	m := &Miner{}

	a := NewAlgorithm(algo)
	if a.id == INVALID {
		return nil, errors.New("NewAlgorithm fail")
	}

	s, err := NewStratum(server, a, usrname, password)
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
