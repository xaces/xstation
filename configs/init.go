package configs

import (
	"errors"

	"github.com/BurntSushi/toml"
)

type localConfigure struct {
	Id     string // 服务Id
	IpAddr string // 服务IP
}

type tomlConfigure struct {
	Port struct {
		Http   uint16
		Access uint16
		Rpc    uint16
	}
	SQL struct {
		Name    string
		Address string
	}
	Superior struct {
		Address string
	}
	Map struct {
		Name string
		Key  string
	}
}

// Default 所有配置参数
var (
	Default tomlConfigure
	Local   localConfigure
)

// Load 初始化配置参数
func Load(path *string) error {
	_, err := toml.DecodeFile(*path, &Default)
	if Default.Superior.Address == "" {
		return errors.New("please set superior address firstly")
	}
	return err
}
