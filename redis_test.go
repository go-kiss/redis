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
