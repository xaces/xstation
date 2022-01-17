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

// Run start server
func Run(port int, user, pswd, root string) error {
	os.MkdirAll(root, os.ModePerm)
	var perm = server.NewSimplePerm("test", "test")
	fopt := &server.ServerOpts{
		Name: "test ftpd",
		Factory: &filedriver.FileDriverFactory{
			RootPath: root,
			Perm:     perm,
		},
		Port: port,
		Auth: &server.SimpleAuth{
			Name:     user,
			Password: pswd,
		},
		Logger: new(server.DiscardLogger),
	}
	ftp = server.NewServer(fopt)
	go func() {
		log.Fatalln(ftp.ListenAndServe())
	}()
	return nil
}

func Stop() {
	if ftp != nil {
		ftp.Shutdown()
	}
}
