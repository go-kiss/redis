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

	// set with expire
	c.Set(ctx, &Item{Key: "set_with_expire", Value: []byte("test"), TTL: 86400})
	item, _ = c.Get(ctx, "set_with_expire")
	if !bytes.Equal(item.Value, []byte("test")) {
		t.Fatal("get foo failed")
	}
	// 获取 TTL
	if ttl, err := c.TTL(ctx, "set_with_expire"); err != nil {
		t.Fatal("ttl failed")
	} else if ttl <= 0 {
		t.Fatal("set with expire failed")
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

	values, _ = c.ZRangeByScore(ctx, "foo", 0, 1, 0, 0)
	if values[0].Member != "a" ||
		values[0].Score != 1 {

		t.Fatal("zrangebyscore faild")
	}

	values, _ = c.ZRevRangeByScore(ctx, "foo", 2, 1, 0, 0)
	if values[1].Member != "a" ||
		values[1].Score != 1 ||
		values[0].Member != "b" ||
		values[0].Score != 2 {

		t.Fatal("zrevrangebyscore faild")
	}

	values, _ = c.ZRevRangeByScore(ctx, "foo", 2, 1, 1, 1)
	if values[0].Member != "a" ||
		values[0].Score != 1 {

		t.Fatal("zrevrangebyscore faild")
	}

	if c, _ := c.ZCard(ctx, "foo"); c != 2 {
		t.Fatal("zcard faild")
	}

	c.ZAdd(ctx, &Item{Key: "foo", ZSetValues: map[string]float64{"c": 2}})

	if c, _ := c.ZCount(ctx, "foo", "(1", "+inf"); c != 2 {
		t.Fatal("zcount faild")
	}

	c.ZIncrBy(ctx, "foo", "a", 4.05)
	if s, _ := c.ZScore(ctx, "foo", "a"); s-5.05 >= 0.000001 {
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

func TestMget(t *testing.T) {
	c := New(Options{
		Address:  os.Getenv("REDIS_HOST"),
		PoolSize: 1,
	})

	err := c.Set(ctx, &Item{Key: "key_m_1", Value: []byte("value_m_1")})
	if err != nil {
		t.Fatal("Set Failed")
	}
	err = c.Set(ctx, &Item{Key: "key_m_2", Value: []byte("value_m_2")})
	if err != nil {
		t.Fatal("Set Failed")
	}
	// 删除数据
	defer c.Del(ctx, "key_m_1", "key_m_2")

	// 批量获取
	items, err := c.MGet(ctx, []string{"key_m_1", "key_m_2", "key_m_3"})
	if err != nil {
		t.Fatal("MGet Failed")
	}

	// 校验获取的值与插入的一致
	if string(items["key_m_1"].Value) != "value_m_1" {
		t.Fatal("MGet Failed")
	}

	if string(items["key_m_2"].Value) != "value_m_2" {
		t.Fatal("MGet Failed")
	}

	if _, ok := items["key_m_3"]; ok {
		t.Fatal("MGet Failed")
	}
}

func TestEval(t *testing.T) {
	c := New(Options{
		Address:  os.Getenv("REDIS_HOST"),
		PoolSize: 1,
	})

	if err := c.Del(ctx, "foo"); err != nil {
		t.Fatal("start faild")
	}

	// array
	val, _ := c.Eval(ctx, "return {\"abcc\",1, {\"b\"} }", []string{})
	if v, err := val.Array(); v[0].(string) != "abcc" ||
		v[1].(int64) != 1 ||
		v[2].([]interface{})[0].(string) != "b" {
		t.Fatal("eval faild", err)
	}

	// int64
	val, _ = c.Eval(ctx, "return 64", []string{})
	if v, err := val.Int64(); v != 64 {
		t.Fatal("eval faild", err)
	}
	// string
	val, _ = c.Eval(ctx, "return ARGV[1]", []string{}, "a\nb\nc")
	if v, err := val.String(); v != "a\nb\nc" {
		t.Fatal("eval faild", err)
	}

	// status
	val, _ = c.Eval(ctx, "return redis.call('set',KEYS[1],ARGV[1])", []string{"foo"}, "hello")
	c.Del(ctx, "foo")
	if v, err := val.String(); v != "OK" {
		t.Fatal("eval faild", err)
	}

	// with no params
	c.Eval(ctx, "return redis.call('set',KEYS[1],ARGV[1])", []string{"foo"}, "hello")
	defer c.Del(ctx, "foo")

	item, _ := c.Get(ctx, "foo")
	if !bytes.Equal(item.Value, []byte("hello")) {
		t.Fatal("eval failed")
	}
}

func TestSet(t *testing.T) {
	c := New(Options{
		Address:  os.Getenv("REDIS_HOST"),
		PoolSize: 1,
	})
	var err error

	if err := c.Del(ctx, "foo"); err != nil {
		t.Fatal("start faild")
	}

	// sadd
	if err := c.SAdd(ctx, "foo", []byte("aaa")); err != nil {
		t.Fatalf("add foo aaa failed, err: %v", err)
	}
	if err := c.SAdd(ctx, "foo", []byte("bbb"), []byte("ccc")); err != nil {
		t.Fatalf("add foo ddd,eee,fff failed, err: %v", err)
	}

	// scard
	if card, err := c.SCard(ctx, "foo"); err != nil || card != 3 {
		t.Fatalf("Key: foo, Scard: %d; Failed, err: %v", card, err)
	}

	// smembers
	items, err := c.SMembers(ctx, "foo")
	if err != nil {
		t.Fatalf("get foo's members failed, err: %v", err)
	}
	t.Logf("items: %#v", items)

	// sismember
	if result, err := c.SIsMember(ctx, "foo", []byte("bbb")); err != nil || result != true {
		t.Fatalf("Key: foo, SIsMember: %t; Failed, err: %v", result, err)
	}

	// spop
	item, _ := c.SPop(ctx, "foo", 1)
	if len(item) < 1 {
		t.Fatal("spop foo empty")
	}
	set := map[string]bool{"aaa": true, "bbb": true, "ccc": true}
	exists := set[string(item[0])]
	if !exists {
		t.Fatal("spop foo failed")
	}

	// srem
	if err := c.SAdd(ctx, "foo", []byte("ddd"), []byte("eee"), []byte("fff")); err != nil {
		t.Fatalf("add foo ddd failed, err: %v", err)
	}
	if result, err := c.SRem(ctx, "foo", []byte("ddd"), []byte("fff")); err != nil || result == 0 {
		t.Fatalf("Key: foo, SRem: %d; Failed, err: %v", result, err)
	}

	// srem 删除不存在的值
	if result, err := c.SRem(ctx, "foo", []byte("zzz"), []byte("xxx")); err != nil || result == 0 {
		t.Logf("Key: foo, SRem: %d; Failed, err: %v", result, err)
	}
}

func TestHash(t *testing.T) {
	c := New(Options{
		Address:  os.Getenv("REDIS_HOST"),
		PoolSize: 1,
	})

	testKey := "hashtest"

	if err := c.Del(ctx, testKey); err != nil {
		t.Fatal("start faild")
	}

	// hset
	added, _ := c.HSet(ctx, testKey, map[string]string{"name": "bilibili", "age": "20"})
	if added != 2 {
		t.Fatalf("hset %s failed", testKey)
	}
	// hget
	item, _ := c.HGet(ctx, testKey, "name")
	if !bytes.Equal(item.Value, []byte("bilibili")) {
		t.Fatalf("hget %s failed", testKey)
	}

	// hgetall
	hgetallItem, _ := c.HGetAll(ctx, testKey)
	for f, v := range hgetallItem.HashValues {
		t.Logf("Key: %s, field: %s, value: %s", testKey, f, v)
	}

}
