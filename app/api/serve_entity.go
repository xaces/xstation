package api

import (
	"errors"
	"xstation/app/mnger"
	"xstation/configs"
	"xstation/internal"
	"xstation/model"

	"github.com/wlgd/xutils"
	"github.com/wlgd/xutils/orm"
)

type applyAuth struct {
	Identity string `json:"identity"` // 身份令牌
	Username string `json:"username"` // 登录
	Password string `json:"password"` // 密码
}

// tryApplyAuth 尝试申请授权
func tryApplyAuth(param *applyAuth) error {
	var data internal.LicensingInf
	if err := xutils.HttpPost(configs.SuperiorAddress+"/station/applyAuth", param, &data); err != nil {
		return errors.New("apply authority failed")
	}
	internal.WriteLicences(&data)
	return nil
}

// serveOpt 更新服务状态
type serveOpt struct {
	Guids  []string `json:"guid"`
	Status int      `json:"status"`
}

func deleteServes(guids []string) error {
	if _, err := orm.DbDeleteBy(&model.Serve{}, "guid in (?)", guids); err != nil {
		return err
	}
	mnger.Serve.Delete(guids)
	return nil
}
