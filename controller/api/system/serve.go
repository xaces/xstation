package system

import (
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"
	"xstation/util"

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
	param.Guid = util.ServeId()
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
	var param service.ServeOpt
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
	var param service.ServeOpt
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if _, err := orm.DbDeleteBy(&model.SysServe{}, "guid in (?)", param.Guids); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

// ServerRouter 服务路由
func ServerRouter(r *gin.RouterGroup) {
	s := Serve{}
	r.GET("/serve/list", s.ListHandler)
	r.GET("/serve/get/:guid", s.GetHandler)
	r.POST("/serve", s.AddHandler)
	r.PUT("/serve", s.UpdateHandler)
	r.PUT("/serve/status", s.UpdateStatusHandler)    // 设置子服务
	r.GET("/serve/status/list", s.StatusListHandler) // 获取子服务状态
	r.DELETE("/serve", s.DeleteHandler)
}
