package test

import (
	"github.com/bilibili/redis"
	"github.com/bilibili/redis/util"
	"testing"
	"time"
)

func TestDel(t *testing.T) {
	err := client.Set(ctx, &redis.Item{
		Key:   TestKey,
		Value: util.StringToBytes(TestValue),
	})
	err = client.Del(ctx, TestKey)
	if err != nil {
		t.Fatalf("string Del error %s", err)
	}
	isExists, _ := client.Exists(ctx, TestKey)
	if isExists {
		t.Fatalf("string Del error the key still exists")
	}
}

func TestDump(t *testing.T) {
	client.Dump(ctx, TestKey)
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
		Key:   TestKey,
		Value: util.StringToBytes(TestValue),
	})
	isExists, err = client.Exists(ctx, TestKey)
	if err != nil || !isExists {
		t.Fatalf("string Exists error %s", err)
	}
}


func TestExpire(t *testing.T) {
	err := client.Set(ctx, &redis.Item{
		Key:   TestKey,
		Value: util.StringToBytes(TestValue),
	})
	var expireTime int32 = 2
	err = client.Expire(ctx, TestKey, expireTime)
	if err != nil {
		t.Fatalf("string Expire error %s", err)
	}
	time.Sleep(time.Duration(expireTime) * time.Second)
	isExists, _ := client.Exists(ctx, TestKey)
	if isExists {
		t.Fatalf("string Expire error the key still exists")
	}
}
