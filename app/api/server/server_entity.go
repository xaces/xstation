package server

import (
	"errors"
	"xstation/configs"
	"xstation/internal"

	"github.com/wlgd/xutils"
)

type applyAuth struct {
	Identity string `json:"identity"` // 身份令牌
	Username string `json:"username"` // 登录
	Password string `json:"password"` // 密码
}

// tryApplyAuth 尝试申请授权
func tryApplyAuth(param *applyAuth) error {
	var data internal.LicensingInf
	if err := xutils.HttpPost(configs.Default.Superior.Address+"/station/applyAuth", param, &data); err != nil {
		return errors.New("apply authority failed")
	}
	internal.WriteLicences(&data)
	return nil
}

// updateStatus 更新服务状态
type updateStatus struct {
	Guid   string `json:"guid"`
	Status int    `json:"status"`
}
