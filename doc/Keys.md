# Keys
- [Redis中文官网-Keys介绍](http://www.redis.cn/commands.html#generic)
- [Redis官网-Keys介绍](https://redis.io/commands#generic)
## 命令列表 
- [ ] [DEL](#DEL)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [DUMP](#DUMP)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [EXISTS](#EXISTS)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [EXPIRE](#EXPIRE)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [EXPIREAT](#EXPIREAT)
- [ ] [KEYS](#KEYS)
- [ ] [MIGRATE](#MIGRATE)
- [ ] [MOVE](#MOVE)
- [ ] [OBJECT](#OBJECT)
- [ ] [PERSIST](#PERSIST)
- [ ] [PEXPIRE](#PEXPIRE)
- [ ] [PEXPIREAT](#PEXPIREAT)
- [ ] [PTTL](#PTTL)
- [ ] [RANDOMKEY](#RANDOMKEY)
- [ ] [RENAME](#RENAME)
- [ ] [RENAMENX](#RENAMENX)
- [ ] [RESTORE](#RESTORE)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [SORT](#SORT)
- [ ] [TTL](#TTL)
    - [x] TestCase
    - [ ] BenchMark
- [ ] [TYPE](#TYPE)
- [ ] [WAIT](#WAIT)
- [ ] [SCAN](#SCAN)
## <span id="DEL">DEL</span>
```go
// set redis bilibili
item, err := client.Del(ctx, "redis")
// nil
fmt.Println(item)
```
## <span id="DUMP">DUMP</span>
```go
// set redis bilibili
content, err := client.Dump(ctx, "redis")
```
## <span id="EXISTS">EXISTS</span>
```go
// set redis bilibili
isExists, err := client.Exists(ctx, "redis")
// true
fmt.Println(isExists)
isExists, err := client.Exists(ctx, "NotExistsKey")
// false
fmt.Println(isExists)
```
## <span id="EXPIRE">EXPIRE</span>
```go
// set redis bilibili
err := client.Expire(ctx, "redis", 3)
// the key will be expired after 3 second
```
## <span id="EXPIREAT">EXPIREAT</span>
## <span id="KEYS">KEYS</span>
## <span id="MIGRATE">MIGRATE</span>
## <span id="MOVE">MOVE</span>
## <span id="OBJECT">OBJECT</span>
## <span id="PERSIST">PERSIST</span>
## <span id="PEXPIRE">PEXPIRE</span>
## <span id="PEXPIREAT">PEXPIREAT</span>
## <span id="PTTL">PTTL</span>
## <span id="RANDOMKEY">RANDOMKEY</span>
## <span id="RENAME">RENAME</span>
## <span id="RENAMENX">RENAMENX</span>
## <span id="RESTORE">RESTORE</span>
```go
// set redis bilibili
// dump redis
// restore content to kaixinbaba
err := client.Restore(ctx, "kaixinbaba", content)
item, err := client.Get(ctx, "kaixinbaba")
// bilibili
fmt.Println(item.Value)
//---------------------------------
// restore content to kaixinbaba with expireTime
err := client.Restore(ctx, "kaixinbaba", 3, content)
item, err := client.Get(ctx, "kaixinbaba")
// bilibili
fmt.Println(item.Value)
// the key kaixinbaba will be expired after 3 second
```
## <span id="SORT">SORT</span>
## <span id="TTL">TTL</span>
```go
// setex redis 3 bilibili
expireTime, err := client.TTL(ctx, "redis")
// 3
fmt.Println(expireTime)
time.sleep(1 * time.Second)
expireTime, err := client.TTL(ctx, "redis")
// 2
fmt.Println(expireTime)
```
## <span id="TYPE">TYPE</span>
## <span id="WAIT">WAIT</span>
## <span id="SCAN">SCAN</span>




































































