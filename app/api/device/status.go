package device

import (
	"xstation/app/mnger"
	"xstation/model"

	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xutils/ctx"

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
	_, m := mnger.Dev.GetModel(param.DeviceNo)
	var rows []model.DevStatus
	totalCount, err := orm.DbPage(m, param.Where()).Scan(param.PageNum, param.PageSize, &rows)
	if err == nil {
		ctx.JSONOk().Write(gin.H{"total": totalCount, "rows": rows}, c)
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
	_, data := mnger.Dev.GetModel(param.DeviceNo)
	if err := orm.DbFirstById(data, param.StatusId); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(data, c)
}
