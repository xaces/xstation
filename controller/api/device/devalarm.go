package device

import (
	"xstation/model"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"
)

type Alarm struct {
}

func (o *Alarm) ListHandler(c *gin.Context) {
	var p Where
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data []model.DevAlarm
	total, _ := orm.DbByWhere(&model.DevAlarm{}, p.Alarm()).Find(&data)
	ctx.JSONWriteData(gin.H{"total": total, "data": data}, c)
}

func (o *Alarm) ListDetailsHandler(c *gin.Context) {
	var p Where
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data []model.DevAlarmDetails
	total, _ := orm.DbByWhere(&model.DevAlarmDetails{}, p.AlarmDetailsPage()).Find(&data)
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
