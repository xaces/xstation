package device

import (
	"xstation/internal/errors"
	"xstation/model"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils/ctx"
)

type Online struct {
}

func (o *Online) ListHandler(c *gin.Context) {
	var p Where
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if p.isDeviceNoInvalid() {
		ctx.JSONWriteError(errors.InvalidDeviceNo, c)
		return
	}
	var data []model.DevOnline
	total, _ := p.Online().Model(&model.DevOnline{}).Find(&data)
	ctx.JSONWriteData(gin.H{"total": total, "data": data}, c)
}

func onlineRouter(r *gin.RouterGroup) {
	on := Online{}
	r.GET("/list", on.ListHandler)
}
