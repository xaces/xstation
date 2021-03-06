package device

import (
	"xstation/internal/errors"
	"xstation/model"

	"github.com/gin-gonic/gin"
	"github.com/xaces/xutils/ctx"
	"github.com/xaces/xutils/orm"
)

type Alarm struct {
}

func (o *Alarm) ListHandler(c *gin.Context) {
	var p Where
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if p.isDeviceNoInvalid() {
		ctx.JSONWriteError(errors.InvalidDeviceNo, c)
		return
	}
	var data []model.DevAlarm
	total, _ := p.Alarm().Model(&model.DevAlarm{}).Find(&data)
	ctx.JSONWriteData(gin.H{"total": total, "data": data}, c)
}

func (o *Alarm) ListDetailsHandler(c *gin.Context) {
	var p Where
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if p.isDeviceNoInvalid() {
		ctx.JSONWriteError(errors.InvalidDeviceNo, c)
		return
	}
	var data []model.DevAlarmDetails
	total, _ := p.AlarmDetails().Model(&model.DevAlarmDetails{}).Find(&data)
	ctx.JSONWriteData(gin.H{"total": total, "data": data}, c)
}

// GetByStatusIdHandler 获取指定id
func (o *Alarm) GetByStatusIdHandler(c *gin.Context) {
	statusId := ctx.ParamUInt(c, "statusId")
	var data model.DevAlarm
	if err := orm.DbFirstBy(&data, "status_id = ?", statusId); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONWriteData(data, c)
}

func alarmRouter(r *gin.RouterGroup) {
	alr := Alarm{}
	r.GET("/list", alr.ListHandler)
	r.GET("/details/list", alr.ListDetailsHandler)
	r.GET("/bystatus/:statusId", alr.GetByStatusIdHandler)
}
