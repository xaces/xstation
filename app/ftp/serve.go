package ftp

import (
	"log"
	"os"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

type Options struct {
	Port uint16
	User string
	Pswd string
}

var (
	ftp *server.Server
)

// NewFtp start server
func New(opt *Options, root string) error {
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

func Shutdown() {
	if ftp != nil {
		ftp.Shutdown()
	}
}
