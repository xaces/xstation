package main

import (
	"fmt"
	"log"
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
	req := gin.H{"serveId": configs.Local.Id, "address": configs.Default.HttpAddr}
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
	log.Printf("Xproto ListenAndServe at %s\n", configs.Default.AccessAddr)
	if err := app.XprotoStart(configs.Default.AccessAddr); err != nil {
		return err
	}
	// if err := rpcxStart(s.RpcPort); err != nil {
	// 	return err
	// }
	// log.Printf("Rpc start on %d\n", s.RpcPort)
	// if err := hook.MqttStart(); err != nil {
	// 	return err
	// }
	log.Printf("Http ListenAndServe at %s\n", configs.Default.HttpAddr)
	go controller.NewServer(configs.Default.HttpAddr).ListenAndServe()
	return nil
}

// AppShutdown 停止
func AppShutdown() {
	app.XprotoStop()
	// rpcxStop()
	// hook.MqttStop()
}
