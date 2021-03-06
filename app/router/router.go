package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"xstation/configs"
	"xstation/controller/api"

	"xstation/controller/api/device"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func newApp() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Logger()) // 日志
	// r.Use(middleware.Cors()) // 跨域
	root := r.Group("/station")
	v1 := root.Group("/api")
	v1.POST("/upload", api.UploadHandler)
	v1.StaticFS("/public", http.Dir(configs.Default.Public))
	device.Routers(v1.Group("/device"))
	return r
}

var (
	s *http.Server
)

// Run
func Run(port uint16) {
	address := fmt.Sprintf(":%d", port)
	r := newApp()
	s = &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	pprof.Register(r)
	go s.ListenAndServe()
	log.Printf("http server listening at: %s\n", address)
}

func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
}
