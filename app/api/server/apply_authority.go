package server

import (
	"errors"
	"xstation/configs"
	"xstation/internal"
	"xstation/pkg/ctx"
	"xstation/pkg/utils"

	"github.com/gin-gonic/gin"
)

type applyAuth struct {
	Identity string `json:"identity"` // 身份令牌
	Username string `json:"username"` // 登录
	Password string `json:"password"` // 密码
}

// ApplyAuthHandler 身份授权
func ApplyAuthHandler(c *gin.Context) {
	var param applyAuth
	if err := c.ShouldBindJSON(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data internal.LicensingInf
	if err := utils.HttpPost(configs.Default.Superior.Address+"/station/applyAuth", &param, &data); err != nil {
		ctx.JSONWriteError(errors.New("apply authority failed"), c)
		return
	}
	internal.WriteLicences(&data)
	ctx.JSONOk().WriteTo(c)
}
