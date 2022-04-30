package access

import (
	"log"
	"xstation/controller/device"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xproto/ho"
	"github.com/wlgd/xproto/jt"
	"github.com/wlgd/xproto/ttx"
)

func protocolAdapter(b []byte) xproto.IClient {
	if c := ho.NewClient(b, &ho.Options{SubAlarmStatus: 0xffff, SubStatus: 0xffff}); c != nil {
		return c
	}
	if c := ttx.NewClient(b); c != nil {
		return c
	}
	if c := jt.NewClient(b); c != nil {
		return c
	}
	return nil
}

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
		Adapter:        protocolAdapter,
	}); err != nil {
		return
	}
	device.Handler.Disptah()
	s.Handle.Access = device.AccessHandler
	s.Handle.Dropped = device.DroppedHandler
	s.Handle.Status = device.StatusHandler
	s.Handle.Alarm = device.AlarmHandler
	s.Handle.Event = device.EventHandler
	s.Handle.Frame = xproto.LogFrame
	go s.ListenTCPAndServe()
	log.Printf("Xproto ListenAndServe at %s:%d\n", host, port)
	return
}

// Stop 停止
func Shutdown() {
	if s != nil {
		s.Release()
	}
	device.Handler.Stop()
}
