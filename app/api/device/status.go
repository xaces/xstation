package device

import (
	"xstation/models"
	"xstation/pkg/ctx"
	"xstation/pkg/orm"
	"xstation/service"

	"github.com/gin-gonic/gin"
)

type Status struct {
}

func (o *Status) ListHandler(c *gin.Context) {
	var param statusPage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var rows []models.XStatus
	totalCount, err := orm.DbPage(service.GetXStatusModel(param.DeviceId), param.Where()).Scan(param.PageNum, param.PageSize, &rows)
	if err == nil {
		ctx.JSONOk().WriteData(gin.H{
			"total": totalCount,
			"rows":  rows}, c)
		return
	}
	ctx.JSONWriteError(err, c)
}

// GetHandler 获取指定id
func (o *Status) GetHandler(c *gin.Context) {
	param := statusGet{}
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	data := service.GetXStatusModel(param.DeviceId)
	if err := orm.DbFirstById(param.StatusId, data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(gin.H{"data": data}, c)
}
