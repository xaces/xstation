package device

import "github.com/gin-gonic/gin"

func InitRouters(r *gin.RouterGroup) {
	api := r.Group("/device")
	DeviceRouter(api)
	RequestRouter(api)
	ControlRouter(api)
	StatusRouter(api)
	OnlineRouter(api)
	AlarmRouter(api)
}
