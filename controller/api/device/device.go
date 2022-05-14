package device

import (
	"fmt"
	"strconv"
	"xstation/configs"
	"xstation/entity/cache"
	"xstation/internal/errors"
	"xstation/model"
	"xstation/util"

	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"
	"gorm.io/gorm"

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

// AddHandler 新增
func (o *Device) AddHandler(c *gin.Context) {
	if configs.MsgProc > 0 {
		ctx.JSONError().WriteTo(c)
		return
	}
	var p model.Device
	//获取参数
	if err := c.ShouldBind(&p.DeviceOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := orm.DbCreate(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	cache.NewDevice(cache.DeviceInfo{ID: p.ID, No: p.No, EffectiveTime: p.EffectiveTime})
	ctx.JSONOk().WriteTo(c)
}

type batchAdd struct {
	Prefix      string `json:"prefix"`
	StartNumber int    `json:"startNumber"`
	Count       int    `json:"count"`
	model.DeviceOpt
}

// BatchAddHandler 新增
func (o *Device) BatchAddHandler(c *gin.Context) {
	if configs.MsgProc > 0 {
		ctx.JSONError().WriteTo(c)
		return
	}
	var p batchAdd
	//获取参数
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	lzero := len(strconv.Itoa(p.StartNumber + p.Count - 1))
	var data []model.Device
	for i := 0; i < p.Count; i++ {
		v := model.Device{}
		v.DeviceOpt = p.DeviceOpt
		v.No = fmt.Sprintf("%s%0*d", p.Prefix, lzero, p.StartNumber+i)
		v.Name = v.No
		data = append(data, v)
	}
	if err := orm.DbCreate(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	for _, v := range data {
		cache.NewDevice(cache.DeviceInfo{ID: v.ID, No: v.No, EffectiveTime: v.EffectiveTime})
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
	var data []model.Device
	ids := util.StringToIntSlice(idstr, ",")
	if err := orm.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id in (?)", ids).Find(&data).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Device{}).Delete(ids).Error; err != nil {
			return err
		}
		for _, v := range data {
			cache.DeviceDel(v.No)
		}
		return nil
	}); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

func deviceRouter(r *gin.RouterGroup) {
	o := Device{}
	r.GET("/list", o.ListHandler)
	r.GET("/:deviceNo", o.GetHandler)
	r.POST("", o.AddHandler)
	r.POST("/batchAdd", o.BatchAddHandler)
	r.DELETE("/:id", o.DeleteHandler)
}
