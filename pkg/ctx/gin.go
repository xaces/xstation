package ctx

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Status 响应状态
type Status int

var (
	StatusFail         Status = -1  // 失败
	StatusOK           Status = 200 // 成功
	StatusError        Status = 500 // 错误
	StatusLoginExpired Status = 401 // 登录过期
	StatusForbidden    Status = 403 // 无权限
)

// Context 响应
type Context struct {
	msg  string
	code Status
}

// JSON 响应
func JSON(status Status) *Context {
	ctx := &Context{}
	ctx.code = status
	switch status {
	case StatusOK:
		ctx.msg = "success"
	case StatusFail:
		ctx.msg = "failed"
	case StatusForbidden:
		ctx.msg = "forbidden"
	}
	return ctx
}

// SetMsg 设置消息体的内容
func (resp *Context) SetMsg(msg string) *Context {
	resp.msg = msg
	return resp
}

// SetCode 设置消息体的编码
func (resp *Context) SetCode(code Status) *Context {
	resp.code = code
	return resp
}

// WriteTo 输出json到客户端
func (resp *Context) WriteTo(c *gin.Context) {
	resp.WriteData(gin.H{}, c)
}

// WriteData 输出json到客户端
func (resp *Context) WriteData(h gin.H, c *gin.Context) {
	h["code"] = resp.code
	h["msg"] = resp.msg
	c.JSON(http.StatusOK, h)
}

// JSONWriteError 错误应答
func JSONWriteError(err error, c *gin.Context) {
	JSONError().SetMsg(err.Error()).WriteTo(c)
}

// JSONError 错误
func JSONError() *Context {
	return JSON(StatusError)
}

// JSONOk 正确
func JSONOk() *Context {
	return JSON(StatusOK)
}

// ParamInt int参数
func ParamInt(c *gin.Context, key string) (uint64, error) {
	idstr := c.Param(key)
	id, err := strconv.Atoi(idstr)
	return uint64(id), err
}

// ParamString string
func ParamString(c *gin.Context, key string) string {
	return c.Param(key)
}

// QueryInt int参数
func QueryInt(c *gin.Context, key string) (uint64, error) {
	idstr := c.Query(key)
	id, err := strconv.Atoi(idstr)
	return uint64(id), err
}
