package serve

import (
	"log"
	"xstation/app/mnger"
	"xstation/configs"
	"xstation/internal"
	"xstation/model"
	"xstation/service"

	"github.com/wlgd/xutils"
	"github.com/wlgd/xutils/orm"
)

// localServe 本地服务
func localServe(name string) *model.Serve {
	var s model.Serve
	if err := orm.DbFirstById(&s, 1); err != nil {
		s.Guid = internal.ServeId()
		orm.DbCreate(&s)
	}
	s.Name = name
	s.HttpPort = configs.Default.Port.Http
	s.RpcPort = configs.Default.Port.Rpc
	s.AccessPort = configs.Default.Port.Access
	s.Status = 1
	orm.DbUpdateModel(&s)
	var devices []model.Device
	if err := orm.DbFind(&devices); err == nil {
		mnger.Dev.Set(devices)
	}
	return &s
}

// Run 启动
func Run() error {
	if err := service.Init(); err != nil {
		return err
	}
	s := localServe("station")
	configs.Local.Id = s.Guid
	configs.Local.IpAddr = xutils.PublicIPAddr()
	if err := xprotoStart(s.AccessPort); err != nil {
		return err
	}
	log.Printf("Xproto start on %d\n", s.AccessPort)
	// if err := rpcxStart(s.RpcPort); err != nil {
	// 	return err
	// }
	// log.Printf("Rpc start on %d\n", s.RpcPort)
	// if err := hook.MqttStart(); err != nil {
	// 	return err
	// }
	return nil
}

// Stop 停止
func Stop() {
	xprotoStop()
	// rpcxStop()
	// hook.MqttStop()
}
