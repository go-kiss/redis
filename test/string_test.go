package test

import (
	"context"
	"github.com/bilibili/redis"
	"github.com/bilibili/redis/util"
	"testing"
	"time"
)

var client = redis.New(redis.Options{
	Address:  "localhost:6379",
	PoolSize: 1,
})

var ctx = context.TODO()

var TestKey = "kaixinbaba"
var TestValue = "bilibili"

func TestSetAndGet(t *testing.T) {
	err := client.Set(ctx, &redis.Item{
		Key:   TestKey,
		Value: util.StringToBytes(TestValue),
	})
	if err != nil {
		t.Fatalf("string Set error %s", err)
	}
	item, err := client.Get(ctx, TestKey)
	if err != nil {
		t.Fatalf("string Get error %s", err)
	}
	if util.BytesToString(item.Value) != TestValue {
		t.Fatalf("string Get result not equal TestValue")
	}
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
