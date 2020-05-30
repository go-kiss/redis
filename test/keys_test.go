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
		t.Fatalf("string Del error %s", err)
	}
	isExists, _ := client.Exists(ctx, "kaixinbaba")
	if isExists {
		t.Fatalf("string Del error the key still exists")
	}
}

func TestDump(t *testing.T) {
	client.Dump(ctx, "kaixinbaba")
}

func TestExists(t *testing.T) {
	canNotExistsKey := "NotExistsKey"
	isExists, err := client.Exists(ctx, canNotExistsKey)
	if err != nil {
		t.Fatalf("string NotExists error %s", err)
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
		t.Fatalf("string Expire error %s", err)
	}
	time.Sleep(time.Duration(expireTime) * time.Second)
	isExists, _ := client.Exists(ctx, "kaixinbaba")
	if isExists {
		t.Fatalf("string Expire error the key still exists")
	}
}
