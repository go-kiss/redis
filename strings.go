package redis

import "context"

func (c *Client) Get(ctx context.Context, key string) (item *Item, err error) {
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

		item = &Item{Value: b}

		return nil
	})
	return
}

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

func (c *Client) DecrBy(ctx context.Context, key string, by int64) (int64, error) {
	return c.IncrBy(ctx, key, -by)
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

			items[key] = &Item{Value: b}
		}

		return nil
	})
	return
}