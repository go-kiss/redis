# Strings
- [Redis中文官网-Strings介绍](http://www.redis.cn/commands.html#string)
- [Redis官网-Strings介绍](https://redis.io/commands#string)
## 命令列表 
- [ ] [APPEND](#APPEND)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [BITCOUNT](#BITCOUNT)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [BITFIELD](#BITFIELD)
- [ ] [BITOP](#BITOP)
- [ ] [BITOS](#BITOS)
- [ ] [DECR](#DECR)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [DECRBY](#DECRBY)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [GET](#GET)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [GETBIT](#GETBIT)
- [ ] [GETRANGE](#GETRANGE)
- [ ] [GETSET](#GETSET)
- [ ] [INCR](#INCR)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [INCRBY](#INCRBY)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [INCRBYFLOAT](#INCRBYFLOAT)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [MGET](#MGET)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [MSET](#MSET)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [MSETNX](#MSETNX)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [PSETEX](#PSETEX)
- [ ] [SET](#SET)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [SETBIT](#SETBIT)
- [ ] [SETEX](#SETEX)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [SETNX](#SETNX)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [SETRANGE](#SETRANGE)
- [ ] [STRLEN](#STRLEN)
    - [x] TestCase
    - [ ] BenchMark
## <span id="APPEND">APPEND</span>
```go
strLen, err := client.Append(ctx, "redis", "bilibili")
// len(bilibili) = 8
fmt.Println(strLen)
strLen, _ = client.Append(ctx, "redis", "1234")
// len(bilibili1234) = 12
fmt.Println(strLen)
```
## <span id="BITCOUNT">BITCOUNT</span>
```go
strLen, err := client.BitCount(ctx, "redis")
// or
strLen, err := client.BitCount(ctx, "redis", 0, 2)
```
## <span id="BITFIELD">BITFIELD</span>
## <span id="BITOP">BITOP</span>
## <span id="BITOS">BITOS</span>
## <span id="DECR">DECR</span>
```go
// set redis 2017
newValue, err := client.Decr(ctx, "redis")
// 2016
fmt.Println(newValue)
```
## <span id="DECRBY">DECRBY</span>
```go
// set redis 2017
newValue, err := client.DecrBy(ctx, "redis", 1002)
// 1015
fmt.Println(newValue)
```
## <span id="GET">GET</span>
```go
// get string
// set redis bilibili
item, err := client.Get(ctx, "redis")
// bilibili
fmt.Println(item.Value)
//---------------------------
// get int
// set redis bilibili
item, err := client.GetInt(ctx, "redis") // got an error
// set redis 2018
item, err := client.GetInt(ctx, "redis")
// 2018
fmt.Println(item.Value)
```
## <span id="GETBIT">GETBIT</span>
## <span id="GETRANGE">GETRANGE</span>
## <span id="GETSET">GETSET</span>
## <span id="INCR">INCR</span>
```go
// set redis 2017
newValue, err := client.Incr(ctx, "redis")
// 2018
fmt.Println(newValue)
```
## <span id="INCRBY">INCRBY</span>
```go
// set redis 2017
newValue, err := client.IncrBy(ctx, "redis", 1002)
// 3019
fmt.Println(newValue)
```
## <span id="INCRBYFLOAT">INCRBYFLOAT</span>
## <span id="MGET">MGET</span>
```go
items, err := client.MGet(ctx, "key1", "key2")
// map key1=value1 key2=value2
fmt.Println(items)
```
## <span id="MSET">MSET</span>
```go
err := client.MSet(ctx, "key1", "value1", "key2", "value2")
```
## <span id="MSETNX">MSETNX</span>
```go
err := client.MSetNX(ctx, "key1", "value1", "key2", "value2")
```
## <span id="PSETEX">PSETEX</span>
## <span id="SET">SET</span>
```go
// set xjj 21 
err := client.Set(ctx, &redis.Item{
            Key:   "xjj",
            Value: 21,
        })
```
## <span id="SETBIT">SETBIT</span>
## <span id="SETEX">SETEX</span>
```go
// setex xjj 3 21 
err := client.Set(ctx, &redis.Item{
            Key:   "xjj",
            Value: 21,
            TTL: 3 // after 3 second, the key will expire
        })
```
## <span id="SETNX">SETNX</span>
```go
// setnx xjj bilibili
err := client.Set(ctx, &redis.Item{
            Key:   "xjj",
            Value: "bilibili",
            Flags: 1,
        })
```
## <span id="SETRANGE">SETRANGE</span>
## <span id="STRLEN">STRLEN</span>
```go
// set redis bilibili
valueLen, err := client.StrLen(ctx, "redis")
// 8
fmt.Println(valueLen)
```




































































