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
	s *xproto.Server = nil
)

// Start 启动
func Run(host string, port uint16) (err error) {
	opt := &xproto.Options{
		RequestTimeout: 50,
		RecvTimeout:    30,
		Port:           uint16(port),
		Host:           host,
		Adapter:        protocolAdapter,
	}
	if s, err = xproto.NewServer(opt); err != nil {
		return
	}
	device.Handler.Disptah()
	s.Handle.Access = device.AccessHandler
	s.Handle.Dropped = device.DroppedHandler
	s.Handle.Status = device.StatusHandler
	s.Handle.Alarm = device.AlarmHandler
	s.Handle.Event = device.EventHandler
	go s.ListenTCPAndServe()
	return
}

// Stop 停止
func Shutdown() {
	if s != nil {
		s.Release()
	}
	device.Handler.Stop()
}
