package device

import (
	"xstation/internal/errors"

	"github.com/wlgd/xutils/ctx"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/wlgd/xproto"
)

type devInfo struct {
	deviceNo string
}

func checkParam(c *gin.Context, param interface{}) (*devInfo, error) {
	dvr := &devInfo{}
	dvr.deviceNo = c.Query("deviceNo")
	if dvr.deviceNo == "" {
		return nil, errors.InvalidDeviceNo
	}
	if param == nil {
		return dvr, nil
	}
	return dvr, c.ShouldBind(&param)
}

// DvrQuery 查询
func QueryHandler(c *gin.Context) {
	var param xproto.Query
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp []xproto.File
	if err := xproto.SyncSend(xproto.Req_Query, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrParameters 参数管理 获取/设置
func ParametersHandler(c *gin.Context) {
	var param interface{}
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Parameters, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrControl 设置控制
func ControlPTZHandler(c *gin.Context) {
	var param xproto.Control
	i, err := checkParam(c, &param.Data)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param.Type = xproto.Ctrl_PTZ
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Control, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// ControlRebootHandler 设置控制
func ControlRebootHandler(c *gin.Context) {
	i, err := checkParam(c, nil)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param := xproto.Control{Type: xproto.Ctrl_Reboot}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Control, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// ControlCaptureHandler 设置控制
func ControlCaptureHandler(c *gin.Context) {
	var param xproto.Control
	i, err := checkParam(c, &param.Data)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param.Type = xproto.Ctrl_Capture
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Control, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// ControlOsdHandler 设置控制
func ControlOsdHandler(c *gin.Context) {
	var param xproto.Control
	i, err := checkParam(c, &param.Data)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param.Type = xproto.Ctrl_OsdSpeed
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Control, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// ControlGsensorHandler 设置控制
func ControlGsensorHandler(c *gin.Context) {
	var param xproto.Control
	i, err := checkParam(c, &param.Data)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param.Type = xproto.Ctrl_Vehi
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Control, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// ControlVehicleHandler 设置控制
func ControlVehicleHandler(c *gin.Context) {
	i, err := checkParam(c, nil)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param := &xproto.Control{Type: xproto.Ctrl_Vehi}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Control, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// ControlResetHandler 设置控制
func ControlResetHandler(c *gin.Context) {
	i, err := checkParam(c, nil)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	param := &xproto.Control{Type: xproto.Ctrl_Reset}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Control, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrLiveStream 实时视频请求
// return url 直接打开播放
func LiveStreamHandler(c *gin.Context) {
	var param xproto.LiveStream
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_LiveStream, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// VoiceHandler 语音业务
// return url 直接打开播放
func VoiceHandler(c *gin.Context) {
	var param xproto.Voice
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Voice, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrPlayback 录像回放
func PlaybackHandler(c *gin.Context) {
	var param xproto.Playback
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_Playback, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrSerialTransparent 串口透传设置
func SerialTransparentHandler(c *gin.Context) {
	var param xproto.SerialTransparent
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_SerialTransparent, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrSerialTransfer 串口透传数据
func SerialTransferHandler(c *gin.Context) {
	session := c.Query("session")
	if session == "" {
		ctx.JSONWriteError(xproto.ErrParam, c)
		return
	}
	var param xproto.RawData
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_SerialTransfer, param, &resp, i.deviceNo, session); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrFileTransfer 文件传输
func FileTransferHandler(c *gin.Context) {
	var param xproto.FileTransfer
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	if param.Action == xproto.Act_Upload {
		f, err := os.Stat(param.FileName)
		if err != nil {
			ctx.JSONWriteError(err, c)
			return
		}
		param.FileSize = int(f.Size())
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_FileTransfer, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrFtpTransfer ftp文件传输
func FtpTransferHandler(c *gin.Context) {
	var param xproto.FtpTransfer
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_FtpTransfer, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// UserDefineHandler
func UserDefineHandler(c *gin.Context) {
	var param xproto.UserDefine
	i, err := checkParam(c, &param)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	var resp interface{}
	if err := xproto.SyncSend(xproto.Req_UserDefine, param, &resp, i.deviceNo); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteData(resp, c)
}

// DvrCloseLink 关闭链路
func CloseLinkHandler(c *gin.Context) {
	i, err := checkParam(c, nil)
	if err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	sessionID := c.Query("session")
	if err := xproto.SyncSend(xproto.Req_Close, nil, nil, i.deviceNo, sessionID); err != nil {
		ctx.JSONWriteError(err, c)
		return
	}
	ctx.JSONOk().WriteTo(c)
}

func requestRouter(r *gin.RouterGroup) {
	r.POST("/liveStream", LiveStreamHandler)
	r.POST("/voice", VoiceHandler)
	r.POST("/playback", PlaybackHandler)
	r.POST("/query", QueryHandler)
	r.POST("/parameters", ParametersHandler)
	r.POST("/fileTransfer", FileTransferHandler)
	r.POST("/ftpTransfer", FtpTransferHandler)
	r.POST("/user", UserDefineHandler)
	r.POST("/close", CloseLinkHandler)

}

func controlRouter(r *gin.RouterGroup) {
	r.POST("/ptz", ControlPTZHandler)
	r.POST("/reboot", ControlRebootHandler)
	r.POST("/capture", ControlCaptureHandler)
	r.POST("/osd", ControlOsdHandler)
	r.POST("/reset", ControlResetHandler)
	r.POST("/vehicle", ControlVehicleHandler)
	r.POST("/gsensor", ControlGsensorHandler)
}
