package redis

import "context"

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

func (c *Client) Expire(ctx context.Context, key string, ttl int32) error {
	args := []interface{}{"expire", key, ttl}

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

func (c *Client) TTL(ctx context.Context, key string) (ttl int32, err error) {
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

		ttl = int32(i)

		return err
	})
	return
}