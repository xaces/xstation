package configs

import (
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
	Guid    string
	Port    struct {
		Http   uint16
		Access uint16
	}
	Super struct {
		Api string
	}
	Sql  db.Option
	Ftp  ftp.Option
	Hook struct {
		Enable  bool
		Options []hook.Option
	}
}

// Default 所有配置参数
var (
	Default configure
	License xutils.License
	FtpAddr string
	absDir  string
	MsgProc int // 0--测试模式 1--正常模式
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
	if err := xutils.YAMLFile(PublicAbs(path), &Default); err != nil {
		return err
	}
	if lice, err := xutils.LicenseRead(PublicAbs(Default.License)); err != nil {
		return err
	} else {
		License = *lice
	}
	Default.Ftp.Path = PublicAbs(Default.Public)
	return nil
}
