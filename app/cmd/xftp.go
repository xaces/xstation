package main

import (
	"log"
	"os"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
	"github.com/wlgd/xutils"

	svc "github.com/kardianos/service"
)

type Options struct {
	Port uint16 `yaml:"port"`
	User string `yaml:"user"`
	Pswd string `yaml:"pswd"`
}

type ymlconfigure struct {
	Ftp    Options `yaml:"ftp"`
	Public string  `yaml:"public"`
}

var (
	ftp *server.Server
	yml ymlconfigure
)

// newFtp start server
func newFtp(opt *Options, root string) error {
	os.MkdirAll(root, os.ModePerm)
	var perm = server.NewSimplePerm("test", "test")
	fopt := &server.ServerOpts{
		Name: "test ftpd",
		Factory: &filedriver.FileDriverFactory{
			RootPath: root,
			Perm:     perm,
		},
		Port: int(opt.Port),
		Auth: &server.SimpleAuth{
			Name:     opt.User,
			Password: opt.Pswd,
		},
		Logger: new(server.DiscardLogger),
	}
	ftp = server.NewServer(fopt)
	go func() {
		log.Fatalln(ftp.ListenAndServe())
	}()
	return nil
}

func ftpShutdown() {
	if ftp != nil {
		ftp.Shutdown()
	}
}

type program struct {
}

func (p *program) Start(s svc.Service) error {

	if err := xutils.YMLConf(".config.yml", &yml); err != nil {
		log.Fatalln(err)
	}
	return p.run()
}

func (p *program) run() error {
	return newFtp(&yml.Ftp, yml.Public)
}

func (p *program) Stop(s svc.Service) error {
	ftpShutdown()
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
