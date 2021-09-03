package configs

import (
	"github.com/wlgd/xutils"
)

type localConfigure struct {
	Id            string // 服务Id
	IpAddr        string // 服务IP
	EffectiveTime int    // 有效时间
	MaxDevNumber  int
}

type ymlConfigure struct {
	Http struct {
		Host string `yaml:"host"`
		Port uint16 `yaml:"port"`
	} `yaml:"http"`

	Access struct {
		Host string `yaml:"host"`
		Port uint16 `yaml:"port"`
	} `yaml:"access"`

	SQL struct {
		Name    string `yaml:"name"`
		Address string `yaml:"address"`
	} `yaml:"sql"`

	Map struct {
		Name string `yaml:"name"`
		Key  string `yaml:"key"`
	} `yaml:"map"`
}

// Default 所有配置参数
var (
	Default      ymlConfigure
	Local        localConfigure
	SuperAddress string
)

// Load 初始化配置参数
func Load(license, path string) error {
	lice, err := xutils.LicenseRead(license)
	if err != nil {
		return err
	}
	Local.Id = lice.ServeGuid
	Local.EffectiveTime = lice.EffectiveTime
	Local.MaxDevNumber = lice.MaxNumber
	//TODO address
	SuperAddress = lice.Address
	return xutils.YMLConf(path, &Default)
}
