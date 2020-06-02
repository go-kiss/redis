// @Title  strings.go
// @Description 包含Keys相关的redis命令
// @Author  kaixinbaba
// @Update  kaixinbaba  2020/05/29
package redis

import (
	"context"
	"errors"
	"github.com/bilibili/redis/util"
)

// @title    	Append
// @description	将value值拼接在目标key原先对应的value后, 如果key不存在就会创建
// @auth      	kaixinbaba      时间（2020/06/01）
// @param     	key							string			"需要操作的key"
// @param     	value						interface{}		"需要拼接的value，都会转换为string拼接在原来的value之后"
// @return    	len(newValue)    			int64           "拼接完成后的新的value的长度"
func (c *Client) Append(ctx context.Context, key string, value interface{}) (strLen int64, err error) {
	args := []interface{}{"append", key, value.(string)}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		if strLen, err = conn.r.ReadIntReply(); err != nil {
			return err
		}
		return nil
	})
	return
}

// @title    	BitCount
// @description	对给定的整个字符串进行字节长度计数，通过指定额外的 start 或 end 参数，可以让计数只在特定的位上进行。
// @auth      	kaixinbaba      时间（2020/06/01）
// @param     	key							string			"需要统计的key"
// @param     	index						int32...		"可以不传或者传两个数字， 不传代表统计整个字符串，传两个对应start和end"
// @return    	len(byte)	 	   			int64           "统计出的字节长度"
func (c *Client) BitCount(ctx context.Context, key string, index ...int32) (byteLen int64, err error) {
	// 只取前两个索引
	args := []interface{}{"bitcount", key}
	if index != nil && len(index) != 2 {
		return -1, errors.New("bitcount both start and end are required or just not pass")
	}
	for _, i := range index {
		args = append(args, i)
	}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		if byteLen, err = conn.r.ReadIntReply(); err != nil {
			return err
		}
		return nil
	})
	return
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
	if b == nil {
		return nil, nil
	}
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

// @title    	GetRange
// @description 获取对应key的value的子字符串, 开始和结束的索引都是闭区间
// @auth      	kaixinbaba      时间（2020/06/02）
// @param     	key							string			"需要操作的key"
// @param     	start						int32			"子字符串的开始索引(包含), 并且可以为负数"
// @param     	end							int32			"子字符串的结束索引(包含), 并且可以为负数"
// @return    	subValue				   	string     		"截取后的子字符串"
func (c *Client) GetRange(ctx context.Context, key string, start int32, end int32) (subValue interface{}, err error) {
	args := []interface{}{"getrange", key, start, end}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var b []byte
		if b, err = conn.r.ReadBytesReply(); err != nil {
			subValue = nil
			return err
		}
		subValue, err = stringConv(b)
		if err != nil {
			subValue = nil
			return err
		}
		return nil
	})
	return
}

// @title    	GetSet
// @description 设置对应key的value值，并且返回旧的value值，如果key原先不存在，则返回nil
// @auth      	kaixinbaba      时间（2020/06/02）
// @param     	key							string			"需要操作的key"
// @param     	newValue					string			"需要设置的value"
// @return    	oldValue				   	interface{}     "对应key原先的值，如果key不存在则返回nil"
func (c *Client) GetSet(ctx context.Context, key string, newValue interface{}) (oldValue interface{}, err error) {
	args := []interface{}{"getset", key, newValue}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var b []byte
		if b, err = conn.r.ReadBytesReply(); err != nil && err != Nil {
			oldValue = nil
			return err
		}
		oldValue, err = stringConv(b)
		if err != nil {
			oldValue = nil
			return err
		}
		return nil
	})
	return
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

// @title    	MGet
// @description	批量Get命令，不会返回不存在的key
// @auth      	kaixinbaba      时间（2020/06/01）
// @param     	keys						[]string			"需要操作的key"
// @return    	map 			   			map[string]string   "返回存在的key和value的map"
func (c *Client) MGet(ctx context.Context, keys []string) (items map[string]string, err error) {
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

		items = make(map[string]string, l)

		for i := 0; i < l; i++ {
			b, err := conn.r.ReadBytesReply()
			if err == Nil {
				continue
			}
			if err != nil {
				return err
			}

			key := keys[i]

			items[key] = string(b)
		}

		return nil
	})
	return
}

// @title    	MSet
// @description	批量Set命令
// @auth      	kaixinbaba      时间（2020/06/01）
// @param     	keys						[]items				"需要设置的item对象"
func (c *Client) MSet(ctx context.Context, items []*Item) error {
	return c.mset(ctx, items, 0)
}

// @title    	MSetNX
// @description	批量SetNX命令
// @auth      	kaixinbaba      时间（2020/06/01）
// @param     	keys						[]items				"需要设置的item对象"
func (c *Client) MSetNX(ctx context.Context, items []*Item) error {
	return c.mset(ctx, items, FlagNX)
}

func (c *Client) mset(ctx context.Context, items []*Item, flag int32) error {
	args := make([]interface{}, 0, len(items)*2+1)

	if flag&FlagNX > 0 {
		args = append(args, "msetnx")
	} else {
		args = append(args, "mset")
	}
	for _, item := range items {
		args = append(args, item.Key, item.Value)
	}
	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}
		var err error
		if flag&FlagNX > 0 {
			_, err = conn.r.ReadIntReply()
		} else {
			_, err = conn.r.ReadStatusReply()
		}
		if err != nil {
			return err
		}

		return nil
	})
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

// @title    	SetRange
// @description	返回设置后的value值的长度
// @auth      	kaixinbaba      时间（2020/06/02）
// @param     	key				string         	"需要查询的key"
// @param     	offset			int         	"偏移量"
// @param     	value			interface{} 	"从offset开始替换原value至value"
// @return    	len(value)    	int64          	"返回设置后的value值的长度"
func (c *Client) SetRange(ctx context.Context, key string, offset int, value interface{}) (strLen int, err error) {
	args := []interface{}{"setrange", key, offset, value}
	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var i int64
		if i, err = conn.r.ReadIntReply(); err != nil {
			strLen = int(i)
			return err
		}
		return nil
	})
	return
}

// @title    	StrLen
// @description	返回value值的长度
// @auth      	kaixinbaba      时间（2020/05/29）
// @param     	key				string         "需要查询的key"
// @return    	len(value)    	int64          "返回key对应value的长度"
func (c *Client) StrLen(ctx context.Context, key string) (valueLen int, err error) {
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
			valueLen = int(i)
			return err
		}
		return nil
	})
	return
}
