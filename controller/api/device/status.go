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
	var data []model.DevStatus
   	total, _ := orm.DbPage(m, param.Where()).Find(param.PageNum, param.PageSize, &data)
	ctx.JSONOk().Write(gin.H{"total": total, "data": data}, c)
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
