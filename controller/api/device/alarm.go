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
	var param service.AlarmPage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var data []model.DevAlarm
	total, _ := orm.DbPage(&model.DevOnline{}, param.Where()).Scan(param.PageNum, param.PageSize, &data)
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
	r.GET("/alarm/bystatus/:statusId", alr.GetByStatusIdHandler)
}
