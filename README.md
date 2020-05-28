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


