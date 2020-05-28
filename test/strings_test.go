package test

import (
	"context"
	"github.com/bilibili/redis"
	"github.com/bilibili/redis/util"
	"testing"
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
