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
	if err := orm.DB().Raw(p.Where()).Scan(&rows).Error; err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	total := len(rows)
	if p.PageSize > 0 {
		rows = rows[p.PageSize*(p.PageNum-1) : p.PageSize*p.PageNum]
	}
	ctx.JSONOk().WriteData(gin.H{"total": total, "rows": rows}, c)
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
