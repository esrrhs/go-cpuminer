package main

import (
	"os"
	"testing"

	"github.com/esrrhs/gohome/loggo"
)

func TestMain(m *testing.M) {
	loggo.Ini(loggo.Config{
		Level:     99,
		Prefix:    "test",
		MaxDay:    1,
		NoLogFile: true,
		NoPrint:   true,
	})
	os.Exit(m.Run())
}
