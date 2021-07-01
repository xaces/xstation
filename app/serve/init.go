package serve

import (
	"errors"
	"log"
	"time"
	"xstation/app/manager"
	"xstation/configs"
	"xstation/internal"
	"xstation/models"

	"github.com/wlgd/xutils"
	"github.com/wlgd/xutils/orm"
	"github.com/wlgd/xutils/rpc"
)

var (
	isLoginServer = false
	rpcToken      = ""
)

// serveInit 初始化服务
func serveInit(name string) *models.XServer {
	var s models.XServer
	s.Role = models.ServeTypeLocal
	if err := orm.DbFirstBy(&s, "role = ?", s.Role); err != nil {
		s.Guid = internal.ServeId()
		orm.DbCreate(&s)
	}
	s.Name = name
	s.HttpPort = configs.Default.Port.Http
	s.RpcPort = configs.Default.Port.Rpc
	s.AccessPort = configs.Default.Port.Access
	s.Address = xutils.PublicIPAddr()
	s.Status = models.ServeStatusRunning
	orm.DbUpdateModel(&s)
	return &s
}

func loginServer() error {
	if isLoginServer {
		return rpc.KeepAlive(&rpc.KeepAliveArgs{
			ServeId:     configs.LocalId,
			Token:       rpcToken,
			UpdatedTime: time.Now()})
	}
	// if err := rpc.Connect("xvms", "127.0.0.1", 10000); err != nil {
	// 	return err
	// }
	// var reply rpc.LoginReply
	// if err := rpc.Login(&rpc.LoginArgs{
	// 	ServeId: global.LocalId,
	// 	Address: global.LocalIpAddr}, &reply); err != nil {
	// 	return err
	// }
	// rpcToken = reply.Token
	// isLoginServer = true
	return nil
}

// Run 启动
func Run() error {
	// 初始化上级服务
	if configs.Default.Superior.Address == "" {
		return errors.New("please set superior address firstly")
	}
	// 初始API服务
	s := serveInit("station")
	configs.LocalId = s.Guid
	configs.LocalIpAddr = s.Address
	manager.Serve.LoadOfDb()
	// 加载设备信息
	manager.Dev.LoadOfDb()
	// 初始化rpc
	go rpcxStart(s.RpcPort)
	go xprotoStart(s.AccessPort)
	go func() {
		loginServer()
		ticker := time.NewTicker(time.Second * 60)
		for {
			<-ticker.C
			if err := loginServer(); err != nil {
				log.Println(err.Error())
				isLoginServer = false
			}
		}
	}()
	return nil
}

// Stop 停止
func Stop() {
	xprotoStop()
	rpcxStop()
}
