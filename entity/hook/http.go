package hook

import (
	"github.com/xaces/xproto"
	"github.com/xaces/xutils"
)

type Http struct {
	Option
}

func NewHttp(o Option) *Http {
	c := &Http{}
	c.Online = o.Address + "/" + o.Online
	c.Status = o.Address + "/" + o.Status
	c.Alarm = o.Address + "/" + o.Alarm
	c.Event = o.Address + "/" + o.Event
	return c
}

func (o *Http) PublishOnline(v *xproto.Access) {
	xutils.HttpPost(o.Address+"/"+o.Online, &online{Code: 50001, Access: v}, nil)
}

func (o *Http) PublishStatus(v *xproto.Status) {
	xutils.HttpPost(o.Status, &status{Code: 50002, Status: v}, nil)
}
func (o *Http) PublishAlarm(v *xproto.Alarm) {
	xutils.HttpPost(o.Alarm, &alarm{Code: 50003, Alarm: v}, nil)
}
func (o *Http) PublishEvent(v *xproto.Event) {
	xutils.HttpPost(o.Event, &event{Code: 50004, Event: v}, nil)
}

func (o *Http) Stop() {
}
