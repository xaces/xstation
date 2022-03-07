package configs

import (
	"fmt"
	"xstation/app/ftp"

	"xstation/app/db"

	"xstation/middleware"

	"github.com/BurntSushi/toml"
	"github.com/wlgd/xutils"
)

type localConfigure struct {
	Id            string // 服务Id
	IpAddr        string // 服务IP
	EffectiveTime int    // 有效时间
	MaxDevNumber  int
}

type configure struct {
	Host    string
	License string
	Public  string
	MsgProc string
	Port    struct {
		Http   uint16
		Access uint16
	}
	Sql db.Options
	Ftp struct {
		Enable bool
		ftp.Options
	}
	Map struct {
		Name string
		Key  string
	}
	RdMQ struct {
		Enable      bool
		Name        string
		middleware.NatsOption
	}
	RdHttp struct {
		Enable bool
		Online string
		Alarm  string
		Status string
		Event  string
	}
}

// Default 所有配置参数
var (
	Default      configure
	Local        localConfigure
	FtpAddr      string
	SuperAddress string
)

// Load 初始化配置参数
func Load(path string) error {
	if _, err := toml.DecodeFile(path, &Default); err != nil {
		return err
	}
	FtpAddr = fmt.Sprintf("ftp://%s:%s@%s:%d", Default.Ftp.User, Default.Ftp.Pswd, Default.Host, Default.Ftp.Port)
	lice, err := xutils.LicenseRead(Default.License)
	if err != nil {
		return err
	}
	Local.Id = lice.ServeGuid
	Local.EffectiveTime = lice.EffectiveTime
	Local.MaxDevNumber = lice.MaxNumber
	//TODO address
	SuperAddress = lice.Address
	return nil
}
