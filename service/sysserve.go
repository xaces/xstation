package service

import (
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// ServeOpt 更新服务状态
type ServeOpt struct {
	Guids  []string `json:"guid"`
	Status int      `json:"status"`
}

func ServesDelete(guids []string) error {
	if _, err := orm.DbDeleteBy(&model.SysServe{}, "guid in (?)", guids); err != nil {
		return err
	}
	return nil
}
