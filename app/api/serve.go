package api

import (
	"xstation/app/mnger"
	"xstation/internal"
	"xstation/model"

	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xutils/ctx"

	"github.com/gin-gonic/gin"
)

// ApplyAuthHandler 身份授权
func ApplyAuthHandler(c *gin.Context) {
	var param applyAuth
	if err := c.ShouldBindJSON(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := tryApplyAuth(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

type Serve struct {
}

// ListHandler
func (o *Serve) ListHandler(c *gin.Context) {
	var serves []model.Serve
	if err := orm.DbFind(&serves); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(serves, c)
}

// ListHandler
func (o *Serve) GetHandler(c *gin.Context) {
	guid := ctx.ParamString(c, "guid")
	var server model.Serve
	if err := orm.DbFirstBy(&server, "guid like ?", guid); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var st model.Serve
	orm.DbFirstBy(&st, "role = ?", 1)
	ctx.JSONOk().WriteData(gin.H{"station": st.ServeOpt, "local": server.ServeOpt}, c)
}

// AddHandler 新增
func (o *Serve) AddHandler(c *gin.Context) {
	var param model.Serve
	if err := c.ShouldBindJSON(&param.ServeOpt); err != nil {
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
	var param model.Serve
	if err := c.ShouldBindJSON(&param.ServeOpt); err != nil {
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
