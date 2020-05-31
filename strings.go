// @Title  strings.go
// @Description 包含Keys相关的redis命令
// @Author  kaixinbaba
// @Update  kaixinbaba  2020/05/29
package redis

import (
	"context"
	"github.com/bilibili/redis/util"
)

func (c *Client) Append(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) BitCount(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) BitField(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) BitOP(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) BitOS(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

// @title    	Decr
// @description	对value值减1，原子操作
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	newValue = oldValue - 1    	int64           "返回减1后的value值"
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.DecrBy(ctx, key, 1)
}

// @title    	DecrBy
// @description	对value值减去by的值，原子操作
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	newValue = oldValue - by   	int64           "返回减by后的value值"
func (c *Client) DecrBy(ctx context.Context, key string, by int64) (int64, error) {
	return c.IncrBy(ctx, key, -by)
}

func stringConv(b []byte) (interface{}, error) {
	return util.BytesToString(b), nil
}

// @title    	Get
// @description 获取目标key的string value值
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	item 					   	*Item           "返回key对应的Item指针, value默认为string"
func (c *Client) Get(ctx context.Context, key string) (item *Item, err error) {
	return c.get(ctx, key, stringConv)
}

func intConv(b []byte) (interface{}, error) {
	i, err := util.Atoi(b)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// @title    	GetInt
// @description 获取目标key的int value值
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	item 					   	*Item           "返回key对应的Item指针, value默认为int"
func (c *Client) GetInt(ctx context.Context, key string) (item *Item, err error) {
	return c.get(ctx, key, intConv)
}

func (c *Client) get(ctx context.Context, key string, typeConv func([]byte) (interface{}, error)) (item *Item, err error) {
	args := []interface{}{"get", key}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var b []byte
		if b, err = conn.r.ReadBytesReply(); err != nil {
			item = nil
			return err
		}
		value, err := typeConv(b)
		if err != nil {
			item = nil
			return err
		}
		item = &Item{Value: value}

		return nil
	})
	return
}

func (c *Client) GetBit(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) GetRange(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) GetSet(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

// @title    	Incr
// @description	对value值加上1的值，原子操作
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	newValue = oldValue + 1   	int64           "返回加1后的value值"
func (c *Client) Incr(ctx context.Context, key string) (i int64, err error) {
	return c.IncrBy(ctx, key, 1)
}

// @title    	IncrBy
// @description	对value值加上by的值，原子操作
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	newValue = oldValue + by   	int64           "返回加by后的value值"
func (c *Client) IncrBy(ctx context.Context, key string, by int64) (i int64, err error) {
	args := []interface{}{"incrby", key, by}

	err = c.do(ctx, args, func(conn *redisConn) error {

		if err = conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		i, err = conn.r.ReadIntReply()

		return err
	})

	return
}

// @title    	IncrByFloat
// @description	对value值加上by的值，原子操作
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	newValue = oldValue + by   	float64         "返回加by后的value值"
func (c *Client) IncrByFloat(ctx context.Context, key string, by float64) (i float64, err error) {
	args := []interface{}{"incrbyfloat", key, by}

	err = c.do(ctx, args, func(conn *redisConn) error {

		if err = conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		i, err = conn.r.ReadFloat()

		return err
	})

	return
}

func (c *Client) MGet(ctx context.Context, keys []string) (items map[string]*Item, err error) {
	args := make([]interface{}, 0, len(keys)+1)

	args = append(args, "mget")
	for _, key := range keys {
		args = append(args, key)
	}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		l, err := conn.r.ReadArrayLenReply()
		if err != nil {
			return err
		}

		items = make(map[string]*Item, l)

		for i := 0; i < l; i++ {
			b, err := conn.r.ReadBytesReply()
			if err == Nil {
				continue
			}
			if err != nil {
				return err
			}

			key := keys[i]

			items[key] = &Item{Value: string(b)}
		}

		return nil
	})
	return
}

func (c *Client) MSet(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) MSetNX(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) PSetEX(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

// @title    	Set
// @description	设置K-V数据到redis中
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	item   						*Item			"需要set的Item对象的指针"
//					Key						string
//					Value					anything
//					TTL						>0
//					Flags					if 1 setnx
// @return 									error
func (c *Client) Set(ctx context.Context, item *Item) error {
	args := make([]interface{}, 0, 6)
	args = append(args, "set", item.Key, item.Value)

	if item.TTL > 0 {
		args = append(args, "EX", item.TTL)
	}

	if item.Flags&FlagNX > 0 {
		args = append(args, "NX")
	} else if item.Flags&FlagXX > 0 {
		args = append(args, "XX")
	}
	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadStatusReply()
		if err != nil {
			return err
		}

		return nil
	})
}

func (c *Client) SetBit(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) SetRange(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

// @title    	StrLen
// @description	返回value值的长度
// @auth      	kaixinbaba      时间（2020/05/29）
// @param     	key				string         "需要查询的key"
// @return    	len(value)    	int64          "返回key对应value的长度"
func (c *Client) StrLen(ctx context.Context, key string) (valueLen int64, err error) {
	args := []interface{}{"strlen", key}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var i int64
		if i, err = conn.r.ReadIntReply(); err != nil {
			return err
		}
		valueLen = i
		return nil
	})
	return
}
