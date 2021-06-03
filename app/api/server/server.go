package server

import (
	"xstation/app/serve"
	"xstation/internal"
	"xstation/models"
	"xstation/pkg/ctx"
	"xstation/pkg/orm"

	"github.com/gin-gonic/gin"
)

type Server struct {
}

// ListHandler
func (o *Server) ListHandler(c *gin.Context) {
	var serves []models.XServer
	if err := orm.DbFind(&serves); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(gin.H{"data": serves}, c)
}

// ListHandler
func (o *Server) GetHandler(c *gin.Context) {
	guid := ctx.ParamString(c, "guid")
	var server models.XServer
	if err := orm.DbFirstBy(&server, "guid like ?", guid); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var st models.XServer
	orm.DbFirstBy(&st, "role = ?", models.ServeTypeLocal)
	ctx.JSONOk().WriteData(gin.H{
		"data": gin.H{
			"station": st.XServerOpt,
			"local":   server.XServerOpt}}, c)
}

// AddHandler 新增
func (o *Server) AddHandler(c *gin.Context) {
	var param models.XServer
	if err := c.ShouldBindJSON(&param.XServerOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param.Guid = internal.ServeId()
	if err := orm.DbCreate(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

// UpdateHandler
func (o *Server) UpdateHandler(c *gin.Context) {
	var param models.XServer
	if err := c.ShouldBindJSON(&param.XServerOpt); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := orm.DbUpdateModelBy(&param, "guid like ?", param.Guid); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

type updateStatus struct {
	Guid   string `json:"guid"`
	Status int    `json:"status"`
}

// UpdateStatusHandler
func (o *Server) UpdateStatusHandler(c *gin.Context) {
	var param updateStatus
	if err := c.ShouldBind(&param); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if err := serve.ChangeStatus(param.Guid, param.Status); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

// StatusListHandler
func (o *Server) StatusListHandler(c *gin.Context) {
	data := serve.LoadAllLServe()
	ctx.JSONOk().WriteData(gin.H{"data": data}, c)
}
