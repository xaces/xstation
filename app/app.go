package app

import (
	"log"
	"xstation/app/access"
	"xstation/app/db"
	"xstation/app/ftp"
	"xstation/app/router"
	"xstation/configs"
	"xstation/controller/device"
)

// func loginServe() error {
// 	url := fmt.Sprintf("http://%s/station/online", configs.SuperAddress)
// 	address := fmt.Sprintf("%s:%d", configs.Default.Host, configs.Default.Port.Http)
// 	req := gin.H{"serveId": configs.Local.Id, "address": address}
// 	return xutils.HttpPost(url, req, nil)
// }

// Run 启动
func Run() error {
	if err := db.Run(&configs.Default.Sql); err != nil {
		return err
	}

	if configs.Default.Ftp.Enable {
		if err := ftp.Run(configs.PublicAbs(configs.Default.Public), &configs.Default.Ftp.Option); err != nil {
			return err
		}
		log.Printf("Xftp ListenAndServe at %s\n", configs.FtpAddr)
	}

	if configs.Default.Hook.Enable {
		device.Hooks(configs.Default.Hook.Options)
	}

	log.Printf("Xproto ListenAndServe at %s:%d\n", configs.Default.Host, configs.Default.Port.Access)
	if err := access.Run(configs.Default.Host, configs.Default.Port.Access); err != nil {
		return err
	}

	log.Printf("Http ListenAndServe at %s:%d\n", configs.Default.Host, configs.Default.Port.Http)
	router.Run(configs.Default.Port.Http)
	return nil
}

// Shutdown 停止
func Shutdown() error {
	access.Shutdown()
	// router.Stop()
	return nil
}
