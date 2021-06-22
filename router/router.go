package router

import (
	"fmt"
	"net/http"
	"time"
	"xstation/app/api"
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

// deiveRouter 设备路由
func deiveRouter(r *gin.RouterGroup) {
	p := r.Group("/device/request")
	p.POST("/liveStream", device.LiveStreamHandler)
	p.POST("/voice", device.VoiceHandler)
	p.POST("/playback", device.PlaybackHandler)
	p.POST("/query", device.QueryHandler)
	p.POST("/parameters", device.ParametersHandler)
	p.POST("/fileTransfer", device.FileTransferHandler)
	p.POST("/ftpTransfer", device.FtpTransferHandler)
	p.POST("/close", device.CloseLinkHandler)

	ctrl := p.Group("/control")
	ctrl.POST("/ptz", device.ControlPTZHandler)
	ctrl.POST("/reboot", device.ControlRebootHandler)
	ctrl.POST("/capture", device.ControlCaptureHandler)
	ctrl.POST("/osd", device.ControlOsdHandler)

	s := device.Device{}
	r.GET("/device/list", s.ListHandler)
	r.GET("/device/get/:id", s.GetHandler)
	r.POST("/device", s.AddHandler)
	r.PUT("/device", s.UpdateHandler)
	r.DELETE("/device/:id", s.DeleteHandler)

	st := device.Status{}
	r.GET("/device/status/list", st.ListHandler)
	r.GET("/device/status", st.GetHandler)
}

func newApp() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.Use(gin.Logger()) // 日志
	// /StandardLoginAction_terminalLogin.action?update=gStream&live=1&server=login
	r.POST("/StandardLoginAction_terminalLogin.action", device.LoginHandler)
	root := r.Group("/xstation")
	root.POST("/applyAuth", server.ApplyAuthHandler)
	v1 := root.Group("/api")
	v1.POST("/upload", api.UploadHandler)
	serverRouter(v1)
	deiveRouter(v1)
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
