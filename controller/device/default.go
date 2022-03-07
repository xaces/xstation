package device

import (
	"github.com/wlgd/xproto"
)

type defaultBroker struct {
}

func (d *defaultBroker) Online(a *xproto.Access) {
	updateDevOnline(a)
}

func (d *defaultBroker) Status(s *xproto.Status) {
	Handler.addStatus(s)
}

func (d *defaultBroker) Alarm(a *xproto.Alarm) {
	Handler.addAlarm(a)
}

func (d *defaultBroker) Event(e *xproto.Event) {
	devEventHandler(e)
}

func (d *defaultBroker) Stop() {
	// 处理缓存数据
}
