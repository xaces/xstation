package configs

import (
	"github.com/BurntSushi/toml"
)

type configure struct {
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
	Default configure
)

// Load 初始化配置参数
func Load(path *string) error {
	_, err := toml.DecodeFile(*path, &Default)
	return err
}
