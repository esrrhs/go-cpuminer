package main

import (
	"github.com/esrrhs/gohome/crypto"
	"github.com/esrrhs/gohome/loggo"
	"github.com/pkg/errors"
	"time"
)

type Benchmark struct {
	exit  bool
	algos []*Algorithm
}

func NewBenchmark(algo string) (*Benchmark, error) {
	var algos []string
	if algo == "all" {
		algos = crypto.Algo()
	} else {
		algos = append(algos, algo)
	}

	b := &Benchmark{}

	for _, alname := range algos {
		al := NewAlgorithm(alname)
		if al.id == INVALID {
			return nil, errors.New("Benchmark create algo fail " + alname)
		}
		if al.supportAlgoName() == "" {
			return nil, errors.New("Benchmark support algo fail " + alname)
		}
		if !crypto.TestSum(al.supportAlgoName()) {
			return nil, errors.New("Benchmark test algo fail " + al.supportAlgoName())
		}
		b.algos = append(b.algos, al)
	}

	return b, nil
}

func (b *Benchmark) Stop() {
	b.exit = true
}

func (b *Benchmark) Run() {
	var input [kMaxBlobSize]byte
	for i, _ := range input {
		input[i] = byte(i)
	}

	cy := crypto.NewCrypto("")

	for !b.exit {
		for _, al := range b.algos {
			start := time.Now()
			n := 0
			for i := 0; i < 1024 && !b.exit; i++ {
				cy.Sum(input[:], al.supportAlgoName(), 0)
				n++
				if time.Now().Sub(start) > time.Second*5 {
					break
				}
			}
			elapse := time.Now().Sub(start)
			speed := float32(n) / float32(elapse/time.Second)
			loggo.Info("Benchmark Algo=%v HashSpeed=%v/s", al.supportAlgoName(), speed)
			if b.exit {
				break
			}
		}
	}
}
