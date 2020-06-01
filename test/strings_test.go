package test

import (
	"context"
	"fmt"
	"github.com/bilibili/redis"
	"testing"
	"time"
)

var client = redis.New(redis.Options{
	Address:  "localhost:6379",
	PoolSize: 1,
})

var ctx = context.TODO()

func TestAppend(t *testing.T) {
	client.Del(ctx, "kaixinbaba")
	valueLen, err := client.Append(ctx, "kaixinbaba", "redis")
	if err != nil {
		t.Fatalf("strings Append error %s", err)
	}
	if valueLen != 5 {
		t.Fatalf("'redis' length must be 5")
	}
	valueLen, _ = client.Append(ctx, "kaixinbaba", "bilibili")
	if valueLen != 13 {
		t.Fatalf("'redisbilibili' length must be 13")
	}
}

func TestBitCount(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	byteLen, err := client.BitCount(ctx, "kaixinbaba")
	if err != nil {
		t.Fatalf("strings Append error %s", err)
	}
	if byteLen != 30 {
		t.Fatalf("'bilibili' byte length must be 30")
	}
	byteLen, _ = client.BitCount(ctx, "kaixinbaba", 0, 1)
	if byteLen != 7 {
		t.Fatalf("'bilibili' byte length must be 7")
	}
	_, err = client.BitCount(ctx, "kaixinbaba", 5)
	if err == nil {
		t.Fatalf("Can't just provide start or end")
	}
	_, err = client.BitCount(ctx, "kaixinbaba", 0, 1, 2, 3)
	if err == nil {
		t.Fatalf("bitcount can only provide 2 arguments")
	}
}

func TestDecr(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	_, err := client.Decr(ctx, "kaixinbaba")
	// value值不是数字类型的就应该抛异常
	if err != nil && err.Error() != "ERR value is not an integer or out of range" {
		t.Fatalf("string value can't use Decr")
	}
	client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: 2017,
	})
	newValue, err := client.Decr(ctx, "xjj")
	if err != nil {
		t.Fatalf("strings Decr error %s", err)
	}
	// 2017 - 1 = 2016
	if newValue != 2016 {
		t.Fatalf("newValue must be 2016")
	}
}

func TestDecrBy(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	_, err := client.DecrBy(ctx, "kaixinbaba", 7)
	// value值不是数字类型的就应该抛异常
	if err != nil && err.Error() != "ERR value is not an integer or out of range" {
		t.Fatalf("string value can't use DecrBy")
	}
	client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: 2017,
	})
	newValue, err := client.DecrBy(ctx, "xjj", 1002)
	if err != nil {
		t.Fatalf("strings DecrBy error %s", err)
	}
	// 2017 - 1002 = 1015
	if newValue != 1015 {
		t.Fatalf("newValue must be 1015")
	}
}

func TestGet(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	item, err := client.GetInt(ctx, "kaixinbaba")
	if err == nil {
		t.Fatalf("Can't GetInt with string value")
	}
	item, err = client.Get(ctx, "kaixinbaba")
	if item.Value != "bilibili" {
		t.Fatalf("value must be bilibili")
	}
}

func TestGetInt(t *testing.T) {

}

func TestIncr(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	_, err := client.Incr(ctx, "kaixinbaba")
	// value值不是数字类型的就应该抛异常
	if err != nil && err.Error() != "ERR value is not an integer or out of range" {
		t.Fatalf("string value can't use Incr")
	}
	client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: 2017,
	})
	newValue, err := client.Incr(ctx, "xjj")
	if err != nil {
		t.Fatalf("strings Incr error %s", err)
	}
	// 2017 + 1 = 2018
	if newValue != 2018 {
		t.Fatalf("newValue must be 2018")
	}
}

func TestIncrBy(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	_, err := client.IncrBy(ctx, "kaixinbaba", 7)
	// value值不是数字类型的就应该抛异常
	if err != nil && err.Error() != "ERR value is not an integer or out of range" {
		t.Fatalf("string value can't use IncrBy")
	}
	client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: 2017,
	})
	newValue, err := client.IncrBy(ctx, "xjj", 1002)
	if err != nil {
		t.Fatalf("strings IncrBy error %s", err)
	}
	// 2017 + 1002 = 3019
	if newValue != 3019 {
		t.Fatalf("newValue must be 3019")
	}
}

