package rpc

import (
	"net"
	"time"

	"github.com/lesismal/arpc/extension/pubsub"
)

type arpclient struct {
	Conn *pubsub.Client
}

func NewArpc(addr string) IRpc {
	c := &arpclient{}
	c.Conn, _ = pubsub.NewClient(func() (net.Conn, error) {
		return net.DialTimeout("tcp", addr, time.Second*3)
	})
	return c
}

func (n *arpclient) Subscribe(topicstr string, handler func([]byte)) error {
	if n.Conn == nil {
		return nil
	}
	return n.Conn.Subscribe(topicstr, func(topic *pubsub.Topic) {
		handler(topic.Data)
	}, time.Second)
}

func (n *arpclient) Publish(topicstr string, v interface{}) error {
	if n.Conn == nil {
		return nil
	}
	return n.Conn.Publish(topicstr, v, 2*time.Second)
}

func (n *arpclient) Release() {
	if n.Conn == nil {
		return
	}
	n.Conn.Stop()
}
