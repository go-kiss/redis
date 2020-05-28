package redis

import "context"

func (c *Client) Eval(ctx context.Context, script string, keys []string, argvs ...interface{}) (result *EvalReturn, err error) {
	args := make([]interface{}, 0, len(keys)+len(argvs)+2)
	args = append(args, "eval", script, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}
	args = append(args, argvs...)

	err = c.do(ctx, args, func(conn *redisConn) error {
		if err := conn.w.WriteArgs(args); err != nil {
			return err
		}

		if err := conn.w.Flush(); err != nil {
			return err
		}

		b, err := conn.r.ReadInterfaceReply()
		if err != nil {
			return err
		}

		result = &EvalReturn{
			val: b,
		}

		return nil
	})

	return
}
