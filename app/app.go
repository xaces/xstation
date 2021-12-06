package app

import (
	"fmt"
	"log"
	"xstation/app/db"
	"xstation/app/gw"
	"xstation/app/web"
	"xstation/configs"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils"
)

func loginServe() error {
	url := fmt.Sprintf("http://%s/station/online", configs.SuperAddress)
	address := fmt.Sprintf("%s:%d", configs.Default.Http.Host, configs.Default.Http.Port)
	req := gin.H{"serveId": configs.Local.Id, "address": address}
	return xutils.HttpPost(url, req, nil)
}

// Run 启动
func Run() error {
	if err := db.Init(); err != nil {
		return err
	}
	// if err := loginServe(); err != nil {
	// 	return err
	// }
	log.Printf("Xproto ListenAndServe at %s:%d\n", configs.Default.Access.Host, configs.Default.Access.Port)
	if err := gw.Start(configs.Default.Access.Host, configs.Default.Access.Port); err != nil {
		return err
	}
	// if err := rpcxStart(s.RpcPort); err != nil {
	// 	return err
	// }
	// log.Printf("Rpc start on %d\n", s.RpcPort)
	// if err := hook.MqttStart(); err != nil {
	// 	return err
	// }
	log.Printf("Http ListenAndServe at %s:%d\n", configs.Default.Http.Host, configs.Default.Http.Port)
	web.Start(configs.Default.Http.Port)
	return nil
}

// Shutdown 停止
func Shutdown() error {
	gw.Stop()
	web.Stop()
	// rpcxStop()
	// hook.MqttStop()
	return nil
}