func TestIncrByFloat(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	_, err := client.IncrByFloat(ctx, "kaixinbaba", 7)
	// value值不是数字类型的就应该抛异常
	if err != nil && err.Error() != "ERR value is not a valid float" {
		t.Fatalf("string value can't use IncrByFloat")
	}
	client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: 2017,
	})
	newValue, err := client.IncrByFloat(ctx, "xjj", 10.02)
	if err != nil {
		t.Fatalf("strings IncrByFloat error %s", err)
	}
	// 2017 + 10.02 = 2027.02
	if newValue != 2027.02 {
		t.Fatalf("newValue must be 2027.02")
	}
}

func TestMGet(t *testing.T) {
	items, err := client.MGet(ctx, []string{"mset1", "mset2"})
	if err != nil {
		t.Fatalf("strings MGet error %s", err)
	}
	fmt.Println(items)
}

func TestMSet(t *testing.T) {
	err := client.MSet(ctx, []*redis.Item{
		{
			Key:   "mset1",
			Value: "123",
		},
		{
			Key:   "mset2",
			Value: "355",
		},
	})
	if err != nil {
		t.Fatalf("strings MSet error %s", err)
	}
}

func TestMSetNX(t *testing.T) {
	err := client.MSetNX(ctx, []*redis.Item{
		{
			Key:   "mset1",
			Value: "bilibili",
		},
		{
			Key:   "mset2",
			Value: "redis",
		},
	})
	if err != nil {
		t.Fatalf("strings MSetNX error %s", err)
	}
	client.MSetNX(ctx, []*redis.Item{
		{
			Key:   "mset1",
			Value: "123",
		},
		{
			Key:   "mset2",
			Value: "345",
		},
	})
	if item, _ := client.Get(ctx, "mset1"); item.Value != "bilibili" {
		t.Fatalf("strings MSetNX can't change exists key")
	}
	if item, _ := client.Get(ctx, "mset2"); item.Value != "redis" {
		t.Fatalf("strings MSetNX can't change exists key")
	}
}

func TestSet(t *testing.T) {
	// set int
	err := client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: 21,
	})
	if err != nil {
		t.Fatalf("strings Set int error %s", err)
	}
	item, _ := client.GetInt(ctx, "xjj")
	if item.Value != 21 {
		t.Fatalf("set int value != 21")
	}
	// set string
	err = client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: "bilibili",
	})
	if err != nil {
		t.Fatalf("strings Set string error %s", err)
	}
	item, _ = client.Get(ctx, "xjj")
	if item.Value != "bilibili" {
		t.Fatalf("set string value != bilibili")
	}
}

func TestSetEX(t *testing.T) {
	// expire 3 second
	err := client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: "bilibili",
		TTL:   3,
	})
	if err != nil {
		t.Fatalf("strings SetEX  error %s", err)
	}
	item, _ := client.Get(ctx, "xjj")
	if item == nil {
		t.Fatalf("key must exists")
	}
	time.Sleep(3 * time.Second)
	item, _ = client.Get(ctx, "xjj")
	if item != nil {
		t.Fatalf("key already expire, item must be nil")
	}
}

func TestSetNX(t *testing.T) {
	// remove key first
	client.Del(ctx, "xjj")
	// this set will be successful
	client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: "bilibili",
		Flags: 1,
	})
	// this set will be failure
	client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: "acfun",
		Flags: 1,
	})
	item, _ := client.Get(ctx, "xjj")
	if item.Value != "bilibili" {
		t.Fatalf("value is not correct, expect 'bilibili', but got '%v'", item.Value)
	}
}

func TestStrLen(t *testing.T) {
	err := client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	valueLen, err := client.StrLen(ctx, "kaixinbaba")
	if err != nil {
		t.Fatalf("strings StrLen error %s", err)
	}
	if valueLen != int64(len("bilibili")) {
		t.Fatalf("strings StrLen result [%d] not equal to len(TestValue) [%d]", valueLen, len("bilibili"))
	}
}
