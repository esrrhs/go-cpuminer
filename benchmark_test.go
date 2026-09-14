package main

import (
	"testing"
	"time"
)

func TestNewBenchmarkSingleAlgo(t *testing.T) {
	b, err := NewBenchmark("cn/0")
	if err != nil {
		t.Fatal(err)
	}
	if b.once || len(b.algos) != 1 {
		t.Fatalf("once=%v algos=%d", b.once, len(b.algos))
	}
}

func TestBenchmarkOneHash(t *testing.T) {
	b, err := NewBenchmark("cn-lite/0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		time.Sleep(200 * time.Millisecond)
		b.Stop()
		close(done)
	}()
	b.Run()
	<-done
}
