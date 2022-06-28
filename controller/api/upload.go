package api

import (
	"github.com/xaces/xutils/ctx"

	"github.com/gin-gonic/gin"
)

// UploadHandler 上传文件
func UploadHandler(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	c.SaveUploadedFile(f, f.Filename)
	ctx.JSONWriteData(gin.H{"fileName": f.Filename, "fileSize": f.Size}, c)
}
