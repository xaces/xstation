package access

import (
	"xstation/controller/device"
	
	_ "github.com/xaces/xproto/jt"
	_ "github.com/xaces/xproto/ttx"
	_ "github.com/xaces/xproto/ho"
	"github.com/xaces/xproto"
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
