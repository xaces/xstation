package dvr

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

// TerminalLoginHandler 终端登录
// {"userGroupId":3,"result":0,"groupName":"HW","time":"2021-02-20 15:21:44","userName":"353661094105089","userId":30}
func TerminalLoginHandler(c *gin.Context) {
	buf := make([]byte, 1024)
	c.Request.Body.Read(buf)
	bodystr := string(buf)
	arrsystr := strings.Split(bodystr, "&")
	bform := make(map[string]string)
	for _, v := range arrsystr {
		iarray := strings.Split(v, "=")
		if len(iarray) == 2 {
			bform[iarray[0]] = iarray[1]
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"result":      0,
		"userName":    bform["userName"],
		"groupName":   "Default Fleet",
		"userGroupId": 12,
		"userId":      32,
	})
}
