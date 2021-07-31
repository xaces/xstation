package device

import (
	"xstation/model"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"
)

type OnLine struct {
}

func (o *OnLine) ListHandler(c *gin.Context) {
	var param onlinePage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var rows []model.OnLine
	totalCount, err := orm.DbPage(&model.OnLine{}, param.Where()).Scan(param.PageNum, param.PageSize, &rows)
	if err == nil {
		ctx.JSONOk().WriteData(gin.H{"total": totalCount, "rows": rows}, c)
		return
	}
	ctx.JSONWriteError(err, c)
}

func (o *OnLine) AddHandler(c *gin.Context) {
	var data model.OnLine
	//获取参数
	if err := c.ShouldBind(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := orm.DbCreate(&data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}