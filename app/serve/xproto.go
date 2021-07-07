package serve

import (
	"github.com/wlgd/xproto"
	"github.com/wlgd/xproto/ho"
	"github.com/wlgd/xproto/jt"
	"github.com/wlgd/xproto/ttx"
)

func protocolAdapter(b []byte) xproto.InterfaceProto {
	if _, err := ho.IsValidProto(b); err == nil {
		return &ho.ProtoImpl{
			SubAlarmStatus: 0xffff,
			SubStatus:      0xffff,
		}
	}
	if _, err := ttx.IsValidProto(b); err == nil {
		return &ttx.ProtoImpl{}
	}
	if _, err := jt.IsValidProto(b); err == nil {
		return &jt.ProtoImpl{}
	}
	return nil
}

var (
	_xproto *xproto.Serve = nil
)

// xprotoStart 启动
func xprotoStart(port uint16) error {
	xnotify := NewXNotify()
	s, err := xproto.NewServe(xproto.Options{
		RequestTimeout:   50,
		RecvTimeout:      30,
		Port:             port,
		Adapter:          protocolAdapter,
		LinkAccessNotify: xnotify.AccessHandler,
		StatusNotify:     xnotify.StatusHandler,
		AlarmNotify:      xnotify.AlarmHandler,
		EventNotify:      xproto.LogEvent,
		AVFrameNotify:    xproto.LogAVFrame,
		RawNotify:        xproto.LogRawFrame,
	})
	_xproto = s
	return err
}

// xprotoStop 停止
func xprotoStop() {
	_xproto.Release()
}
