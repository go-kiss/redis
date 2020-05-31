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
	Value interface{}

	ZSetValues map[string]float64

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

func (c *Client) Stats() *pool.Stats {
	return c.pool.Stats()
}
