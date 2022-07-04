package rpc

import "github.com/xaces/xutils/mq"

func NewNats(addr string) IRpc {
	c, err := mq.NewNats(addr)
	if err != nil {
		return nil
	}
	c.Run()
	return c
}
