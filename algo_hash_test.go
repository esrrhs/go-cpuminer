package main

import (
	"testing"

	"github.com/esrrhs/gohome/crypto"
)

func TestAllSupportedAlgosHashVectors(t *testing.T) {
	names := supportedAlgoNames()
	if len(names) == 0 {
		t.Fatal("no supported algorithms")
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			if !crypto.TestSum(name) {
				t.Fatalf("hash vectors failed for %s", name)
			}
		})
	}
}
