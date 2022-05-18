package device

import (
	"xstation/entity/cache"
	"xstation/internal/errors"
	"xstation/model"

	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xutils/ctx"

	"github.com/gin-gonic/gin"
)

type Status struct {
}

func (o *Status) ListHandler(c *gin.Context) {
	var p Where
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var (
		data  []model.DevStatus
		total int64
	)
	m := cache.Device(p.DeviceNo)
	if m != nil {
		total, _ = orm.DbByWhere(m.Model(), p.Status()).Find(&data)
	}
	ctx.JSONWriteData(gin.H{"total": total, "data": data}, c)
}

// statusGet 获取
type statusGet struct {
	DeviceNo string `form:"deviceNo"` //
	StatusId uint   `form:"statusId"` //
}

// GetHandler 获取指定id
func (o *Status) GetHandler(c *gin.Context) {
	p := statusGet{}
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	m := cache.Device(p.DeviceNo)
	if m != nil {
		ctx.JSONWriteError(errors.InvalidDeviceNo, c)
		return
	}
	data := m.Model()
	if err := orm.DbFirstById(data, p.StatusId); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONWriteData(data, c)
}

func statusRouter(r *gin.RouterGroup) {
	s := Status{}
	r.GET("/list", s.ListHandler)
	r.GET("", s.GetHandler)
}
