package hook

import (
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/mq"
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
	c.Conn, _ = mq.NewNatsClient(o.Address, false)
	return c
}

func (n *Nats) Online(deviceId uint, a *xproto.Access) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.Opt.Online, &Online{DeviceId: deviceId, Access: a})
}

func (n *Nats) Status(deviceId uint, s *xproto.Status) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.Opt.Status, &Status{DeviceId: deviceId, Status: s})
}

func (n *Nats) Alarm(deviceId uint, a *xproto.Alarm) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.Opt.Alarm, &Alarm{DeviceId: deviceId, Alarm: a})
}

func (n *Nats) Event(deviceId uint, e *xproto.Event) {
	if n.Conn == nil {
		return
	}
	n.Conn.Publish(n.Opt.Event, &Event{DeviceId: deviceId, Event: e})
}

func (n *Nats) Stop() {
	if n.Conn == nil {
		return
	}
	n.Conn.Release()
}
