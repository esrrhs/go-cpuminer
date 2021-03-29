package main

import (
	"flag"
	"github.com/esrrhs/go-engine/src/common"
	"github.com/esrrhs/go-engine/src/loggo"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"sync/atomic"
	"time"
)

func main() {

	defer common.CrashLog()

	algo := flag.String("algo", "", "algo name")
	name := flag.String("name", "g", "worker name")
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

	ms := make([]*Miner, *thread)
	for i := 0; i < *thread; i++ {
		m, err := NewMiner(*server, *algo, *username, *password, *name+strconv.Itoa(i))
		if err != nil {
			loggo.Error("Error initializing miner: %v", err)
			return
		}
		ms[i] = m
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer common.CrashLog()
		<-c
		loggo.Warn("Got Control+C, exiting...")
		for i := 0; i < *thread; i++ {
			ms[i].Stop()
		}
	}()

	num := int32(*thread)
	for i := 0; i < *thread; i++ {
		index := i
		go func() {
			defer common.CrashLog()
			defer atomic.AddInt32(&num, -1)
			ms[index].Run()
		}()
	}

	for num > 0 {
		time.Sleep(time.Second)
	}

	loggo.Info("exit...")
}
