package service

import (
	"github.com/gin-gonic/gin"
	"github.com/wlgd/xutils/ctx"
	"github.com/wlgd/xutils/orm"
)

func QueryById(v interface{}, c *gin.Context) {
	id := ctx.ParamUInt(c, "id")
	if err := orm.DbFirstById(v, id); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(v, c)
}
