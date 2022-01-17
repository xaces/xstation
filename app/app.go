package app

import (
	"fmt"
	"log"
	"strings"
	"xstation/app/access"
	"xstation/app/ftp"
	"xstation/app/mnger"
	"xstation/app/router"
	"xstation/configs"
	"xstation/internal"
	"xstation/service"
)

// func loginServe() error {
// 	url := fmt.Sprintf("http://%s/station/online", configs.SuperAddress)
// 	address := fmt.Sprintf("%s:%d", configs.Default.Host, configs.Default.Port.Http)
// 	req := gin.H{"serveId": configs.Local.Id, "address": address}
// 	return xutils.HttpPost(url, req, nil)
// }

// Run 启动
func Run() error {
	if err := service.Init(); err != nil {
		return err
	}
	if err := mnger.Init(); err != nil {
		return err
	}
	if configs.Default.Ftp.Enable {
		log.Printf("Xftp ListenAndServe at %s\n", configs.Default.Ftp.Address)
		port, user, pswd := internal.StringParseFtpUri(configs.Default.Ftp.Address)
		if err := ftp.Run(port, user, pswd, configs.Default.Public); err != nil {
			return err
		}
		if strings.Contains(configs.Default.Ftp.Address, "127.0.0.1") {
			configs.Default.Ftp.Address = fmt.Sprintf("ftp://%s:%s@%s:%d", user, pswd, configs.Default.Host, port)
		}
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
	return nil
}
