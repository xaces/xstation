package device

import (
	"xstation/app/mnger"
	"xstation/model"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"
)

type Alarm struct {
}

type alarm struct {
	Alarm  model.DevAlarm `json:"alarm"`
	Status interface{} `json:"status"`
}

func (o *Alarm) ListHandler(c *gin.Context) {
	var param alarmPage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var rows []model.DevAlarm
	totalCount, err := orm.DbPage(&model.DevAlarm{}, param.Where()).Scan(param.PageNum, param.PageSize, &rows)
	if err == nil {
		var data []alarm
		for _, v := range rows {
			alr := alarm{Alarm: v}
			_, m := mnger.Dev.GetModel(param.DeviceNo)
			if err := orm.DbFirstById(m, v.StatusId); err == nil {
				alr.Status = m
			}
			data = append(data, alr)
		}
		ctx.JSONOk().WriteData(gin.H{"total": totalCount, "rows": data}, c)
		return
	}
	ctx.JSONWriteError(err, c)
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
