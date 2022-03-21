package system

import "github.com/gin-gonic/gin"

func InitRouters(r *gin.RouterGroup) {
	api := r.Group("/system")
	SysServerRouter(api)
}