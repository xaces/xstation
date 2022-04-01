package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"xstation/app/db"
	"xstation/app/ftp"
	"xstation/entity/hook"

	"github.com/wlgd/xutils"
)

type configure struct {
	Host    string
	License string
	Public  string
	Port    struct {
		Http   uint16
		Access uint16
	}
	Sql db.Option
	Ftp struct {
		Enable bool
		Option ftp.Option
	}
	Hook struct {
		Enable  bool
		Options []hook.Option
	}
}

// Default 所有配置参数
var (
	Default      configure
	License      xutils.License
	FtpAddr      string
	SuperAddress string
	absDir       string
)

func PublicAbs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return absDir + "/" + path
}

// Load 初始化配置参数
func Load(path string) error {
	absDir = filepath.Dir(os.Args[0])
	if err := xutils.YMLConf(PublicAbs(path), &Default); err != nil {
		return err
	}
	if lice, err := xutils.LicenseRead(PublicAbs(Default.License)); err != nil {
		return err
	} else {
		License = *lice
	}
	FtpAddr = fmt.Sprintf("ftp://%s:%s@%s:%d", Default.Ftp.Option.User, Default.Ftp.Option.Pswd, Default.Host, Default.Ftp.Option.Port)
	return nil
}
