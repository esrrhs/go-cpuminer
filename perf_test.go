package main

import (
	"strconv"
	"strings"
	"testing"

	"github.com/esrrhs/gohome/crypto"
)

func benchName(algo string, threads int) string {
	return strings.ReplaceAll(algo, "/", "-") + "/" + strconv.Itoa(threads) + "t"
}

// perfInput 与 benchmark.go 保持一致，用于可重复的哈希率基准。
func perfInput() []byte {
	input := make([]byte, 76)
	for i := range input {
		input[i] = byte(i)
	}
	return input
}

// BenchmarkSum 逐算法测量单线程哈希率，等价于 -type benchmark 的核心循环。
// 运行方式：go test -run '^$' -bench BenchmarkSum -benchmem
func BenchmarkSum(b *testing.B) {
	input := perfInput()
	for _, name := range supportedAlgoNames() {
		name := name
		b.Run(name, func(b *testing.B) {
			cy := crypto.NewCrypto("")
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cy.Sum(input, name, 0)
			}
		})
	}
}

// BenchmarkSumParallel 按 miner 的线程模型测量扩展性：每个 goroutine 独占一个
// Crypto 实例（各自持有 scratchpad），互不共享状态。
// 运行方式：go test -run '^$' -bench 'BenchmarkSumParallel/cn-heavy-xhv' -benchtime 8s
func BenchmarkSumParallel(b *testing.B) {
	input := perfInput()
	for _, name := range []string{"cn/0", "cn-heavy/xhv"} {
		for _, threads := range []int{1, 2, 4, 8} {
			name, threads := name, threads
			b.Run(benchName(name, threads), func(b *testing.B) {
				b.ReportAllocs()
				done := make(chan struct{})
				per := b.N / threads
				if per == 0 {
					per = 1
				}
				b.ResetTimer()
				for t := 0; t < threads; t++ {
					go func() {
						cy := crypto.NewCrypto("")
						for i := 0; i < per; i++ {
							cy.Sum(input, name, 0)
						}
						done <- struct{}{}
					}()
				}
				for t := 0; t < threads; t++ {
					<-done
				}
			})
		}
	}
}
