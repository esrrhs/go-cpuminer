package main

import "testing"

func TestFromHexToHex(t *testing.T) {
	ok, data := fromHex("0a0b0c")
	if !ok {
		t.Fatal("fromHex failed")
	}
	if len(data) != 3 || data[0] != 0x0a || data[1] != 0x0b || data[2] != 0x0c {
		t.Fatalf("data = %v", data)
	}
	ok, s := toHex(data)
	if !ok || s != "0a0b0c" {
		t.Fatalf("toHex = %q", s)
	}
	if ok, _ := fromHex("zz"); ok {
		t.Fatal("invalid hex should fail")
	}
}

func TestFromHexWithBufferTooSmall(t *testing.T) {
	buf := make([]byte, 1)
	if fromHexWithBuffer(buf, "0a0b") {
		t.Fatal("overflow should fail")
	}
}

func TestToDiff(t *testing.T) {
	if toDiff(0) != 0 {
		t.Fatal("zero target should yield 0")
	}
	if toDiff(1) != 0xFFFFFFFFFFFFFFFF {
		t.Fatalf("toDiff(1) = %d", toDiff(1))
	}
	if toDiff(2) != 0x7FFFFFFFFFFFFFFF {
		t.Fatalf("toDiff(2) = %d", toDiff(2))
	}
}
