package ftp

import (
	"errors"
	"log"
	"os"

	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
)

var (
	ftp *server.Server
)

type Option struct {
	Enable bool
	Port   int
	User   string
	Pswd   string
	Path   string
}

// Run start server
func Run(o *Option) error {
	if !o.Enable {
		return errors.New("Xftp disable")
	}
	os.MkdirAll(o.Path, os.ModePerm)
	var perm = server.NewSimplePerm("test", "test")
	fopt := &server.ServerOpts{
		Name: "test ftpd",
		Factory: &filedriver.FileDriverFactory{
			RootPath: o.Path,
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
