package app

import (
	"fmt"
	"log"
	"xstation/app/access"
	"xstation/app/db"
	"xstation/app/ftp"
	"xstation/app/router"
	"xstation/configs"
	"xstation/controller/device"
	"xstation/entity/cache"
	"xstation/entity/task"

	"github.com/wlgd/xutils"
)

func getLocalVehicle() error {
	if configs.Default.Super.Api == "" {
		return nil // 测试
	}
	// 获取设备信息
	url := fmt.Sprintf("%s/devices/%s", configs.Default.Super.Api, configs.Default.Guid)
	var vehis []cache.Vehicle
	if err := xutils.HttpGet(url, &vehis); err != nil {
		return err
	}
	for _, v := range vehis {
		cache.NewDevice(v)
	}
	return nil
}

// Run 启动
func Run() error {
	conf := configs.Default
	if err := getLocalVehicle(); err != nil {
		return err
	}
	task.Timer.Run() // 启动定时任务
	if err := db.Run(&conf.Sql); err != nil {
		return err
	}
	if err := ftp.Run(&conf.Ftp); err == nil {
		configs.FtpAddr = fmt.Sprintf("ftp://%s:%s@%s:%d", conf.Ftp.User, conf.Ftp.Pswd, conf.Host, conf.Ftp.Port)
		log.Printf("Xftp ListenAndServe at %s\n", configs.FtpAddr)
	}
	if conf.Hook.Enable {
		device.NewHooks(conf.Hook.Options)
	}
	if err := access.Run(conf.Host, conf.Port.Access); err != nil {
		return err
	}
	router.Run(conf.Port.Http)
	return nil
}

// Shutdown 停止
func Shutdown() error {
	access.Shutdown()
	// router.Stop()
	return nil
}
