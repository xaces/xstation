package configs

import (
	"github.com/BurntSushi/toml"
)

var (
	LocalId     string // 服务Id
	LocalIpAddr string // 服务IP
)

type tomlConfigure struct {
	Port struct {
		Api    uint16
		Access uint16
		Rpc    uint16
	}
	SQL struct {
		Name       string
		Address    string
		LiteDB     string
		Postgre    string
		StTableNum int
	}
	Superior struct {
		Address string
	}
}

// Default 所有配置参数
var (
	Default tomlConfigure
)

// Load 初始化配置参数
func Load(path *string) error {
	_, err := toml.DecodeFile(*path, &Default)
	return err
}
