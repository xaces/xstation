package device

import (
	"xstation/model"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"
)

type Online struct {
}

func (o *Online) ListHandler(c *gin.Context) {
	var p Where
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data []model.DevOnline
	total, _ := orm.DbByWhere(&model.DevOnline{}, p.Online()).Find(&data)
	ctx.JSONWriteData(gin.H{"total": total, "data": data}, c)
}

func onlineRouter(r *gin.RouterGroup) {
	on := Online{}
	r.GET("/list", on.ListHandler)
}
