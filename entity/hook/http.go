package hook

import (
	"github.com/xaces/xproto"
	"github.com/xaces/xutils"
)

type Http struct {
	urlOnline string
	urlStatus string
	urlAlarm  string
	urlEvent  string
}

func NewHttp(o Option) *Http {
	return &Http{
		urlOnline: o.Address + "/" + o.Online,
		urlStatus: o.Address + "/" + o.Status,
		urlAlarm:  o.Address + "/" + o.Alarm,
		urlEvent:  o.Address + "/" + o.Event,
	}
}

func (h *Http) Online(deviceId uint, a *xproto.Access) {
	xutils.HttpPost(h.urlOnline, newOnline(deviceId, a), nil)
}

func (h *Http) Status(deviceId uint, s *xproto.Status) {
	xutils.HttpPost(h.urlStatus, newStatus(deviceId, s), nil)
}

func (h *Http) Alarm(deviceId uint, a *xproto.Alarm) {
	xutils.HttpPost(h.urlAlarm, newAlarm(deviceId, a), nil)
}

func (h *Http) Event(deviceId uint, e *xproto.Event) {
	xutils.HttpPost(h.urlEvent, newEvent(deviceId, e), nil)
}

func (h *Http) Stop() {
}
