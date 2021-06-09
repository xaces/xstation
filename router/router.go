package router

import (
	"fmt"
	"net/http"
	"time"
	"xstation/app/api/dvr"
	"xstation/app/api/server"

	"xstation/app/api/device"

	"github.com/gin-gonic/gin"
)

// serverRouter 服务路由
func serverRouter(r *gin.RouterGroup) {
	s := server.Server{}
	r.GET("/serve/list", s.ListHandler)
	r.GET("/serve/get/:guid", s.GetHandler)
	r.POST("/serve", s.AddHandler)
	r.PUT("/serve", s.UpdateHandler)
	r.PUT("/serve/status", s.UpdateStatusHandler)    // 设置子服务
	r.GET("/serve/status/list", s.StatusListHandler) // 获取子服务状态
}

// xprotoRouter 设备路由
func xprotoRouter(r *gin.RouterGroup) {
	p := r.Group("/dvr")
	p.POST("/liveStream", dvr.LiveStreamHandler)
	p.POST("/playback", dvr.PlaybackHandler)
	p.POST("/query", dvr.QueryHandler)
	p.POST("/parameters", dvr.ParametersHandler)
	p.POST("/control", dvr.ControlHandler)
	p.POST("/fileTransfer", dvr.FileTransferHandler)
	p.POST("/ftpTransfer", dvr.FtpTransferHandler)
	p.POST("/close", dvr.CloseLinkHandler)
}

// deiveRouter 设备路由
func deiveRouter(r *gin.RouterGroup) {
	s := device.Device{}
	r.GET("/device/list", s.ListHandler)
	r.GET("/device/get/:id", s.GetHandler)
	r.POST("/device", s.AddHandler)
	r.PUT("/device", s.UpdateHandler)
	r.DELETE("/device/:id", s.DeleteHandler)
}

func newApp() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Logger()) // 日志
	// /StandardLoginAction_terminalLogin.action?update=gStream&live=1&server=login
	r.POST("/StandardLoginAction_terminalLogin.action", dvr.TerminalLoginHandler)
	root := r.Group("/xstation")
	root.POST("/applyAuth", server.ApplyAuthHandler)
	api := root.Group("/api")
	serverRouter(api)
	xprotoRouter(api)
	deiveRouter(api)
	return r
}

// Init 路由初始化
func Init(port uint16) *http.Server {
	r := newApp()
	address := fmt.Sprintf(":%d", port)
	return &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
