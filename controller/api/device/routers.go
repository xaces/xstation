package device

import "github.com/gin-gonic/gin"

func Routers(r *gin.RouterGroup) {
	deviceRouter(r)
	requestRouter(r.Group("/request"))
	controlRouter(r.Group("/control"))
	statusRouter(r.Group("/status"))
	onlineRouter(r.Group("/online"))
	alarmRouter(r.Group("/alarm"))
}
