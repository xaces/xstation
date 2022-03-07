package app

import (
	"log"
	"xstation/app/access"
	"xstation/app/db"
	"xstation/app/ftp"
	"xstation/app/router"
	"xstation/configs"
	"xstation/controller/device"
	"xstation/middleware"
)

// func loginServe() error {
// 	url := fmt.Sprintf("http://%s/station/online", configs.SuperAddress)
// 	address := fmt.Sprintf("%s:%d", configs.Default.Host, configs.Default.Port.Http)
// 	req := gin.H{"serveId": configs.Local.Id, "address": address}
// 	return xutils.HttpPost(url, req, nil)
// }

func thirdHooks() error {
	if configs.Default.RdMQ.Enable {
		switch configs.Default.RdMQ.Name {
		case "nats":
			n := &middleware.Nats{}
			if err := n.Start(configs.Default.RdMQ.NatsOption); err != nil {
				return err
			}
			device.Handler.Handle(n)
		}
	}
	return nil
}

// Run 启动
func Run() error {
	if err := db.Init(&configs.Default.Sql); err != nil {
		return err
	}
	if configs.Default.Ftp.Enable {
		if err := ftp.Run(&configs.Default.Ftp.Options); err != nil {
			return err
		}
		log.Printf("Xftp ListenAndServe at %s\n", configs.FtpAddr)
	}
	if err := device.Handler.Run(configs.Default.MsgProc); err != nil {
		return err
	}
	if err := thirdHooks(); err != nil {
		return err
	}
	log.Printf("Xproto ListenAndServe at %s:%d\n", configs.Default.Host, configs.Default.Port.Access)
	if err := access.Start(configs.Default.Host, configs.Default.Port.Access); err != nil {
		return err
	}
	log.Printf("Http ListenAndServe at %s:%d\n", configs.Default.Host, configs.Default.Port.Http)
	router.Run(configs.Default.Port.Http)
	return nil
}

// Shutdown 停止
func Shutdown() error {
	access.Stop()
	router.Stop()
	device.Handler.Stop()
	return nil
}
