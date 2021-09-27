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
	var p alarmPage
	if err := c.ShouldBind(&p); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var rows []alarmData
	total, err := orm.DbPageRawScan(p.Where(), &rows, p.PageNum, p.PageSize)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().Write(gin.H{"total": total, "rows": rows}, c)
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
