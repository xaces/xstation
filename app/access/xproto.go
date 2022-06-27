package access

import (
	"xstation/controller/device"

	"github.com/wlgd/xproto"
	_ "github.com/wlgd/xproto/ho"
	_ "github.com/wlgd/xproto/jt"
	_ "github.com/wlgd/xproto/ttx"
)

var (
	s *xproto.Server = nil
)

// Start 启动
func Run(host string, port uint16) (err error) {
	if s, err = xproto.NewServer(&xproto.Options{
		RequestTimeout: 50,
		RecvTimeout:    30,
		Port:           uint16(port),
		Host:           host,
	}); err != nil {
		return
	}
	device.Handler.Disptah()
	s.Handler.Access = device.AccessHandler
	s.Handler.Dropped = device.DroppedHandler
	s.Handler.Status = device.StatusHandler
	s.Handler.Alarm = device.AlarmHandler
	s.Handler.Event = device.EventHandler
	s.Handler.Frame = xproto.LogFrame
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
