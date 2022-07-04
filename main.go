package main

import (
	"log"
	"os"
	"xstation/app"
	"xstation/configs"

	svc "github.com/kardianos/service"
)

type program struct {
}

func (p *program) Start(s svc.Service) error {
	if err := configs.Load("config.yaml"); err != nil {
		return err
	}
	return p.run()
}

func (p *program) run() error {
	return app.Run()
}

func (p *program) Stop(s svc.Service) error {
	return app.Shutdown()
}

func main() {
	svconf := &svc.Config{
		Name:        "xstation",
		DisplayName: "xstation",
		Description: "This is mdvr station application",
	}
	s, err := svc.New(&program{}, svconf)
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
