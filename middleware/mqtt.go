package middleware

import (
	"github.com/wlgd/xutils/mq"
)

type Nats struct {
	nclient *mq.NatsClient
}

func NewNatsClient() *Nats {
	return new(Nats)
}

func (n *Nats) Notify(topic string, data interface{}) {
	if n.nclient == nil {
		return
	}
	n.nclient.Publish(topic, data)
}

func (n *Nats) Subscribe(topic string, handle func([]byte)) {
	if n.nclient == nil {
		return
	}
	n.nclient.Subscribe(topic, handle)
}

func (n *Nats) Start() error {
	client, err := mq.NewNatsClient(mq.DefaultURL, false)
	if err != nil {
		return err
	}
	n.nclient = client
	return nil
}

func (n *Nats) Connect(addr string) error {
	client, err := mq.NewNatsClient(addr, false)
	if err != nil {
		return err
	}
	n.nclient = client
	return nil
}

func (n *Nats) Stop() {
	if n.nclient == nil {
		return
	}
	n.nclient.Release()
}
