package serve

import (
	"fmt"
	"log"
	"xstation/app/mnger"
	"xstation/configs"
	"xstation/model"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils"
	"github.com/wlgd/xutils/orm"
)

// localData 本地数据
func localData() {
	var devices []model.Device
	if err := orm.DbFind(&devices); err == nil {
		mnger.Dev.Set(devices)
	}
}

func loginServe() error {
	url := fmt.Sprintf("http://%s/stationLogin", configs.SuperiorAddress)
	addr := fmt.Sprintf("%s:%d", configs.Default.Host, configs.Default.Port.Http)
	req := gin.H{"serveId": configs.Local.Id, "address": addr}
	return xutils.HttpPost(url, req, nil)
}

// Run 启动
func Run() error {
	localData()
	loginServe()
	if err := xprotoStart(configs.Default.Port.Access); err != nil {
		return err
	}
	log.Printf("Xproto start on %d\n", configs.Default.Port.Access)
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
