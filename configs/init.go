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
	HttpAddr   string //
	AccessAddr string //
	SQL        struct {
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
	Default      tomlConfigure
	Local        localConfigure
	SuperAddress string
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
	SuperAddress = lice.Address
	if _, err := toml.DecodeFile(*path, &Default); err != nil {
		return err
	}
	return err
}
