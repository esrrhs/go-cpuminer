package main

import (
	"encoding/binary"
	"testing"
)

func TestJobSetBlobAndTarget(t *testing.T) {
	job := &Job{algorithm: NewAlgorithm("cn-heavy/xhv")}
	if job.setBlob("") || job.setBlob("abc") {
		t.Fatal("invalid blob should fail")
	}
	if !job.setBlob(TEST_BLOB) {
		t.Fatal("TEST_BLOB should parse")
	}
	if job.size != len(TEST_BLOB)/2 {
		t.Fatalf("size = %d, want %d", job.size, len(TEST_BLOB)/2)
	}
	if job.nonceOffset() != 39 {
		t.Fatalf("nonceOffset = %d, want 39", job.nonceOffset())
	}
	if job.nonceSize() != 4 {
		t.Fatalf("nonceSize = %d, want 4", job.nonceSize())
	}

	if job.setTarget("") || job.setTarget("zz") {
		t.Fatal("invalid target should fail")
	}
	if !job.setTarget(TEST_TARGET) {
		t.Fatal("TEST_TARGET should parse")
	}
	if job.target == 0 || job.diff == 0 {
		t.Fatalf("target=%d diff=%d", job.target, job.diff)
	}
	if job.diff != toDiff(job.target) {
		t.Fatalf("diff = %d, want %d", job.diff, toDiff(job.target))
	}

	ok, raw := fromHex(TEST_TARGET)
	if !ok {
		t.Fatal("TEST_TARGET hex")
	}
	want := 0xFFFFFFFFFFFFFFFF / (0xFFFFFFFF / uint64(binary.LittleEndian.Uint32(raw)))
	if job.target != want {
		t.Fatalf("target = %d, want %d", job.target, want)
	}
}

func TestJobSetTarget64(t *testing.T) {
	job := &Job{algorithm: NewAlgorithm("cn/0")}
	if !job.setTarget("0100000000000000") {
		t.Fatal("8-byte target should parse")
	}
	if job.target != 1 {
		t.Fatalf("target = %d, want 1", job.target)
	}
}

func TestJobSetSeedHash(t *testing.T) {
	job := &Job{}
	if job.setSeedHash("") || job.setSeedHash("aa") || job.setSeedHash("zz") {
		t.Fatal("invalid seed hash should fail")
	}
	seed := "00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
	if !job.setSeedHash(seed) {
		t.Fatal("valid seed hash should parse")
	}
	if len(job.seed) != kMaxSeedSize {
		t.Fatalf("seed len = %d, want %d", len(job.seed), kMaxSeedSize)
	}
}

func TestJobKawpowNonceLayout(t *testing.T) {
	job := &Job{algorithm: NewAlgorithm("kawpow")}
	if job.algorithm.family() != KAWPOW {
		t.Fatalf("family = %d, want KAWPOW", job.algorithm.family())
	}
	if job.nonceOffset() != 32 {
		t.Fatalf("nonceOffset = %d, want 32", job.nonceOffset())
	}
	if job.nonceSize() != 8 {
		t.Fatalf("nonceSize = %d, want 8", job.nonceSize())
	}
}

func TestJobNicehashMask(t *testing.T) {
	job := &Job{algorithm: NewAlgorithm("cn/0"), nicehash: true}
	if job.nonceMask() != 0xFFFFFF {
		t.Fatalf("nicehash mask = %d", job.nonceMask())
	}
	job.nicehash = false
	if job.nonceMask() != 0xFFFFFFFF {
		t.Fatalf("cn mask = %d", job.nonceMask())
	}
}
