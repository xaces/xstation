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
	"xstation/model"

	"github.com/wlgd/xutils"
	"github.com/wlgd/xutils/orm"
)

func getLocalDeivce() (err error) {
	var data []cache.DeviceInfo
	// 获取设备信息
	if configs.Default.Super.Api == "" {
		err = orm.DB().Model(&model.Device{}).Scan(&data).Error
	} else {
		configs.MsgProc = 1
		url := fmt.Sprintf("%s/devices/%s", configs.Default.Super.Api, configs.Default.Guid)
		err = xutils.HttpGet(url, &data)
	}
	for _, v := range data {
		cache.NewDevice(v)
	}
	return
}

// Run 启动
func Run() error {
	conf := configs.Default
	task.Timer.Run() // 启动定时任务
	if err := db.Run(&conf.Sql); err != nil {
		return err
	}
	if err := getLocalDeivce(); err != nil {
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
