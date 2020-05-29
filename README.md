# redis

轻量 redis 客户端

支持传入 ctx 对象。纯单机版 SDK，集群相关内容请使用 envoy 等中间件。

## 更新日志
- 1.0.0


## BasicUsage
创建客户端
```go
import "redis"

client := redis.New(Options{
    Address:  "localhost:6379",
    PoolSize: 1,
})
```
## Command Group
- [Cluster](doc/Cluster.md)
- [Geo](doc/Geo.md)
- [Hashes](doc/Hashes.md)
- [HyperLogLog](doc/HyperLogLog.md)
- [Keys](doc/Keys.md)
- [Lists](doc/Lists.md)
- [Pub/Sub](doc/PubSub.md)
- [Scripting](doc/Scripting.md)
- [Server](doc/Server.md)
- [Sets](doc/Sets.md)
- [Sorted Sets](doc/SortedSets.md)
- [Streams](doc/Streams.md)
- [Strings](doc/Strings.md)
- [Transactions](doc/Transactions.md)


