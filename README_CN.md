# go-cpuminer

[<img src="https://img.shields.io/github/license/esrrhs/go-cpuminer">](https://github.com/esrrhs/go-cpuminer)
[<img src="https://img.shields.io/github/languages/top/esrrhs/go-cpuminer">](https://github.com/esrrhs/go-cpuminer)
[<img src="https://img.shields.io/github/v/release/esrrhs/go-cpuminer">](https://github.com/esrrhs/go-cpuminer/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/go-cpuminer/total">](https://github.com/esrrhs/go-cpuminer/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/go-cpuminer">](https://hub.docker.com/repository/docker/esrrhs/go-cpuminer)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/go-cpuminer/go.yml?branch=master">](https://github.com/esrrhs/go-cpuminer/actions)

[English](./README.md) | 中文

用纯 Go 实现的 CPU 挖矿工具。

## 特性

- 纯 Go 实现，可编译到 Go 支持的任意平台
- 算法：`cn/0`、`cn/1`、`cn/2`、`cn/r`、`cn/fast`、`cn/half`、`cn/xao`、`cn/rto`、`cn/rwz`、`cn/zls`、`cn/double`、`cn-lite/0`、`cn-lite/1`、`cn-heavy/0`、`cn-heavy/tube`、`cn-heavy/xhv`、`cn-pico`、`cn-pico/tlo`
- 支持 Stratum 2.0 协议

## 编译

```
go build
```

交叉编译各平台并打包：

```
./pack.sh
```

## 用法

| 参数 | 默认值 | 说明 |
| --- | --- | --- |
| `-v`、`-version` | | 打印版本、git commit、分支和编译时间 |
| `-type` | `miner` | `miner` / `benchmark` / `test` |
| `-algo` | `cn-heavy/xhv` | 算法名称 |
| `-server` | `pool.hashvault.pro:80` | 矿池地址 |
| `-user` | | 钱包 / 用户名 |
| `-pass` | `x` | 密码 |
| `-thread` | `1` | 工作线程数 |
| `-nolog` | `0` | `>0` 时不写日志文件 |
| `-noprint` | `0` | `>0` 时不输出到 stdout |
| `-loglevel` | `info` | 日志级别 |
| `-profile` | `0` | pprof HTTP 端口 |
| `-cpuprofile` | | CPU profile 文件 |
| `-memprofile` | | 内存 profile 文件 |

### Hashvault Haven 挖矿

```
./go-cpuminer -server pool.hashvault.pro:80 -user hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG -pass x -algo cn-heavy/xhv
```

### Herominers Haven 挖矿

```
./go-cpuminer -server hk.haven.herominers.com:10450 -user hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG -pass x -algo cn-heavy/xhv
```

### 性能测试 / 挖矿测试

```
./go-cpuminer -type benchmark -algo cn-heavy/xhv
./go-cpuminer -type benchmark -algo all
./go-cpuminer -type test -algo cn-heavy/xhv
```

### Docker

```
docker run --name go-cpuminer -d --restart=always esrrhs/go-cpuminer ./go-cpuminer -server pool.hashvault.pro:80 -user hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG -pass x -algo cn-heavy/xhv
```

## 性能

单线程哈希速度对比 [xmrig](https://github.com/xmrig/xmrig) 6.26.0。测试环境：Windows amd64，Intel Core i7-6700K @ 4.00GHz，均为 1 线程；xmrig 使用 `--no-huge-pages`。go-cpuminer 一列可用 `./go-cpuminer -type benchmark -algo all` 复现。

| 算法 | xmrig (H/s) | go-cpuminer (H/s) | go / xmrig |
| --- | ---: | ---: | ---: |
| `cn/0` | 41.9 | 15.5 | 37% |
| `cn/1` | 39.6 | 25.3 | 64% |
| `cn/2` | 47.4 | 31.3 | 66% |
| `cn/r` | 28.1 | 5.0 | 18% |
| `cn/fast` | 67.8 | 53.8 | 79% |
| `cn/half` | 76.7 | 54.7 | 71% |
| `cn/xao` | 19.2 | 17.5 | 91% |
| `cn/rto` | 39.2 | 31.6 | 81% |
| `cn/rwz` | 57.6 | 37.0 | 64% |
| `cn/zls` | 55.4 | 33.8 | 61% |
| `cn/double` | 20.0 | 11.1 | 56% |
| `cn-lite/0` | 102.2 | 75.8 | 74% |
| `cn-lite/1` | 119.4 | 82.9 | 69% |
| `cn-heavy/0` | 27.3 | 16.9 | 62% |
| `cn-heavy/tube` | 22.2 | 17.0 | 77% |
| `cn-heavy/xhv` | 21.1 | 17.8 | 84% |
| `cn-pico` | 612.3 | 327.9 | 54% |
| `cn-pico/tlo` | 628.5 | 345.7 | 55% |

go-cpuminer 整体慢于 xmrig，因为它是纯 Go 实现，而 xmrig 使用手写汇编。差距最大的是 `cn/r`。

## 参考

- https://github.com/sammy007/monero-stratum
- https://github.com/decred/gominer
- https://github.com/xmrig/xmrig
- https://github.com/gurupras/go-cryptonight-miner
- https://github.com/gurupras/go-stratum-client
- https://git.dero.io/DERO_Foundation/RandomX
- https://haven.herominers.com/#how-to-mine-haven-xhv
- https://github.com/tevador/RandomX
- https://medium.com/novamining/in-depth-analysis-on-stratum-protocol-and-its-known-vulnerabilities-3ef139495608
- https://zhuanlan.zhihu.com/p/34441197
- https://github.com/Equim-chan/cryptonight
- https://www.cs.cmu.edu/~dga/crypto/xmr/cryptonight.png
