package redis

import "context"

func (c *Client) ZAdd(ctx context.Context, item *Item) (added int64, err error) {
	args := make([]interface{}, 0, 4+len(item.ZSetValues))
	args = append(args, "zadd", item.Key)

	if item.Flags&FlagNX > 0 {
		args = append(args, "NX")
	} else if item.Flags&FlagXX > 0 {
		args = append(args, "XX")
	}

	if item.Flags&FlagCH > 0 {
		args = append(args, "CH")
	}

	for member, score := range item.ZSetValues {
		args = append(args, score, member)
	}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		added, err = conn.r.ReadIntReply()
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func (c *Client) ZIncrBy(ctx context.Context, key, member string, by float64) error {
	args := []interface{}{"zincrby", key, by, member}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadFloat()
		return err
	})
}

func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) (values []*ZSetValue, err error) {
	return c.zrange(ctx, "zrange", key, start, stop, 0, 0)
}

func (c *Client) ZRevRange(ctx context.Context, key string, start, stop int64) (values []*ZSetValue, err error) {
	return c.zrange(ctx, "zrevrange", key, start, stop, 0, 0)
}

func (c *Client) ZRangeByScore(ctx context.Context, key string, min, max float64, offset, count int64) (values []*ZSetValue, err error) {
	return c.zrange(ctx, "zrangebyscore", key, min, max, offset, count)
}

func (c *Client) ZRevRangeByScore(ctx context.Context, key string, max, min float64, offset, count int64) (values []*ZSetValue, err error) {
	return c.zrange(ctx, "zrevrangebyscore", key, max, min, offset, count)
}

func (c *Client) zrange(ctx context.Context, cmd, key string, start, stop interface{}, offset, count int64) (values []*ZSetValue, err error) {
	args := []interface{}{cmd, key, start, stop, "WITHSCORES"}
	if count > 0 {
		args = append(args, "LIMIT", offset, count)
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

		values = make([]*ZSetValue, 0, l)
		for i := 0; i < l/2; i++ {
			b, err := conn.r.ReadBytesReply()
			if err != nil {
				return err
			}
			f, err := conn.r.ReadFloat()
			if err != nil {
				return err
			}

			values = append(values, &ZSetValue{Member: string(b), Score: f})
		}

		return nil
	})

	return
}

func (c *Client) ZRank(ctx context.Context, key, member string) (rank int64, err error) {
	return c.zrank(ctx, "zrank", key, member)
}
func (c *Client) ZRevRank(ctx context.Context, key, member string) (rank int64, err error) {
	return c.zrank(ctx, "zrevrank", key, member)
}

func (c *Client) zrank(ctx context.Context, cmd, key, member string) (rank int64, err error) {
	args := []interface{}{cmd, key, member}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		rank, err = conn.r.ReadIntReply()
		return err
	})

	return
}

func (c *Client) ZScore(ctx context.Context, key, member string) (score float64, err error) {
	args := []interface{}{"zscore", key, member}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		score, err = conn.r.ReadFloat()
		return err
	})

	return
}

func (c *Client) ZCard(ctx context.Context, key string) (card int64, err error) {
	args := []interface{}{"zcard", key}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		card, err = conn.r.ReadIntReply()

		return err
	})
	return
}

func (c *Client) ZCount(ctx context.Context, key, min, max string) (i int64, err error) {
	args := []interface{}{"zcount", key, min, max}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
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

func (c *Client) ZRem(ctx context.Context, keys ...string) error {
	args := make([]interface{}, 0, 1+len(keys))

	args = append(args, "zrem")
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

func (c *Client) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (i int64, err error) {
	args := []interface{}{"zremrangebyrank", key, start, stop}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
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

func (c *Client) ZRemRangeByScore(ctx context.Context, key, min, max string) (i int64, err error) {
	args := []interface{}{"zremrangebyscore", key, min, max}

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
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
