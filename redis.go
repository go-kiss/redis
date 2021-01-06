package redis

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/bilibili/net/pool"
	"github.com/bilibili/redis/protocol"
)

const (
	FlagNX = 1 << 0
	FlagXX = 1 << 1
	FlagCH = 1 << 2
)

type EvalReturn struct {
	val interface{}
}

func (e EvalReturn) Int64() (int64, error) {
	v, ok := e.val.(int64)
	if !ok {
		return 0, fmt.Errorf("redis: unexpected type=%T for int64", e.val)
	}
	return v, nil
}

func (e EvalReturn) String() (string, error) {
	v, ok := e.val.(string)
	if !ok {
		return "", fmt.Errorf("redis: unexpected type=%T for string", e.val)
	}
	return v, nil
}

func (e EvalReturn) Array() ([]interface{}, error) {
	v, ok := e.val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type=%T for Array", e.val)
	}
	return v, nil
}

func (e EvalReturn) Interface() interface{} {
	return e.val
}

type Options struct {
	Address      string
	PoolSize     int
	MinIdleConns int

	MaxConnAge  time.Duration
	PoolTimeout time.Duration
	IdleTimeout time.Duration

	IdleCheckFrequency time.Duration

	OnPreCmd  func(context.Context, []interface{}) context.Context
	OnPostCmd func(context.Context, error)
}

type Client struct {
	opts Options
	pool pool.Pooler
}

func New(opts Options) Client {
	poolOpts := pool.Options{
		PoolSize:           opts.PoolSize,
		MinIdleConns:       opts.MinIdleConns,
		MaxConnAge:         opts.MaxConnAge,
		PoolTimeout:        opts.PoolTimeout,
		IdleTimeout:        opts.IdleTimeout,
		IdleCheckFrequency: opts.IdleCheckFrequency,
	}

	poolOpts.Dialer = func(ctx context.Context) (pool.Closer, error) {
		d := net.Dialer{}
		conn, err := d.DialContext(ctx, "tcp", opts.Address)
		if err != nil {
			return nil, err
		}

		rw := redisConn{
			c: conn,
			r: protocol.NewReader(conn),
			w: protocol.NewWriter(conn),
		}

		return &rw, nil
	}

	return Client{pool: pool.New(poolOpts), opts: opts}
}

type redisConn struct {
	c net.Conn
	r *protocol.Reader
	w *protocol.Writer
}

func (rc *redisConn) Close() error {
	return rc.c.Close()
}

type ZSetValue struct {
	Member string
	Score  float64
}

type Item struct {
	// Key is the Item's key (250 bytes maximum).
	Key string

	// Value is the Item's value.
	Value []byte

	ZSetValues map[string]float64

	HashValues map[string]string

	// Flags 一些 redis 标记位，请参考 Flag 开头的常量定义
	Flags uint32

	// TTL 缓存时间，秒，0 表示不过期
	TTL int32
}

var noDeadline = time.Time{}

// PoolStats 返回连接池状态
func (c *Client) PoolStats() *pool.Stats {
	return c.pool.Stats()
}

func (c *Client) do(ctx context.Context, args []interface{}, fn func(conn *redisConn) error) error {
	if c.opts.OnPreCmd != nil {
		ctx = c.opts.OnPreCmd(ctx, args)
	}

	conn, err := c.pool.Get(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if IsBadConn(err, false) {
			c.pool.Remove(conn)
		} else {
			c.pool.Put(conn)
		}

	}()

	rc := conn.C.(*redisConn)

	if t, ok := ctx.Deadline(); ok {
		err = rc.c.SetDeadline(t)
	} else {
		err = rc.c.SetDeadline(noDeadline)
	}

	if err != nil {
		return err
	}

	// 此处赋值给 defer 函数用的，不要去掉
	err = fn(rc)

	if c.opts.OnPostCmd != nil {
		c.opts.OnPostCmd(ctx, err)
	}

	return err
}

func (c *Client) Get(ctx context.Context, key string) (item *Item, err error) {
	args := []interface{}{"get", key}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var b []byte
		if b, err = conn.r.ReadBytesReply(); err != nil {
			item = nil
			return err
		}

		item = &Item{Value: b}

		return nil
	})
	return
}

func (c *Client) MGet(ctx context.Context, keys []string) (items map[string]*Item, err error) {
	args := make([]interface{}, 0, len(keys)+1)

	args = append(args, "mget")
	for _, key := range keys {
		args = append(args, key)
	}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		l, err := conn.r.ReadArrayLenReply()
		if err != nil {
			return err
		}

		items = make(map[string]*Item, l)

		for i := 0; i < l; i++ {
			b, err := conn.r.ReadBytesReply()
			if err == Nil {
				continue
			}
			if err != nil {
				return err
			}

			key := keys[i]

			items[key] = &Item{Value: b}
		}

		return nil
	})
	return
}

