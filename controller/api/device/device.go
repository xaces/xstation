package device

import (
	"xstation/entity/cache"
	"xstation/internal/errors"

	"github.com/wlgd/xutils/ctx"

	"github.com/gin-gonic/gin"
)

type Device struct {
}

func (o *Device) ListHandler(c *gin.Context) {
	data := cache.DeviceList()
	ctx.JSONOk().Write(gin.H{"total": len(data), "data": data}, c)
}

// GetHandler 获取指定id
func (o *Device) GetHandler(c *gin.Context) {
	deviceNo := ctx.ParamString(c, "deviceNo")
	if data := cache.Device(deviceNo); data != nil {
		ctx.JSONOk().WriteData(data, c)
		return
	}
	ctx.JSONWriteError(errors.InvalidDeviceNo, c)
}

func DeviceRouter(r *gin.RouterGroup) {
	d := Device{}
	r.GET("/list", d.ListHandler)
	r.GET("/:deviceNo", d.GetHandler)
}
