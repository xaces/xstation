package hook

import (
	"github.com/xaces/xproto"
	"github.com/xaces/xutils/mq"
)

type Nats struct {
	Cli *mq.Client
	Option
}

func NewNats(o Option) *Nats {
	c := &Nats{Option: o}
	if conn, err := mq.NewPublish(&mq.Options{Address: o.Address, Goc: 4}, mq.NewNats); err == nil {
		c.Cli = conn
	}
	return c
}

func (o *Nats) PublishOnline(v *xproto.Access) {
	if o.Cli == nil {
		return
	}
	o.Cli.Publish(o.Online, &online{Code: 50001, Access: v})
}

func (o *Nats) PublishStatus(v *xproto.Status) {
	if o.Cli == nil {
		return
	}
	o.Cli.Publish(o.Status, &status{Code: 50002, Status: v})
}

func (o *Nats) PublishAlarm(v *xproto.Alarm) {
	if o.Cli == nil {
		return
	}
	o.Cli.Publish(o.Alarm, &alarm{Code: 50003, Alarm: v})
}

func (o *Nats) PublishEvent(v *xproto.Event) {
	if o.Cli == nil {
		return
	}
	o.Cli.Publish(o.Event, &event{Code: 50004, Event: v})
}

func (o *Nats) Stop() {
	if o.Cli == nil {
		return
	}
	o.Cli.Shutdown()
}
