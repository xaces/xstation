package sys

import (
	"xstation/internal"
	"xstation/mnger"
	"xstation/model"

	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xutils/ctx"

	"github.com/gin-gonic/gin"
)

type Serve struct {
}

// ListHandler
func (o *Serve) ListHandler(c *gin.Context) {
	var serves []model.SysServe
	if err := orm.DbFind(&serves); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(serves, c)
}

// ListHandler
func (o *Serve) GetHandler(c *gin.Context) {
	guid := ctx.ParamString(c, "guid")
	var server model.SysServe
	if err := orm.DbFirstBy(&server, "guid like ?", guid); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var st model.SysServe
	orm.DbFirstBy(&st, "role = ?", 1)
	ctx.JSONOk().WriteData(gin.H{"station": st.SysServeOpt, "local": server.SysServeOpt}, c)
}

// AddHandler 新增
func (o *Serve) AddHandler(c *gin.Context) {
	var param model.SysServe
	if err := c.ShouldBindJSON(&param.SysServeOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param.Guid = internal.ServeId()
	if err := orm.DbCreate(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	mnger.Serve.Add(&param)
	// 同步信息到服务管理
	ctx.JSONOk().WriteTo(c)
}

// UpdateHandler
func (o *Serve) UpdateHandler(c *gin.Context) {
	var param model.SysServe
	if err := c.ShouldBindJSON(&param.SysServeOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := orm.DbUpdateModelBy(&param, "guid like ?", param.Guid); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

// UpdateStatusHandler
func (o *Serve) UpdateStatusHandler(c *gin.Context) {
	var param serveOpt
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := mnger.Serve.UpdateStatus(param.Guids, param.Status); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

// StatusListHandler
func (o *Serve) StatusListHandler(c *gin.Context) {
	data := mnger.Serve.GetAll()
	ctx.JSONOk().WriteData(data, c)
}

// DeleteHandler 删除
func (o *Serve) DeleteHandler(c *gin.Context) {
	var param serveOpt
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := deleteServes(param.Guids); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}
