package main

import "testing"

func TestNonceNext(t *testing.T) {
	n := &Nonce{}
	ok, n0, n1 := n.next(0, 0, 10, 0xFFFFFFFF)
	if !ok {
		t.Fatal("first next should succeed")
	}
	if n0 != 0 || n1 != 0 {
		t.Fatalf("first nonce = %d/%d, want 0/0", n0, n1)
	}
	ok, n0, n1 = n.next(0, 0, 10, 0xFFFFFFFF)
	if !ok {
		t.Fatal("second next should succeed")
	}
	if n0 != 10 || n1 != 0 {
		t.Fatalf("second nonce = %d/%d, want 10/0", n0, n1)
	}
}

func TestNonceNextRejectsBadReserve(t *testing.T) {
	n := &Nonce{}
	if ok, _, _ := n.next(0, 0, 0, 0xFFFFFFFF); ok {
		t.Fatal("reserveCount 0 should fail")
	}
	if ok, _, _ := n.next(0, 0, 8, 3); ok {
		t.Fatal("mask smaller than reserveCount should fail")
	}
}

func TestAddNonceSequence(t *testing.T) {
	a := addNonceSequence()
	b := addNonceSequence()
	if b != a+1 {
		t.Fatalf("sequence %d -> %d", a, b)
	}
}
