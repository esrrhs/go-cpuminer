package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/esrrhs/go-engine/src/crypto"
	"github.com/esrrhs/go-engine/src/loggo"
	"sync/atomic"
	"time"
)

const (
	TEST_BLOB   = "1010f1c2eb830614c310e956721ef3e47e882d077c2d2f161a9a2ae5567419b4d7b54e40d5800a0000000070b340380900000020123822000000000000000000000000000000000000000000000000000000000000000000000000308f1744ee050000c06891b8c30000000000000000000000000000000000000000000000000000000000000000000000202708d10516000000c78abc7b1600000025b190ff16000000fda372b1170000931ecff705873a888568c0fb25a4a979ca5ca798f2f9994a6266ec197d175eee4398a8d498c4b875f8cbc9147cda00d98bbe60103a293aad30157bd2a51a4fd5b99d52ce015b5ee297c58f79d409ff2b793e08e05a1f94bf76b6013e120db18001"
	TEST_TARGET = "711b0d00"
	TEST_HEIGHT = 835301
)

type Tester struct {
	exit bool
	algo *Algorithm
}

func NewTester(alname string) (*Tester, error) {

	t := &Tester{}

	al := NewAlgorithm(alname)
	if al.id == INVALID {
		return nil, errors.New("Tester create algo fail " + alname)
	}
	if al.supportAlgoName() == "" {
		return nil, errors.New("Tester support algo fail " + alname)
	}
	if !crypto.TestSum(al.supportAlgoName()) {
		return nil, errors.New("Tester test algo fail " + al.supportAlgoName())
	}
	t.algo = al

	return t, nil
}

func (t *Tester) Stop() {
	t.exit = true
}

func (t Tester) Run() {

	job := &Job{}
	job.algorithm = t.algo
	job.setBlob(TEST_BLOB)
	job.setTarget(TEST_TARGET)
	job.height = TEST_HEIGHT

	non := &Nonce{}

	wj := WorkerJob{}
	wj.non = non
	wj.add(job, 1, kReserveCount)

	done := 0
	n := uint32(0)

	cy := crypto.NewCrypto("")

	start := time.Now()

	for !t.exit {
		job := wj.currentJob()
		currentJobNonces := wj.nonce0()

		algo := job.algorithm.supportAlgoName()
		hash := cy.Sum(wj.blob()[0:job.size], algo, job.height)

		if !wj.nextRound(kReserveCount, 1) {
			break
		}

		value := binary.LittleEndian.Uint64(hash[24:])
		if value < job.target {
			loggo.Warn("Tester find hash Algo=%v Nonce=%v Blob=%v Hash=%v Value=%v Target=%v",
				t.algo.supportAlgoName(), currentJobNonces, hex.EncodeToString(wj.blob()[0:job.size]),
				hex.EncodeToString(hash), value, job.target)
			done++
		}

		atomic.AddUint32(&n, 1)

		elapse := time.Now().Sub(start)
		if elapse > time.Minute {
			start = time.Now()
			speed := float32(n) / float32(elapse/time.Second)
			n = 0
			loggo.Info("Tester Algo=%v Nonces=%v HashSpeed=%v/s Done=%v", t.algo.supportAlgoName(), currentJobNonces,
				speed, done)
		}
	}
}
