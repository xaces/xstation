package service

import (
	"github.com/gin-gonic/gin"
	"github.com/xaces/xutils/ctx"
	"github.com/xaces/xutils/orm"
)

func QueryByID(v interface{}, c *gin.Context) {
	id := ctx.ParamUInt(c, "id")
	if err := orm.DbFirstById(v, id); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONWriteData(v, c)
}
