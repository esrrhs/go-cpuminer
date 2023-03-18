# Go-cpuminer

[<img src="https://img.shields.io/github/license/esrrhs/go-cpuminer">](https://github.com/esrrhs/go-cpuminer)
[<img src="https://img.shields.io/github/languages/top/esrrhs/go-cpuminer">](https://github.com/esrrhs/go-cpuminer)
[![Go Report Card](https://goreportcard.com/badge/github.com/esrrhs/go-cpuminer)](https://goreportcard.com/report/github.com/esrrhs/go-cpuminer)
[<img src="https://img.shields.io/github/v/release/esrrhs/go-cpuminer">](https://github.com/esrrhs/go-cpuminer/releases)
[<img src="https://img.shields.io/github/downloads/esrrhs/go-cpuminer/total">](https://github.com/esrrhs/go-cpuminer/releases)
[<img src="https://img.shields.io/docker/pulls/esrrhs/go-cpuminer">](https://hub.docker.com/repository/docker/esrrhs/go-cpuminer)
[<img src="https://img.shields.io/github/actions/workflow/status/esrrhs/go-cpuminer/go.yml?branch=master">](https://github.com/esrrhs/go-cpuminer/actions)

go-cpuminer is a cpu miner tool implemented by pure go

# Feature
* Pure golang implementation, can support any platform
* Supported algorithms: CN / 0, CN / 1, CN / 2, CN / R, CN / FAST, CN / HALF, CN / XAO, CN / RTO, CN / RWZ, CN / DOUBLE, CN-Lite / 0 , CN-Lite / 1, CN-Heavy / 0, CN-Heavy / Tube, CN-Heavy / XHV, CN-Pico, CN-Pico / TLO
* Support STRATUM 2.0 protocol

# Compilation
```
# go build
```

# Example
* Hashvault Haven mining
```
./go-cpuminer -server pool.hashvault.pro:80 -user hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG -pass x -algo cn-heavy/xhv
```
* Herominers Haven mining
```
./go-cpuminer -server hk.haven.herominers.com:10450 -user hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG -pass x -algo cn-heavy/xhv
```
* HAVEN performance test
```
./go-cpuminer -type benchmark -algo cn-heavy/xhv
```
* Haven excavation test
```
./go-cpuminer -type test -algo cn-heavy/xhv
```
* use docker
```
docker run --name go-cpuminer -d --restart=always esrrhs/go-cpuminer ./go-cpuminer -server pool.hashvault.pro:80 -user hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG -pass x -algo cn-heavy/xhv
```

# Performance
Comparison with XMRIG's hash speed

|    platform    | xmrig     | go-cpuminer   |
| ------ | -------- | -------- |
| linux amd64 | 31H/s | 18H/s |

# Reference
* https://github.com/sammy007/monero-stratum
* https://github.com/decred/gominer
* https://github.com/xmrig/xmrig
* https://github.com/gurupras/go-cryptonight-miner
* https://github.com/gurupras/go-stratum-client
* https://git.dero.io/DERO_Foundation/RandomX
* https://haven.herominers.com/#how-to-mine-haven-xhv
* https://github.com/tevador/RandomX
* https://medium.com/novamining/in-depth-analysis-on-stratum-protocol-and-its-known-vulnerabilities-3ef139495608
* https://zhuanlan.zhihu.com/p/34441197
* https://github.com/Equim-chan/cryptonight
* https://www.cs.cmu.edu/~dga/crypto/xmr/cryptonight.png