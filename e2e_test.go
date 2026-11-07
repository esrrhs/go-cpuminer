package main

import (
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// TestMiningLoopHashRate 测量完整 worker 挖矿循环的哈希率：
// nonce 迭代 -> 任务 blob 切片 -> 哈希 -> 目标判断 -> 统计计数。
// 默认跳过，避免拖慢 CI；需要时用
//
//	CPUMINER_E2E=1 CPUMINER_E2E_SECONDS=8 CPUMINER_E2E_THREADS=1,8 \
//	  go test -run TestMiningLoopHashRate -v
//
// 运行。结果以 H/s 打印。CPUMINER_E2E_THREADS 不指定时默认 1,8。
func TestMiningLoopHashRate(t *testing.T) {
	if os.Getenv("CPUMINER_E2E") != "1" {
		t.Skip("set CPUMINER_E2E=1 to run the end-to-end mining loop benchmark")
	}
	secs := 8
	if v := os.Getenv("CPUMINER_E2E_SECONDS"); v != "" {
		fmt.Sscanf(v, "%d", &secs)
	}
	algos := []string{"cn/0", "cn-heavy/xhv"}
	if v := os.Getenv("CPUMINER_E2E_ALGOS"); v != "" {
		algos = splitComma(v)
	}
	threads := []int{1, 8}
	if v := os.Getenv("CPUMINER_E2E_THREADS"); v != "" {
		threads = nil
		for _, s := range splitComma(v) {
			n, err := strconv.Atoi(s)
			if err != nil {
				t.Fatalf("bad CPUMINER_E2E_THREADS %q", v)
			}
			threads = append(threads, n)
		}
	}
	reps := 1
	if v := os.Getenv("CPUMINER_E2E_REPS"); v != "" {
		fmt.Sscanf(v, "%d", &reps)
	}

	for _, algo := range algos {
		for _, n := range threads {
			var rates []float64
			for r := 0; r < reps; r++ {
				rates = append(rates, measureMiningLoop(t, algo, n, time.Duration(secs)*time.Second))
				t.Logf("%-14s %2dt rep%d %6.1f H/s", algo, n, r, rates[len(rates)-1])
				fmt.Printf("RESULT %s %dt rep%d %.1f H/s\n", algo, n, r, rates[len(rates)-1])
			}
			fmt.Printf("MEDIAN %s %dt %.1f H/s (min %.1f max %.1f)\n",
				algo, n, median(rates), minFloat(rates), maxFloat(rates))
		}
	}
}

func splitComma(s string) []string {
	var out []string
	cur := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			if cur != "" {
				out = append(out, cur)
			}
			cur = ""
			continue
		}
		cur += string(s[i])
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}

func median(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	s := append([]float64(nil), v...)
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
	return s[len(s)/2]
}

func minFloat(v []float64) float64 {
	m := v[0]
	for _, x := range v {
		if x < m {
			m = x
		}
	}
	return m
}

func maxFloat(v []float64) float64 {
	m := v[0]
	for _, x := range v {
		if x > m {
			m = x
		}
	}
	return m
}

func measureMiningLoop(t *testing.T, algo string, threads int, d time.Duration) float64 {
	t.Helper()

	job := &Job{algorithm: NewAlgorithm(algo), size: 76, height: 0, target: 0}
	for i := 0; i < job.size; i++ {
		job.blob[i] = byte(i*7 + 3)
	}

	seq := addNonceSequence()
	stat := &Stat{}
	result := make(chan *JobResult, 1024)

	var workers []*Worker
	for i := 0; i < threads; i++ {
		w := NewWorker(result, stat)
		workers = append(workers, w)
		w.setJob(job, seq, &Nonce{})
		go w.start()
	}

	before := atomic.LoadUint32(&stat.hash)
	time.Sleep(d)
	after := atomic.LoadUint32(&stat.hash)

	// 注意：worker.start() 在测量结束后会带着旧 seq 的 wj 在外层循环里无 sleep
	// 自旋（gSequence 已变，但它不会清 wj），因此同一个进程内的下一次测量会
	// 被这些自旋 goroutine 污染。此测试因此约定「一个进程只测一组配置」：
	// 多组配置请用 shell 循环分别启动，进程退出即可回收自旋 goroutine。
	atomic.AddUint64(&gSequence, 1)

	return float64(after-before) / d.Seconds()
}
