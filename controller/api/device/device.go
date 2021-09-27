package device

import (
	"xstation/internal"
	"xstation/mnger"
	"xstation/model"

	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"

	"github.com/gin-gonic/gin"
)

type Device struct {
}

func (o *Device) ListHandler(c *gin.Context) {
	var param devicePage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var rows []model.Device
	totalCount, err := orm.DbPage(&model.Device{}, param.Where()).Find(param.PageNum, param.PageSize, &rows)
	if err == nil {
		ctx.JSONOk().Write(gin.H{"total": totalCount, "rows": rows}, c)
		return
	}
	ctx.JSONWriteError(err, c)
}

// GetHandler 获取指定id
func (o *Device) GetHandler(c *gin.Context) {
	id, err := ctx.ParamInt(c, "id")
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data model.Device
	err = orm.DbFirstById(&data, id)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(data, c)
}

// AddHandler 新增
func (o *Device) AddHandler(c *gin.Context) {
	var data model.Device
	//获取参数
	if err := c.ShouldBind(&data.DeviceOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	data.Guid = internal.UUID()
	if err := orm.DbCreate(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	mnger.Dev.Add(&data)
	ctx.JSONOk().WriteTo(c)
}

// UpdateHandler 修改
func (o *Device) UpdateHandler(c *gin.Context) {
	var param deviceUpdate
	//获取参数
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	data := &model.Device{
		DeviceOpt: param.DeviceOpt,
	}
	if err := orm.DbUpdateById(&data, param.Id); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

// DeleteHandler 删除
func (o *Device) DeleteHandler(c *gin.Context) {
	idstr := ctx.ParamString(c, "id")
	if idstr == "" {
		ctx.JSONError().WriteTo(c)
		return
	}
	ids := internal.StringToIntSlice(idstr, ",")
	if err := deleteDevices(ids); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}
