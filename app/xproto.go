package app

import (
	"strconv"
	"strings"

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

// XprotoStart 启动
func XprotoStart(addr string) error {
	as := strings.Split(addr, ":")
	port, _ := strconv.Atoi(as[1])
	xnotify := NewXNotify()
	s, err := xproto.NewServe(&xproto.Options{
		RequestTimeout:   50,
		RecvTimeout:      30,
		Port:             uint16(port),
		Host:             as[0],
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

// XprotoStop 停止
func XprotoStop() {
	_xproto.Release()
}
