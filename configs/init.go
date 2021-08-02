package configs

import (
	"github.com/BurntSushi/toml"
	"github.com/wlgd/xutils"
)

type localConfigure struct {
	Id            string // 服务Id
	IpAddr        string // 服务IP
	EffectiveTime int    // 有效时间
}

type tomlConfigure struct {
	Host string //
	Port struct {
		Http   uint16
		Access uint16
		Rpc    uint16
	}
	SQL struct {
		Name    string
		Address string
	}
	Map struct {
		Name string
		Key  string
	}
}

// Default 所有配置参数
var (
	Default         tomlConfigure
	Local           localConfigure
	SuperiorAddress string
)

// Load 初始化配置参数
func Load(licences, path *string) error {
	lice, err := xutils.ReadLicences(*licences)
	if err != nil {
		return err
	}
	Local.Id = lice.ServeGuid
	Local.EffectiveTime = lice.EffectiveTime
	//TODO address
	SuperiorAddress = lice.Address
	if _, err := toml.DecodeFile(*path, &Default); err != nil {
		return err
	}
	if Default.Host == "" {
		Default.Host = xutils.PublicIPAddr()
	}
	return err
}
