package test

import (
	"context"
	"fmt"
	"github.com/bilibili/redis"
	"testing"
)

var client = redis.New(redis.Options{
	Address:  "localhost:6379",
	PoolSize: 1,
})

var ctx = context.TODO()

func TestDecr(t *testing.T) {
	client.Set(ctx, &redis.Item{
		Key:   "kaixinbaba",
		Value: "bilibili",
	})
	_, err := client.Decr(ctx, "kaixinbaba")
	if err != nil && err.Error() != "ERR value is not an integer or out of range" {
		t.Fatalf("string value can't use decr")
	}
	client.Set(ctx, &redis.Item{
		Key:   "xjj",
		Value: "2017",
	})
	afterDecr, err := client.Decr(ctx, "xjj")
	if err != nil {
		t.Fatalf("strings Decr error %s", err)
	}
	if afterDecr != 2016 {
		t.Fatalf("afterDecr must be 2016")
	}
}

func TestIncrBy(t *testing.T) {
	//client.Get(ctx, "incrbytest")
	afterIncr, err := client.IncrBy(ctx, "incrbytest", 7)
	if err != nil {
		t.Fatalf("strings IncrBy error %s", err)
	}
	fmt.Println(afterIncr)
}

func TestSet(t *testing.T) {
	//intValue, _ := util.IntToBytes(7, 4)
	client.Set(ctx, &redis.Item{
		Key:        "xjj",
		Value:      "21",
		ZSetValues: nil,
		Flags:      0,
		TTL:        0,
	})
	item, _ := client.Get(ctx, "xjj")
	fmt.Println(item.Value)
	//if item.Value == 21 {
	//	fmt.Println(item.Value)
	//}

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

