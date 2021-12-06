package device

import (
	"xstation/app/mnger"
	"xstation/model"
	"xstation/service"

	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xutils/ctx"

	"github.com/gin-gonic/gin"
)

type Status struct {
}

func (o *Status) ListHandler(c *gin.Context) {
	var param service.StatusPage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	_, m := mnger.Devs.Model(param.DeviceNo)
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
	param := service.StatusGet{}
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	_, data := mnger.Devs.Model(param.DeviceNo)
	if err := orm.DbFirstById(data, param.StatusId); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(data, c)
}

func StatusRouter(r *gin.RouterGroup) {
	s := Status{}
	r.GET("/status/list", s.ListHandler)
	r.GET("/status", s.GetHandler)
}
