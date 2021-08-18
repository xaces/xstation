package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"xstation/configs"

	"github.com/urfave/cli/v2"
)

// 启动流程
// 1、初始化数据库
// 2、获取向中心服务配置信息
func main() {
	app := &cli.App{
		Name:        "mdvr's workstation",
		Version:     "0.3.0",
		Description: "This is mdvr access application",
		Authors: []*cli.Author{{
			Name:  "don.wang",
			Email: "wanguandong@126.com",
		},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "c",
				Usage: "config file",
				Value: "./conf/config.toml",
			},
			&cli.StringFlag{
				Name:  "ces",
				Usage: "licences file",
				Value: "localtest.lice",
			},
		},
		Before: func(c *cli.Context) error {
			return configs.Load(c.String("ces"), c.String("c"))
		},
		Action: func(c *cli.Context) error {
			if err := AppRun(); err != nil {
				return err
			}
			quit := make(chan os.Signal)
			// kill (no param) default send syscanll.SIGTERM
			// kill -2 is syscall.SIGINT
			// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit
			AppShutdown()
			log.Println("AppShutdown done")
			return nil
		},
	}
	log.Println(app.Run(os.Args))
}
