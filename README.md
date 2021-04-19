# go-cpuminer
go-cpuminer是go实现的cpu挖矿工具

# 特性
* 纯golang实现，可支持任意平台
* 支持的算法：cn/0，cn/1，cn/2，cn/r，cn/fast，cn/half，cn/xao，cn/rto，cn/rwz，cn/double，cn-lite/0，cn-lite/1，cn-heavy/0，cn-heavy/tube，cn-heavy/xhv，cn-pico，cn-pico/tlo
* 支持stratum 2.0协议

# 编译
```
# go build
```

# 示例
* hashvault haven挖矿
```
./go-cpuminer -server pool.hashvault.pro:80 -user hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG -pass x -algo cn-heavy/xhv
```
* herominers haven挖矿
```
./go-cpuminer -server hk.haven.herominers.com:10450 -user hvxxwtgSqXaH9AZYYed9NbijK8hydEVtpb2k8SLv39ZrQxHacwP8QeeYriNunavkRf5fYbdf6BPj6g7yGmh2kS2i4toHRp4pdG -pass x -algo cn-heavy/xhv
```
* haven性能测试
```
./go-cpuminer -type benchmark -algo cn-heavy/xhv
```
* haven挖矿测试
```
./go-cpuminer -type test -algo cn-heavy/xhv
```

# 参考
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
