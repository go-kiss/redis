# Strings
- [Redis中文官网-Strings介绍](http://www.redis.cn/commands.html#string)
- [Redis官网-Strings介绍](https://redis.io/commands#string)
## 命令列表 
- [ ] [APPEND](#APPEND)
- [ ] [BITCOUNT](#BITCOUNT)
- [ ] [BITFIELD](#BITFIELD)
- [ ] [BITOP](#BITOP)
- [ ] [BITOS](#BITOS)
- [ ] [DECR](#DECR)
- [ ] [DECRBY](#DECRBY)
- [ ] [GET](#GET)
- [ ] [GETBIT](#GETBIT)
- [ ] [GETRANGE](#GETRANGE)
- [ ] [GETSET](#GETSET)
- [ ] [INCR](#INCR)
- [ ] [INCRBY](#INCRBY)
- [ ] [INCRBYFLOAT](#INCRBYFLOAT)
- [ ] [MGET](#MGET)
- [ ] [MSET](#MSET)
- [ ] [MSETNX](#MSETNX)
- [ ] [PSETEX](#PSETEX)
- [ ] [SET](#SET)
- [ ] [SETBIT](#SETBIT)
- [ ] [SETEX](#SETEX)
- [ ] [SETNX](#SETNX)
- [ ] [SETRANGE](#SETRANGE)
- [ ] [STRLEN](#STRLEN)
    - [x] TestCase
    - [ ] BenchMark
## <span id="APPEND">APPEND</span>
## <span id="BITCOUNT">BITCOUNT</span>
## <span id="BITFIELD">BITFIELD</span>
## <span id="BITOP">BITOP</span>
## <span id="BITOS">BITOS</span>
## <span id="DECR">DECR</span>
## <span id="DECRBY">DECRBY</span>
## <span id="GET">GET</span>
## <span id="GETBIT">GETBIT</span>
## <span id="GETRANGE">GETRANGE</span>
## <span id="GETSET">GETSET</span>
## <span id="INCR">INCR</span>
## <span id="INCRBY">INCRBY</span>
## <span id="INCRBYFLOAT">INCRBYFLOAT</span>
## <span id="MGET">MGET</span>
## <span id="MSET">MSET</span>
## <span id="MSETNX">MSETNX</span>
## <span id="PSETEX">PSETEX</span>
## <span id="SET">SET</span>
## <span id="SETBIT">SETBIT</span>
## <span id="SETEX">SETEX</span>
## <span id="SETNX">SETNX</span>
## <span id="SETRANGE">SETRANGE</span>
## <span id="STRLEN">STRLEN</span>
```go
// set redis bilibili
valueLen, err := client.StrLen(ctx, "redis")
// 8
fmt.Println(valueLen)
```




































































