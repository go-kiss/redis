// @Title  keys.go
// @Description 包含Keys相关的redis命令
// @Author  kaixinbaba
// @Update  kaixinbaba  2020/05/29
package redis

import (
	"context"
)

// @title    	Del
// @description	删除一个或多个key
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	keys						...string			"需要删除的key，可变长度"
// @return    						    	error
func (c *Client) Del(ctx context.Context, keys ...string) error {
	args := make([]interface{}, 0, 1+len(keys))

	args = append(args, "del")
	for _, key := range keys {
		args = append(args, key)
	}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadIntReply()

		return err
	})
}

// @title    	Dump
// @description	把指定key值dump的value序列化，转换成字节数组
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	b					    	[]byte          "value值序列化后的字节数组"
func (c *Client) Dump(ctx context.Context, key string) (b []byte, err error) {
	args := []interface{}{"dump", key}

	c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		if b, err = conn.r.ReadBytesReply(); err != nil {
			return err
		}
		return nil
	})
	return
}

// @title    	Exists
// @description	判断指定key是否存在
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	isExists			    	bool 	        "true 存在，false 不存在"
func (c *Client) Exists(ctx context.Context, key string) (isExists bool, err error) {
	args := []interface{}{"exists", key}

	c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var i int64
		if i, err = conn.r.ReadIntReply(); err != nil {
			isExists = false
			return err
		}
		isExists = i == 1
		return nil
	})
	return
}

// @title    	Expire
// @description	对指定的key设置超时时间
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @param     	expireTime					int32			"超时时间，单位秒"
// @return    								error
func (c *Client) Expire(ctx context.Context, key string, expireTime int32) error {
	args := []interface{}{"expire", key, expireTime}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadIntReply()
		return err
	})
}

func (c *Client) ExpireAt(ctx context.Context, key string, ttl int32) error {
	// TODO
	return nil
}

func (c *Client) Keys(ctx context.Context, pattern string) (keys []string, err error) {
	// TODO
	return nil, nil
}

func (c *Client) Migrate(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) Move(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) Object(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) Persist(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) PExpire(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) PExpireAt(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) PTTL(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) RandomKey(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) Rename(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) RenameNX(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

// @title    	Restore
// @description	将字节数组反序列化到指定的key中，不设置该key的过期时间
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要反序列化的key"
// @param     	content						[]byte			"序列化的字节数组"
// @return    								error
func (c *Client) Restore(ctx context.Context, key string, content []byte) (err error) {
	return c.RestoreEX(ctx, key, 0, content)
}

// @title    	Restore
// @description	将字节数组反序列化到指定的key中，可以设置该key的过期时间
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要反序列化的key"
// @param     	expireTime					int32			"key的过期时间，单位秒"
// @param     	content						[]byte			"序列化的字节数组"
// @return    								error
func (c *Client) RestoreEX(ctx context.Context, key string, expireTime int32, content []byte) (err error) {
	// restore的ttl是毫秒单位的
	args := []interface{}{"restore", key, expireTime * 1000, content}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err = conn.r.ReadStatusReply()
		return err
	})
}


func (c *Client) Sort(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

// @title    	TTL
// @description	查询指定key的超时剩余时间
// @auth      	kaixinbaba      时间（2020/05/30）
// @param     	key							string			"需要操作的key"
// @return    	expireTime			    	int32 	        "目标key剩余过期时间，如果返回-1代表不会过期"
func (c *Client) TTL(ctx context.Context, key string) (expireTime int32, err error) {
	args := []interface{}{"ttl", key}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		var i int64
		i, err = conn.r.ReadIntReply()
		if err != nil {
			return err
		}

		if i == -2 {
			err = Nil
			return err
		}

		expireTime = int32(i)

		return err
	})
	return
}

func (c *Client) Type(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) Wait(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}

func (c *Client) Scan(ctx context.Context, args ...interface{}) error {
	// TODO
	return nil
}