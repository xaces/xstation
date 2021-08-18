package controller

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"xstation/controller/api"
	"xstation/middleware"

	"xstation/controller/api/device"
	"xstation/controller/api/sys"

	"github.com/gin-gonic/gin"
)

// serverRouter 服务路由
func serverRouter(r *gin.RouterGroup) {
	s := sys.Serve{}
	r.GET("/serve/list", s.ListHandler)
	r.GET("/serve/get/:guid", s.GetHandler)
	r.POST("/serve", s.AddHandler)
	r.PUT("/serve", s.UpdateHandler)
	r.PUT("/serve/status", s.UpdateStatusHandler)    // 设置子服务
	r.GET("/serve/status/list", s.StatusListHandler) // 获取子服务状态
	r.DELETE("/serve", s.DeleteHandler)
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

	ctrl := r.Group("/device/control")
	ctrl.POST("/ptz", device.ControlPTZHandler)
	ctrl.POST("/reboot", device.ControlRebootHandler)
	ctrl.POST("/capture", device.ControlCaptureHandler)
	ctrl.POST("/osd", device.ControlOsdHandler)
	ctrl.POST("/reset", device.ControlResetHandler)
	ctrl.POST("/vehicle", device.ControlVehicleHandler)
	ctrl.POST("/gsensor", device.ControlGsensorHandler)

	s := device.Device{}
	r.GET("/device/list", s.ListHandler)
	r.GET("/device/get/:id", s.GetHandler)
	r.POST("/device", s.AddHandler)
	r.PUT("/device", s.UpdateHandler)
	r.DELETE("/device/:id", s.DeleteHandler)

	st := device.Status{}
	r.GET("/device/status/list", st.ListHandler)
	r.GET("/device/status", st.GetHandler)

	on := device.OnLine{}
	r.POST("/device/online", on.AddHandler)
	r.GET("/device/online/list", on.ListHandler)

	alr := device.Alarm{}
	r.GET("/device/alarm/list", alr.ListHandler)
	r.GET("/device/alarm/bystatus/:statusId", alr.GetByStatusIdHandler)
}

func newRouters() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Logger())      // 日志
	r.Use(middleware.Cors()) // 跨域
	root := r.Group("/station")
	v1 := root.Group("/api")
	v1.POST("/upload", api.UploadHandler)
	serverRouter(v1)
	deiveRouter(v1)
	return r
}

// NewServer
func NewServer(address string) *http.Server {
	as := strings.Split(address, ":")
	addr := fmt.Sprintf(":%s", as[1])
	r := newRouters()
	return &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
