package serve

import (
	"log"
	"time"
	"xstation/app/mnger"
	"xstation/configs"
	"xstation/internal"
	"xstation/models"

	"github.com/wlgd/xutils"
	"github.com/wlgd/xutils/orm"
)

// localServe 本地服务
func localServe(name string) *models.XServer {
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
	s.Status = models.ServeStatusWork
	orm.DbUpdateModel(&s)
	return &s
}

func loginServer() error {
	return nil
}

// Run 启动
func Run() error {
	s := localServe("station")
	configs.LocalId = s.Guid
	configs.LocalIpAddr = xutils.PublicIPAddr()
	mnger.Serve.LoadOfDb()
	mnger.Dev.LoadOfDb()
	go xprotoStart(s.AccessPort)
	go func() {
		loginServer()
		ticker := time.NewTicker(time.Second * 60)
		for {
			<-ticker.C
			if err := loginServer(); err != nil {
				log.Println(err.Error())
			}
			mnger.Serve.UpdateAll() // 检测服务状态
		}
	}()
	return nil
}

// Stop 停止
func Stop() {
	xprotoStop()
}
