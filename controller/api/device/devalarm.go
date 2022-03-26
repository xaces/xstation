package device

import (
	"xstation/model"
	"xstation/service"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"
)

type Alarm struct {
}

func (o *Alarm) ListHandler(c *gin.Context) {
	var p service.AlarmPage
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data []model.DevAlarm
	total, _ := orm.DbByWhere(&model.DevAlarm{}, p.Where()).Find(&data)
	ctx.JSONOk().Write(gin.H{"total": total, "data": data}, c)
}

func (o *Alarm) ListDetailsHandler(c *gin.Context) {
	var p service.AlarmDetailsPage
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data []model.DevAlarmDetails
	total, _ := orm.DbByWhere(&model.DevAlarmDetails{}, p.Where()).Find(&data)
	ctx.JSONOk().Write(gin.H{"total": total, "data": data}, c)
}

// GetByStatusIdHandler 获取指定id
func (o *Alarm) GetByStatusIdHandler(c *gin.Context) {
	statusId, err := ctx.ParamInt(c, "statusId")
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data model.DevAlarm
	err = orm.DbFirstBy(&data, "status_id = ?", statusId)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(data, c)
}

func AlarmRouter(r *gin.RouterGroup) {
	alr := Alarm{}
	r.GET("/alarm/list", alr.ListHandler)
	r.GET("/alarm/details/list", alr.ListDetailsHandler)
	r.GET("/alarm/bystatus/:statusId", alr.GetByStatusIdHandler)
}
