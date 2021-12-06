package device

import (
	"xstation/model"
	"xstation/service"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"
)

type Online struct {
}

func (o *Online) ListHandler(c *gin.Context) {
	var param service.OnlinePage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var rows []model.DevOnline
	totalCount, err := orm.DbPage(&model.DevOnline{}, param.Where()).Scan(param.PageNum, param.PageSize, &rows)
	if err == nil {
		ctx.JSONOk().WriteData(gin.H{"total": totalCount, "rows": rows}, c)
		return
	}
	ctx.JSONWriteError(err, c)
}

func (o *Online) AddHandler(c *gin.Context) {
	var data model.DevOnline
	//获取参数
	if err := c.ShouldBind(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := orm.DbCreate(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}


func OnlineRouter(r *gin.RouterGroup) {
	on := Online{}
	r.POST("/online", on.AddHandler)
	r.GET("/online/list", on.ListHandler)
}