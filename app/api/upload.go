package api

import (
	"xstation/pkg/ctx"

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
	ctx.JSONOk().WriteData(gin.H{"fileName": f.Filename, "fileSize": f.Size}, c)
}
