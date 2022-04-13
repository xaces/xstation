package device

import (
	"xstation/entity/cache"

	"github.com/wlgd/xutils/ctx"

	"github.com/gin-gonic/gin"
)

type Device struct {
}

func (o *Device) ListHandler(c *gin.Context) {
	data := cache.DeviceList()
	ctx.JSONOk().Write(gin.H{"total": len(data), "data": data}, c)
}

func DeviceRouter(r *gin.RouterGroup) {
	d := Device{}
	r.GET("/list", d.ListHandler)
}
