package router

import (
	"fmt"
	"net/http"
	"time"
	"xstation/configs"
	"xstation/controller/api"
	"xstation/middleware"

	"xstation/controller/api/device"
	"xstation/controller/api/sys"

	"github.com/gin-gonic/gin"
)

func newRouters() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Logger())      // 日志
	r.Use(middleware.Cors()) // 跨域
	root := r.Group("/station")
	v1 := root.Group("/api")
	v1.POST("/upload", api.UploadHandler)
	v1.StaticFS("/public", http.Dir(configs.Default.Public))
	sys.ServerRouter(v1)
	dev := v1.Group("/device")
	device.Router(dev)
	device.RequestRouter(dev)
	device.ControlRouter(dev)
	device.StatusRouter(dev)
	device.OnlineRouter(dev)
	device.AlarmRouter(dev)
	return r
}

// NewServer
func NewServer(port uint16) *http.Server {
	address := fmt.Sprintf(":%d", port)
	r := newRouters()
	s := &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go s.ListenAndServe()
	return s
}
