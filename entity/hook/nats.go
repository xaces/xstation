package hook

import (
	"github.com/xaces/xproto"
	"github.com/xaces/xutils/mq"
)

type Nats struct {
	Conn *mq.NatsClient
	Opt  Option
}

func NewNats(o Option) *Nats {
	if o.Address == "" {
		o.Address = mq.DefaultURL
	}
	c := &Nats{Opt: o}
	if conn, err := mq.NewNatsClient(o.Address, false); err == nil {
		c.Conn = conn
	}
	return c
}

func (n *Nats) Online(deviceId uint, a *xproto.Access) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.Opt.Online, newOnline(deviceId, a))
}

func (n *Nats) Status(deviceId uint, s *xproto.Status) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.Opt.Status, newStatus(deviceId, s))
}

func (n *Nats) Alarm(deviceId uint, a *xproto.Alarm) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.Opt.Alarm, newAlarm(deviceId, a))
}

func (n *Nats) Event(deviceId uint, e *xproto.Event) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.Opt.Event, newEvent(deviceId, e))
}

func (n *Nats) Stop() {
	if n.Conn == nil {
		return
	}
	n.Conn.Release()
}
