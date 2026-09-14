package main

import (
	"time"

	"github.com/esrrhs/gohome/crypto"
	"github.com/esrrhs/gohome/loggo"
	"github.com/pkg/errors"
)

type Benchmark struct {
	exit  bool
	once  bool
	algos []*Algorithm
}

func NewBenchmark(algo string) (*Benchmark, error) {
	var algos []string
	if algo == "all" {
		algos = crypto.Algo()
	} else {
		algos = append(algos, algo)
	}

	b := &Benchmark{once: algo == "all"}

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
	input := make([]byte, 76)
	for i := range input {
		input[i] = byte(i)
	}

	cy := crypto.NewCrypto("")

	for !b.exit {
		for _, al := range b.algos {
			name := al.supportAlgoName()
			start := time.Now()
			n := 0
			for !b.exit {
				cy.Sum(input, name, 0)
				n++
				if time.Since(start) >= 8*time.Second {
					break
				}
			}
			elapse := time.Since(start).Seconds()
			if elapse <= 0 {
				elapse = 1
			}
			speed := float64(n) / elapse
			loggo.Info("Benchmark Algo=%v HashSpeed=%.2f/s", name, speed)
			if b.exit {
				break
			}
		}
		if b.once || b.exit {
			break
		}
	}
}
