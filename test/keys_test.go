package test

import (
	"github.com/bilibili/redis"
	"testing"
	"time"
)

func TestDel(t *testing.T) {
	err := client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	err = client.Del(ctx, "kaixinbaba")
	if err != nil {
		t.Fatalf("keys Del error %s", err)
	}
	isExists, _ := client.Exists(ctx, "kaixinbaba")
	if isExists {
		t.Fatalf("string Del error the key still exists")
	}
}

func TestDump(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	b, err := client.Dump(ctx, "kaixinbaba")
	if err != nil {
		t.Fatalf("keys Dump error %s", err)
	}
	client.RestoreEX(ctx, "restorexjj", 10, b)
}

func TestExists(t *testing.T) {
	canNotExistsKey := "NotExistsKey"
	isExists, err := client.Exists(ctx, canNotExistsKey)
	if err != nil {
		t.Fatalf("keys NotExists error %s", err)
	}
	if isExists {
		t.Fatalf("The key [%s] exists!", canNotExistsKey)
	}
	err = client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	isExists, err = client.Exists(ctx, "kaixinbaba")
	if err != nil || !isExists {
		t.Fatalf("string Exists error %s", err)
	}
}


func TestExpire(t *testing.T) {
	err := client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	var expireTime int32 = 2
	err = client.Expire(ctx, "kaixinbaba", expireTime)
	if err != nil {
		t.Fatalf("keys Expire error %s", err)
	}
	time.Sleep(time.Duration(expireTime) * time.Second)
	isExists, _ := client.Exists(ctx, "kaixinbaba")
	if isExists {
		t.Fatalf("string Expire error the key still exists")
	}
}

func TestTTL(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	ttl, err := client.TTL(ctx, "kaixinbaba")
	if err != nil {
		t.Fatalf("keys TTL error %s", err)
	}
	if ttl != -1 {
		t.Fatalf("expireTime should be -1, mean no expire, but got %d", ttl)
	}
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
		TTL: 3,
	})
	ttl, err = client.TTL(ctx, "kaixinbaba")
	if ttl != 3 {
		t.Fatalf("expireTime should be 3")
	}
	time.Sleep(1*time.Second)
	ttl, err = client.TTL(ctx, "kaixinbaba")
	if ttl != 2 {
		t.Fatalf("expireTime should be 2")
	}
}
