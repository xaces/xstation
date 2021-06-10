package device

import (
	"xstation/app/api/page"
	"xstation/models"
	"xstation/pkg/ctx"
	"xstation/pkg/orm"
	"xstation/service"

	"github.com/gin-gonic/gin"
)

type Status struct {
}

type statusPage struct {
	page.Page
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceId  uint64 `form:"deviceId"` //
	Descs     string `form:"descs"`    //
}

// Where 初始化
func (s *statusPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.Append("device_id = ?", s.DeviceId)
	where.Append("dtu >= ?", s.StartTime)
	where.Append("dtu <= ?", s.EndTime)
	where.Orders = append(where.Orders, s.Descs+" desc")
	return &where
}

func (o *Status) ListHandler(c *gin.Context) {
	var param statusPage
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var rows []models.XStatus
	totalCount, err := orm.DbPage(service.GetXStatusModel(param.DeviceId), param.Where()).Scan(param.PageNum, param.PageSize, &rows)
	if err == nil {
		ctx.JSONOk().WriteData(gin.H{
			"total": totalCount,
			"rows":  rows}, c)
		return
	}
	ctx.JSONWriteError(err, c)
}

type statusGetParam struct {
	DeviceId uint64 `form:"deviceId"` //
	StatusId uint64 `form:"statusId"` //
}

// GetHandler 获取指定id
func (o *Status) GetHandler(c *gin.Context) {
	var param statusGetParam
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	data := service.GetXStatusModel(param.DeviceId)
	if err := orm.DbFirstById(param.StatusId, data); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(gin.H{"data": data}, c)
}
