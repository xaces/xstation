package app

import (
	"fmt"
	"log"
	"xstation/app/db"
	"xstation/app/ftp"
	"xstation/app/gw"
	"xstation/app/web"
	"xstation/configs"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils"
)

func loginServe() error {
	url := fmt.Sprintf("http://%s/station/online", configs.SuperAddress)
	address := fmt.Sprintf("%s:%d", configs.Default.Host, configs.Default.Port.Http)
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
	if configs.Default.Ftp.Enable {
		fopt := &ftp.Options{
			User: configs.Default.Ftp.User,
			Port: configs.Default.Ftp.Port,
			Pswd: configs.Default.Ftp.Pswd,
		}
		if configs.Default.Ftp.Url == "" {
			configs.Default.Ftp.Url = fmt.Sprintf("ftp://%s:%s@%s:%d", configs.Default.Ftp.User, configs.Default.Ftp.Pswd, configs.Default.Host, configs.Default.Ftp.Port)
		}
		log.Printf("Xftp ListenAndServe at %s\n", configs.Default.Ftp.Url)
		if err := ftp.New(fopt, configs.Default.Public); err != nil {
			return err
		}
	}
	log.Printf("Xproto ListenAndServe at %s:%d\n", configs.Default.Host, configs.Default.Port.Access)
	if err := gw.Start(configs.Default.Host, configs.Default.Port.Access); err != nil {
		return err
	}
	// if err := rpcxStart(s.RpcPort); err != nil {
	// 	return err
	// }
	// log.Printf("Rpc start on %d\n", s.RpcPort)
	// if err := hook.MqttStart(); err != nil {
	// 	return err
	// }
	log.Printf("Http ListenAndServe at %s:%d\n", configs.Default.Host, configs.Default.Port.Http)
	web.Start(configs.Default.Port.Http)
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
