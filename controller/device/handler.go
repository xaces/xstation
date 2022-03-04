package device

import (
	"fmt"
	"log"
	"xstation/configs"
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"

	"github.com/wlgd/xproto"
)

var (
	msgProc string = ""
)

type Handler struct {
}

func NewHandler(msgproc string) *Handler {
	if msgproc == "nats" {
		natsRun()
	}
	msgProc = msgproc
	devTaskRun()
	return &Handler{}
}

func onlineHandler(a *xproto.Access, online bool) error {
	if a.LinkType != xproto.LINK_Signal {
		return xproto.ErrUnSupport
	}
	m := mnger.Device.Get(a.DeviceNo)
	if m == nil {
		return fmt.Errorf("[%s] invalid", a.DeviceNo)
	}
	o := toDevOnlineModel(a, online)
	switch msgProc {
	case "nats":
		nats.Notify(topicDevOnline, o)
	case "default":
		service.DeviceUpdate(m, online, o.Version, o.DevType)
		service.DevOnlineUpdate(o)
	}
	// 第三方hook
	return nil
}

func (h *Handler) AccessHandler(b []byte, a *xproto.Access) (interface{}, error) {
	log.Printf("%s\n", b)
	if a.LinkType == xproto.LINK_FileTransfer {
		filename, act := xproto.FileOfSess(a.Session)
		if act == xproto.ACTION_Upload {
			return nil, xproto.UploadFile(a, filename, true)
		}
		return xproto.DownloadFile(configs.Default.Public+"/"+filename, nil)
	}
	return nil, onlineHandler(a, true)

}

func (h *Handler) DroppedHandler(v interface{}, a *xproto.Access, err error) {
	log.Println(err)
	if a.LinkType == xproto.LINK_FileTransfer {
		_, act := xproto.FileOfSess(a.Session)
		if act == xproto.ACTION_Upload {
			xproto.DownloadFile("", v)
		}
		return
	}
	onlineHandler(a, false)
}
func (h *Handler) StatusHandler(tag string, s *xproto.Status) {
	// 转换
	xproto.LogStatus(tag, s)
	o := toDevStatusModel(s)
	if o == nil {
		return
	}
	switch msgProc {
	case "nats":
		nats.Notify(topicDevStatus, o)
	case "default":
		devtask.AddStatus(*o)
	}
	// 第三方hook
}
func (h *Handler) AlarmHandler(b []byte, a *xproto.Alarm) {
	xproto.LogAlarm(b, a)
	// 转换
	status := toDevStatusModel(a.Status)
	if status == nil {
		return
	}
	o := toDevAlarmModel(a)
	o.DevStatus = model.JDevStatus(*status)
	switch msgProc {
	case "nats":
		nats.Notify(topicDevAlarm, o)
	case "default":
		devtask.AddStatus(*status)
		devtask.AddAlarm(*o)
	}
	// 第三方hook
}
func (h *Handler) EventHandler(data []byte, e *xproto.Event) {
	xproto.LogEvent(data, e)
	//
	switch msgProc {
	case "nats":
		nats.Notify(topicDevEvent, e)
	case "default":
		devEventHandler(e)
	}
	// 第三方hook
}

func devEventHandler(e *xproto.Event) {
	m := mnger.Device.Get(e.DeviceNo)
	if m == nil || !m.AutoFtp {
		return
	}
	switch e.Type {
	case xproto.EVENT_FtpTransfer:
	case xproto.EVENT_FileLittle:
		ftpLittleFile(e)
	case xproto.EVENT_FileTimedCapture:
		ftpTimedCapture(e)
	}
}
