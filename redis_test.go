package redis

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"
)

var ctx = context.Background()

func TestBasic(t *testing.T) {
	c := New(Options{
		Address:  os.Getenv("REDIS_HOST"),
		PoolSize: 1,
	})

	if err := c.Del(ctx, "foo"); err != nil {
		t.Fatal("start faild")
	}

	c.Set(ctx, &Item{Key: "foo", Value: []byte("hello")})

	item, _ := c.Get(ctx, "foo")
	if !bytes.Equal(item.Value, []byte("hello")) {
		t.Fatal("get foo failed")
	}

	// add
	if err := c.Set(ctx, &Item{Key: "foo", Value: []byte("123"), Flags: FlagNX}); err != Nil {
		t.Fatal("add foo failed")
	}

	// replace
	if err := c.Set(ctx, &Item{Key: "foo", Value: []byte("123"), Flags: FlagXX}); err == Nil {
		t.Fatal("replace foo failed")
	}

	item, _ = c.Get(ctx, "foo")
	if !bytes.Equal(item.Value, []byte("123")) {
		t.Fatal("replace foo failed")
	}

	c.Del(ctx, "foo")
	if _, err := c.Get(ctx, "foo"); err != Nil {
		t.Fatal("del foo failed")
	}

	if _, err := c.TTL(ctx, "foo"); err != Nil {
		t.Fatal("ttl foo failed")
	}

	c.IncrBy(ctx, "foo", 3)
	item, _ = c.Get(ctx, "foo")
	if !bytes.Equal(item.Value, []byte("3")) {
		t.Fatal("increment foo failed")
	}

	c.DecrBy(ctx, "foo", 4)
	item, _ = c.Get(ctx, "foo")
	if !bytes.Equal(item.Value, []byte("-1")) {
		t.Fatal("decrement foo failed")
	}

	if ttl, _ := c.TTL(ctx, "foo"); ttl != -1 {
		t.Fatal("ttl foo failed")
	}

	c.Expire(ctx, "foo", 10)
	if ttl, _ := c.TTL(ctx, "foo"); ttl != 10 {
		t.Fatal("ttl foo failed")
	}
	time.Sleep(1 * time.Second)
	if ttl, _ := c.TTL(ctx, "foo"); ttl != 9 {
		t.Fatal("ttl foo failed")
	}
}

func TestZSet(t *testing.T) {
	c := New(Options{
		Address:  os.Getenv("REDIS_HOST"),
		PoolSize: 1,
	})

	if err := c.Del(ctx, "foo"); err != nil {
		t.Fatal("start faild")
	}

	c.ZAdd(ctx, &Item{Key: "foo", ZSetValues: map[string]float64{"a": 1, "b": 2}})

	values, _ := c.ZRange(ctx, "foo", 0, -1)
	if values[0].Member != "a" ||
		values[0].Score != 1 ||
		values[1].Member != "b" ||
		values[1].Score != 2 {

		t.Fatal("zrange faild")
	}

	values, _ = c.ZRevRange(ctx, "foo", 0, -1)
	if values[1].Member != "a" ||
		values[1].Score != 1 ||
		values[0].Member != "b" ||
		values[0].Score != 2 {

		t.Fatal("zrange faild")
	}

	if c, _ := c.ZCard(ctx, "foo"); c != 2 {
		t.Fatal("zcard faild")
	}

	c.ZAdd(ctx, &Item{Key: "foo", ZSetValues: map[string]float64{"c": 2}})

	if c, _ := c.ZCount(ctx, "foo", "(1", "+inf"); c != 2 {
		t.Fatal("zcount faild")
	}

	c.ZIncrBy(ctx, "foo", "a", 4)
	if s, _ := c.ZScore(ctx, "foo", "a"); s != 5 {
		t.Fatal("zincrby faild")
	}

	if r, _ := c.ZRank(ctx, "foo", "a"); r != 2 {
		t.Fatal("zrank faild")
	}

	if r, err := c.ZRevRank(ctx, "foo", "a"); err != nil || r != 0 {
		t.Fatal("zrevrank faild")
	}

	c.ZRem(ctx, "foo", "a", "b", "c")
	if c, err := c.ZCard(ctx, "foo"); err != nil || c != 0 {
		t.Fatal("zrem faild")
	}

	c.ZAdd(ctx, &Item{Key: "foo", ZSetValues: map[string]float64{"a": 1, "b": 2, "c": 3, "d": 4}})
	c.ZRemRangeByRank(ctx, "foo", 2, -1)
	values, _ = c.ZRange(ctx, "foo", 0, -1)
	if values[0].Member != "a" || values[1].Member != "b" {
		t.Fatal("zremrangebyrank faild")
	}

	c.ZAdd(ctx, &Item{Key: "foo", ZSetValues: map[string]float64{"a": 1, "b": 2, "c": 3, "d": 4}})
	c.ZRemRangeByScore(ctx, "foo", "0", "2")
	values, _ = c.ZRange(ctx, "foo", 0, -1)
	if values[0].Member != "c" || values[1].Member != "d" {
		t.Fatal("zremrangebyscore faild")
	}
}

func TestOnCmd(t *testing.T) {
	c := New(Options{
		Address:  os.Getenv("REDIS_HOST"),
		PoolSize: 1,
		OnPreCmd: func(ctx context.Context, args []interface{}) context.Context {
			if len(args) == 0 {
				t.Fatal("OnPre failed")
			}

			return context.WithValue(ctx, "foo", "bar")
		},
		OnPostCmd: func(ctx context.Context, err error) {
			if ctx.Value("foo").(string) != "bar" {
				t.Fatal("OnPostCmd failed")
			}
		},
	})

	if err := c.Del(ctx, "foo"); err != nil {
		t.Fatal("start faild")
	}

	c.Set(ctx, &Item{Key: "foo", Value: []byte("123"), Flags: FlagXX})
}
