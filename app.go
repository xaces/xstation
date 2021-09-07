package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"xstation/app"
	"xstation/configs"
	"xstation/controller"
	"xstation/mnger"
	"xstation/model"
	"xstation/service"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils"
	"github.com/wlgd/xutils/orm"
)

var (
	hs *http.Server
)

// loadData 本地数据
func loadData() error {
	var devices []model.Device
	orm.DbFind(&devices)
	if len(devices) > configs.Local.MaxDevNumber {
		return fmt.Errorf("please do not modify the database")
	}
	mnger.Dev.Set(devices)
	return nil
}

func loginServe() error {
	url := fmt.Sprintf("http://%s/station/online", configs.SuperAddress)
	address := fmt.Sprintf("%s:%d", configs.Default.Http.Host, configs.Default.Http.Port)
	req := gin.H{"serveId": configs.Local.Id, "address": address}
	return xutils.HttpPost(url, req, nil)
}

// Run 启动
func AppRun() error {
	if err := service.Init(); err != nil {
		return err
	}
	if err := loadData(); err != nil {
		return err
	}
	loginServe()
	log.Printf("Xproto ListenAndServe at %s:%d\n", configs.Default.Access.Host, configs.Default.Access.Port)
	if err := app.XprotoStart(configs.Default.Access.Host, configs.Default.Access.Port); err != nil {
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
	hs = controller.NewServer(configs.Default.Http.Port)
	return nil
}

// AppShutdown 停止
func AppShutdown() {
	app.XprotoStop()
	// rpcxStop()
	// hook.MqttStop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := hs.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
}
