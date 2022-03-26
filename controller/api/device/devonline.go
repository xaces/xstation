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
	var data []model.DevOnline
	total, _ := orm.DbByWhere(&model.DevOnline{}, param.Where()).Find(&data)
	ctx.JSONOk().WriteData(gin.H{"total": total, "data": data}, c)
}

func OnlineRouter(r *gin.RouterGroup) {
	on := Online{}
	r.GET("/online/list", on.ListHandler)
}
