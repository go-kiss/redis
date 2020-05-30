// @Title  other.go
// @Description 其他命令诸如 flushall flushdb
// @Author  kaixinbaba
// @Update  kaixinbaba  2020/05/30
package redis

import (
	"context"
)

func (c *Client) FlushAll(ctx context.Context) error {
	args := []interface{}{"flushall"}

	return c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		_, err := conn.r.ReadStatusReply()
		return err
	})
}
