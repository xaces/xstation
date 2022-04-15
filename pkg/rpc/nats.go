package rpc

import "github.com/wlgd/xutils/mq"

type natsclient struct {
	Conn *mq.NatsClient
}

func NewNats(addr string) Interface {
	c := &natsclient{}
	c.Conn, _ = mq.NewNatsClient(addr, false)
	return c
}

func (n *natsclient) Subscribe(topicstr string, handler func([]byte)) {
	if n.Conn == nil {
		return
	}
	n.Conn.Subscribe(topicstr, handler)
}

func (n *natsclient) Publish(topicstr string, v interface{}) error {
	if n.Conn == nil {
		return nil
	}
	return n.Conn.Publish(topicstr, v)
}

func (n *natsclient) Relase() {
	if n.Conn == nil {
		return
	}
	n.Conn.Release()
}
