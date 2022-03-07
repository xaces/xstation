package middleware

import (
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/mq"
)

type NatsOption struct {
	Address     string
	TopicOnline string
	TopicStatus string
	TopicAlarm  string
	TopicEvent  string
}

type Nats struct {
	Conn *mq.NatsClient
	opt  NatsOption
}

func (n *Nats) Online(a *xproto.Access) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.opt.TopicOnline, a)
}

func (n *Nats) Status(s *xproto.Status) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.opt.TopicStatus, s)
}

func (n *Nats) Alarm(a *xproto.Alarm) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.opt.TopicAlarm, a)
}

func (n *Nats) Event(e *xproto.Event) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.opt.TopicEvent, e)
}

func (n *Nats) Start(o NatsOption) (err error) {
	n.opt = o
	if o.Address == "" {
		o.Address = mq.DefaultURL
	}
	n.Conn, err = mq.NewNatsClient(o.Address, false)
	return
}

func (n *Nats) Stop() {
	if n.Conn == nil {
		return
	}
	n.Conn.Release()
}
