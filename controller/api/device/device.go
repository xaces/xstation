package device

import (
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"
	"xstation/util"

	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"

	"github.com/gin-gonic/gin"
)

type Device struct {
}

func (o *Device) ListHandler(c *gin.Context) {
	var param service.DevicePage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data []model.Device
	total, _ := orm.DbByWhere(&data, param.Where()).Find(&data)
	ctx.JSONOk().Write(gin.H{"total": total, "data": data}, c)
}

// GetHandler 获取指定id
func (o *Device) GetHandler(c *gin.Context) {
	service.BasicGet(&model.Device{}, c)
}

// AddHandler 新增
func (o *Device) AddHandler(c *gin.Context) {
	var data model.Device
	//获取参数
	if err := c.ShouldBind(&data.DeviceOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	data.Guid = util.UUID()
	if err := orm.DbCreate(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	mnger.Device.Add(&data)
	ctx.JSONOk().WriteTo(c)
}

// UpdateHandler 修改
func (o *Device) UpdateHandler(c *gin.Context) {
	var data model.Device
	//获取参数
	if err := c.ShouldBind(&data.DeviceOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := orm.DbUpdateModel(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

// resetOrganize 更新
type resetOrganize struct {
	Devices    string `json:"devices"`
	OrganizeId int    `json:"organizeId"`
}

func (o *Device) ResetOrganizeHandler(c *gin.Context) {
	var param resetOrganize
	//获取参数
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ids := util.StringToIntSlice(param.Devices, ",")
	if err := orm.DbUpdateByIds(&model.Device{}, ids, map[string]interface{}{"organize_id": param.OrganizeId}); err != nil {
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
	ids := util.StringToIntSlice(idstr, ",")
	var devs []model.Device
	if _, err := orm.DbFindBy(&devs, "id in (?)", ids); err != nil {
		ctx.JSONError().WriteTo(c)
		return
	}
	if err := orm.DbDeletes(&devs); err != nil {
		ctx.JSONError().WriteTo(c)
		return
	}
	for _, v := range devs {
		mnger.Device.Delete(v.DeviceNo)
	}
	ctx.JSONOk().WriteTo(c)
}

func DeviceRouter(r *gin.RouterGroup) {
	d := Device{}
	r.GET("/list", d.ListHandler)
	r.GET("/:id", d.GetHandler)
	r.POST("", d.AddHandler)
	r.PUT("", d.UpdateHandler)
	r.PUT("/resetOrganize", d.ResetOrganizeHandler)
	r.DELETE("/:id", d.DeleteHandler)
}
