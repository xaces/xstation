package ftp

import (
	"log"
	"os"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

var (
	ftp *server.Server
)

type Option struct {
	Port int
	User string
	Pswd string
}

// Run start server
func Run(path string, o *Option) error {
	os.MkdirAll(path, os.ModePerm)
	var perm = server.NewSimplePerm("test", "test")
	fopt := &server.ServerOpts{
		Name: "test ftpd",
		Factory: &filedriver.FileDriverFactory{
			RootPath: path,
			Perm:     perm,
		},
		Port: o.Port,
		Auth: &server.SimpleAuth{
			Name:     o.User,
			Password: o.Pswd,
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
