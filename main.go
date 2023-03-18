package main

import (
	"flag"
	"github.com/esrrhs/gohome/common"
	"github.com/esrrhs/gohome/loggo"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"time"
)

type Runner interface {
	Stop()
	Run()
}

func main() {

	defer common.CrashLog()

	ty := flag.String("type", "miner", "miner/benchmark/test")
	algo := flag.String("algo", "cn-heavy/xhv", "algo name")
	username := flag.String("user", "hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG", "username")
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

	var r Runner
	if *ty == "benchmark" {
		b, err := NewBenchmark(*algo)
		if err != nil {
			loggo.Error("Error initializing Benchmark: %v", err)
			return
		}
		r = b
	} else if *ty == "test" {
		t, err := NewTester(*algo)
		if err != nil {
			loggo.Error("Error initializing Benchmark: %v", err)
			return
		}
		r = t
	} else if *ty == "miner" {
		m, err := NewMiner(*server, *algo, *username, *password, *thread)
		if err != nil {
			loggo.Error("Error initializing miner: %v", err)
			return
		}
		r = m
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer common.CrashLog()
		<-c
		loggo.Warn("Got Control+C, exiting...")
		r.Stop()
	}()

	r.Run()

	loggo.Info("exit...")
}