func (c *Client) Eval(ctx context.Context, script string, keys []string, argvs ...interface{}) (result *EvalReturn, err error) {
	args := make([]interface{}, 0, len(keys)+len(argvs)+2)
	args = append(args, "eval", script, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}
	args = append(args, argvs...)

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		b, err := conn.r.ReadInterfaceReply()
		if err != nil {
			return err
		}

		result = &EvalReturn{
			val: b,
		}

		return nil
	})

	return
}

func (c *Client) Set(ctx context.Context, item *Item) error {
	args := make([]interface{}, 0, 6)
	args = append(args, "set", item.Key, item.Value)

	if item.TTL > 0 {
		args = append(args, "EX", item.TTL)
	}

	if item.Flags&FlagNX > 0 {
		args = append(args, "NX")
	} else if item.Flags&FlagXX > 0 {
		args = append(args, "XX")
	}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadStatusReply()
		if err != nil {
			return err
		}

		return nil
	})
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	args := make([]interface{}, 0, 1+len(keys))

	args = append(args, "del")
	for _, key := range keys {
		args = append(args, key)
	}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadIntReply()

		return err
	})
}

func (c *Client) IncrBy(ctx context.Context, key string, by int64) (i int64, err error) {
	args := []interface{}{"incrby", key, by}

	err = c.do(ctx, args, func(conn *redisConn) error {

		if err = conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		i, err = conn.r.ReadIntReply()

		return err
	})

	return
}

func (c *Client) DecrBy(ctx context.Context, key string, by int64) (int64, error) {
	return c.IncrBy(ctx, key, -by)
}

func (c *Client) Expire(ctx context.Context, key string, ttl int32) error {
	args := []interface{}{"expire", key, ttl}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadIntReply()

		return err
	})
}

func (c *Client) TTL(ctx context.Context, key string) (ttl int32, err error) {
	args := []interface{}{"ttl", key}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var i int64
		i, err = conn.r.ReadIntReply()
		if err != nil {
			return err
		}

		if i == -2 {
			err = Nil
			return err
		}

		ttl = int32(i)

		return err
	})

	return
}

func (c *Client) ZAdd(ctx context.Context, item *Item) (added int64, err error) {
	args := make([]interface{}, 0, 4+len(item.ZSetValues))
	args = append(args, "zadd", item.Key)

	if item.Flags&FlagNX > 0 {
		args = append(args, "NX")
	} else if item.Flags&FlagXX > 0 {
		args = append(args, "XX")
	}

	if item.Flags&FlagCH > 0 {
		args = append(args, "CH")
	}

	for member, score := range item.ZSetValues {
		args = append(args, score, member)
	}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		added, err = conn.r.ReadIntReply()
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (c *Client) ZIncrBy(ctx context.Context, key, member string, by float64) error {
	args := []interface{}{"zincrby", key, by, member}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadFloat()
		return err
	})
}

func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) (values []*ZSetValue, err error) {
	return c.zrange(ctx, "zrange", key, start, stop, 0, 0)
}

func (c *Client) ZRevRange(ctx context.Context, key string, start, stop int64) (values []*ZSetValue, err error) {
	return c.zrange(ctx, "zrevrange", key, start, stop, 0, 0)
}

func (c *Client) ZRangeByScore(ctx context.Context, key string, min, max float64, offset, count int64) (values []*ZSetValue, err error) {
	return c.zrange(ctx, "zrangebyscore", key, min, max, offset, count)
}

func (c *Client) ZRevRangeByScore(ctx context.Context, key string, max, min float64, offset, count int64) (values []*ZSetValue, err error) {
	return c.zrange(ctx, "zrevrangebyscore", key, max, min, offset, count)
}

func (c *Client) zrange(ctx context.Context, cmd, key string, start, stop interface{}, offset, count int64) (values []*ZSetValue, err error) {
	args := []interface{}{cmd, key, start, stop, "WITHSCORES"}
	if count > 0 {
		args = append(args, "LIMIT", offset, count)
	}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		l, err := conn.r.ReadArrayLenReply()
		if err != nil {
			return err
		}

		values = make([]*ZSetValue, 0, l)
		for i := 0; i < l/2; i++ {
			b, err := conn.r.ReadBytesReply()
			if err != nil {
				return err
			}
			f, err := conn.r.ReadFloat()
			if err != nil {
				return err
			}

			values = append(values, &ZSetValue{Member: string(b), Score: f})
		}

		return nil
	})

	return
}

func (c *Client) ZRank(ctx context.Context, key, member string) (rank int64, err error) {
	return c.zrank(ctx, "zrank", key, member)
}
func (c *Client) ZRevRank(ctx context.Context, key, member string) (rank int64, err error) {
	return c.zrank(ctx, "zrevrank", key, member)
}

