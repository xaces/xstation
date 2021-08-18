package sys

import (
	"xstation/mnger"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// serveOpt 更新服务状态
type serveOpt struct {
	Guids  []string `json:"guid"`
	Status int      `json:"status"`
}

func deleteServes(guids []string) error {
	if _, err := orm.DbDeleteBy(&model.SysServe{}, "guid in (?)", guids); err != nil {
		return err
	}
	mnger.Serve.Delete(guids)
	return nil
}
