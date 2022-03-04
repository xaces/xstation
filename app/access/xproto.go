package access

import (
	"xstation/controller/device"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xproto/ho"
	"github.com/wlgd/xproto/jt"
	"github.com/wlgd/xproto/ttx"
)

func protocolAdapter(b []byte) xproto.InterfaceProtocol {
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

// Start 启动
func Start(msgproc, host string, port uint16) error {
	device.Handler.Run(msgproc)
	s, err := xproto.NewServe(&xproto.Options{
		RequestTimeout: 50,
		RecvTimeout:    30,
		Port:           uint16(port),
		Host:           host,
		Adapter:        protocolAdapter,
		AccessNotify:   device.AccessHandler,
		DroppedNotify:  device.DroppedHandler,
		StatusNotify:   device.StatusHandler,
		AlarmNotify:    device.AlarmHandler,
		EventNotify:    device.EventHandler,
		AVFrameNotify:  xproto.LogAVFrame,
		RawNotify:      xproto.LogRawFrame,
	})
	_xproto = s
	return err
}

// Stop 停止
func Stop() {
	_xproto.Release()
}