func (c *Client) zrank(ctx context.Context, cmd, key, member string) (rank int64, err error) {
	args := []interface{}{cmd, key, member}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		rank, err = conn.r.ReadIntReply()
		return err
	})

	return
}

func (c *Client) ZScore(ctx context.Context, key, member string) (score float64, err error) {
	args := []interface{}{"zscore", key, member}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		score, err = conn.r.ReadFloat()
		return err
	})

	return
}

func (c *Client) ZCard(ctx context.Context, key string) (card int64, err error) {
	args := []interface{}{"zcard", key}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		card, err = conn.r.ReadIntReply()

		return err
	})
	return
}

func (c *Client) ZCount(ctx context.Context, key, min, max string) (i int64, err error) {
	args := []interface{}{"zcount", key, min, max}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		i, err = conn.r.ReadIntReply()

		return err
	})
	return
}

func (c *Client) ZRem(ctx context.Context, keys ...string) error {
	args := make([]interface{}, 0, 1+len(keys))

	args = append(args, "zrem")
	for _, key := range keys {
		args = append(args, key)
	}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadIntReply()

		return err
	})
}

func (c *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (i int64, err error) {
	args := []interface{}{"zremrangebyrank", key, start, stop}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		i, err = conn.r.ReadIntReply()

		return err
	})
	return
}

func (c *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) (i int64, err error) {
	args := []interface{}{"zremrangebyscore", key, min, max}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		i, err = conn.r.ReadIntReply()

		return err
	})
	return
}

func (c *Client) Stats() *pool.Stats {
	return c.pool.Stats()
}

func (c *Client) SAdd(ctx context.Context, key string, data ...[]byte) (err error) {
	args := []interface{}{"sadd", key}
	for _, d := range data {
		args = append(args, d)
	}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err = conn.r.ReadIntReply()

		return err
	})
	return
}

func (c *Client) SPop(ctx context.Context, key string, cnt int32) (data [][]byte, err error) {
	args := []interface{}{"spop", key, cnt}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		l, err := conn.r.ReadArrayLenReply()
		if err != nil {
			return err
		}

		data = make([][]byte, l)
		for i := 0; i < l; i++ {
			b, err := conn.r.ReadBytesReply()
			if err != nil {
				return err
			}
			data[i] = b
		}

		return nil
	})
	return
}

func (c *Client) SCard(ctx context.Context, key string) (card int64, err error) {
	args := []interface{}{"scard", key}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		card, err = conn.r.ReadIntReply()

		return err
	})
	return
}

func (c *Client) SIsMember(ctx context.Context, key string, data []byte) (result bool, err error) {
	args := []interface{}{"sismember", key, data}

	var resultInt int64
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		resultInt, err = conn.r.ReadIntReply()

		return err
	})
	if resultInt == 1 {
		result = true
	}
	return
}

func (c *Client) SRem(ctx context.Context, key string, data ...[]byte) (result int64, err error) {
	args := []interface{}{"srem", key}
	for _, d := range data {
		args = append(args, d)
	}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		result, err = conn.r.ReadIntReply()

		return err
	})
	return
}

func (c *Client) SMembers(ctx context.Context, key string) (items [][]byte, err error) {
	args := []interface{}{"smembers", key}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		l, err := conn.r.ReadArrayLenReply()
		if err != nil {
			return err
		}

		items = make([][]byte, l)
		for i := 0; i < l; i++ {
			b, err := conn.r.ReadBytesReply()
			if err != nil {
				return err
			}
			items[i] = b
		}

		return nil
	})
	return
}

// Hashes
func (c *Client) HSet(ctx context.Context, key string, hashes map[string]string) (added int64, err error) {
	args := []interface{}{"hset", key}

	for f, v := range hashes {
		args = append(args, f)
		args = append(args, v)
	}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		added, err = conn.r.ReadIntReply()
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func (c *Client) HGet(ctx context.Context, key string, field string) (item *Item, err error) {
	args := []interface{}{"hget", key, field}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var b []byte
		if b, err = conn.r.ReadBytesReply(); err != nil {
			item = nil
			return err
		}

		item = &Item{Value: b}

		return nil
	})
	return
}

func (c *Client) HGetAll(ctx context.Context, key string) (item *Item, err error) {
	args := []interface{}{"hgetall", key}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		l, err := conn.r.ReadArrayLenReply()
		if err != nil {
			return err
		}
		hashes := make(map[string] string, l / 2)
		for i := 0; i < l; i += 2 {
			field, err := conn.r.ReadBytesReply()
			if err != nil {
				return err
			}
			value, err := conn.r.ReadBytesReply()
			if err != nil {
				return err
			}
			hashes[string(field)] = string(value)
		}
		item = &Item{HashValues: hashes}
		return nil
	})
	return
}
