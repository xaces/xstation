package device

import (
	"xstation/app/api/page"
	"xstation/app/serve"
	"xstation/internal"
	"xstation/models"
	"xstation/pkg/ctx"
	"xstation/pkg/orm"

	"github.com/gin-gonic/gin"
)

type Device struct {
}

func (o *Device) ListHandler(c *gin.Context) {
	var param page.Page
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var rows []models.XDevice
	totalCount, err := orm.DbPage(&models.XDevice{}, param.Where()).Find(param.PageNum, param.PageSize, &rows)
	if err == nil {
		ctx.JSONOk().WriteData(gin.H{
			"total": totalCount,
			"rows":  rows}, c)
		return
	}
	ctx.JSONWriteError(err, c)
}

// GetHandler 获取指定id
func (o *Device) GetHandler(c *gin.Context) {
	ID, err := ctx.ParamInt(c, "id")
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data models.XDevice
	err = orm.DbFirstById(ID, &data)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(gin.H{"data": data}, c)
}

// AddHandler 新增
func (o *Device) AddHandler(c *gin.Context) {
	var data models.XDevice
	//获取参数
	if err := c.ShouldBind(&data.XDeviceOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	data.Guid = internal.UUID()
	if err := orm.DbCreate(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	serve.DefaultDevsManager.Add(&data)
	ctx.JSONOk().WriteTo(c)
}

// UpdateHandler 修改
func (o *Device) UpdateHandler(c *gin.Context) {
	var data models.XDevice
	//获取参数
	if err := c.ShouldBind(&data.XDeviceOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := orm.DbUpdateModel(&data); err != nil {
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
	var devs []models.XDevice
	if _, err := orm.DbFindBy(&devs, "id in (?)", ids); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := orm.DbDeletes(models.XDevice{}, ids); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	for _, dev := range devs {
		serve.DefaultDevsManager.Delete(dev.VehiNo)
	}
	ctx.JSONOk().WriteTo(c)
}
