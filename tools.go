package main

import (
	"encoding/hex"
	"github.com/esrrhs/go-engine/src/loggo"
)

func fromHexWithBuffer(data []byte, s string) bool {
	d, err := hex.DecodeString(s)
	if err != nil {
		loggo.Error("fromHex fail %v", err)
		return false
	}
	if len(d) > len(data) {
		loggo.Error("fromHex not enough %v %v %v", s, len(d), len(data))
		return false
	}
	copy(data, d)
	return true
}

func fromHex(s string) (bool, []byte) {
	data := make([]byte, len(s)/2)
	ret := fromHexWithBuffer(data[:], s)
	if !ret {
		return false, nil
	}
	return true, data
}

func toDiff(target uint64) uint64 {
	if target != 0 {
		return 0xFFFFFFFFFFFFFFFF / target
	}
	return 0
}
