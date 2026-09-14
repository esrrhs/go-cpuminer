package main

import "testing"

func TestWorkerJobRounds(t *testing.T) {
	job := &Job{algorithm: NewAlgorithm("cn/0")}
	if !job.setBlob(TEST_BLOB) {
		t.Fatal("blob")
	}
	wj := &WorkerJob{non: &Nonce{}}
	wj.add(job, 1, 4)
	if wj.currentJob() != job {
		t.Fatal("current job")
	}
	if wj.sequence() != 1 {
		t.Fatalf("seq = %d", wj.sequence())
	}
	first := wj.nonce0()
	if !wj.nextRound(4, 1) {
		t.Fatal("nextRound")
	}
	if wj.nonce0() != first+1 {
		t.Fatalf("nonce = %d, want %d", wj.nonce0(), first+1)
	}
	if len(wj.blob()) != kMaxBlobSize {
		t.Fatalf("blob buffer = %d", len(wj.blob()))
	}
}
