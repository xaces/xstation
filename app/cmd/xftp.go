package main

import (
	"log"
	"os"
	"xstation/app/ftp"

	"github.com/wlgd/xutils"

	svc "github.com/kardianos/service"
)

type ymlconfigure struct {
	Ftp    ftp.Options `yaml:"ftp"`
	Public string      `yaml:"public"`
}

var (
	yml ymlconfigure
)

type program struct {
}

func (p *program) Start(s svc.Service) error {

	if err := xutils.YMLConf(".config.yml", &yml); err != nil {
		log.Fatalln(err)
	}
	return p.run()
}

func (p *program) run() error {
	return ftp.New(&yml.Ftp, yml.Public)
}

func (p *program) Stop(s svc.Service) error {
	ftp.Shutdown()
	return nil
}

func main() {
	svvconfig := &svc.Config{
		Name:        "xvms.ftp",
		DisplayName: "xftp",
		Description: "This is xvms ftp server",
	}
	s, err := svc.New(&program{}, svvconfig)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err = s.Install()
		} else if os.Args[1] == "uninstall" {
			err = s.Uninstall()
		}
		log.Println(err)
		return
	}
	err = s.Run()
	if err != nil {
		log.Fatal(err)
	}
}
