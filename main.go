package main

import (
	"flag"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/crypto"
	"github.com/esrrhs/go-engine/src/loggo"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"time"
)

func main() {

	defer common.CrashLog()

	algo := flag.String("algo", "", "algo name")
	username := flag.String("user", "my", "username")
	password := flag.String("pass", "x", "password")
	server := flag.String("server", "pool.hashvault.pro:80", "pool server addr")
	thread := flag.Int("thread", 1, "thread num")

	nolog := flag.Int("nolog", 0, "write log file")
	noprint := flag.Int("noprint", 0, "print stdout")
	loglevel := flag.String("loglevel", "info", "log level")
	profile := flag.Int("profile", 0, "open profile")
	cpuprofile := flag.String("cpuprofile", "", "open cpuprofile")
	memprofile := flag.String("memprofile", "", "open memprofile")

	flag.Parse()

	level := loggo.LEVEL_INFO
	if loggo.NameToLevel(*loglevel) >= 0 {
		level = loggo.NameToLevel(*loglevel)
	}
	loggo.Ini(loggo.Config{
		Level:     level,
		Prefix:    "gocpuminer",
		MaxDay:    3,
		NoLogFile: *nolog > 0,
		NoPrint:   *noprint > 0,
	})
	loggo.Info("start...")

	var al *Algorithm
	if *algo != "" {
		al = NewAlgorithm(*algo)
		if al.id == INVALID {
			loggo.Error("Unable to create algo %v", *algo)
			return
		}
		if al.supportAlgoName() == "" {
			loggo.Error("Unable to support algo %v", *algo)
			return
		}
		if !crypto.TestSum(al.supportAlgoName()) {
			loggo.Error("test algo %v fail", al.supportAlgoName())
			return
		}
	}

	if *profile > 0 {
		go http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*profile), nil)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			loggo.Error("Unable to create cpu profile: %v", err)
			return
		}
		pprof.StartCPUProfile(f)
		defer f.Close()
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			loggo.Error("Unable to create cpu profile: %v", err)
			return
		}
		timer := time.NewTimer(time.Minute * 20) // 20 minutes
		go func() {
			defer common.CrashLog()
			<-timer.C
			pprof.WriteHeapProfile(f)
			f.Close()
		}()
	}

	m, err := NewMiner(*server, al, *username, *password, *thread)
	if err != nil {
		loggo.Error("Error initializing miner: %v", err)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer common.CrashLog()
		<-c
		loggo.Warn("Got Control+C, exiting...")
		m.Stop()
	}()

	m.Run()

	loggo.Info("exit...")
}
