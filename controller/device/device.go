package device

import (
	"fmt"
	"log"
	"xstation/configs"
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"
	"xstation/util"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

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

func AccessHandler(b []byte, a *xproto.Access) (interface{}, error) {
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

func DroppedHandler(v interface{}, a *xproto.Access, err error) {
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
func StatusHandler(tag string, s *xproto.Status) {
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
		Handler.AddStatus(*o)
	}
	// 第三方hook
}
func AlarmHandler(b []byte, a *xproto.Alarm) {
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
		Handler.AddStatus(*status)
		Handler.AddAlarm(*o)
	}
	// 第三方hook
}
func EventHandler(data []byte, e *xproto.Event) {
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

// ftpLittleFile 新文件通知
func ftpLittleFile(e *xproto.Event) {
	v := e.Data.(xproto.EventFileLittle)
	fpName := util.FilePath(v.FileName, e.DeviceNo)
	ftp := xproto.FtpTransfer{
		FtpURL:   configs.FtpAddr,
		FileSrc:  v.FileName,
		FileDst:  configs.Default.Public + "/" + fpName,
		Action:   xproto.ACTION_Download,
		FileType: v.FileType,
		Session:  e.Session,
	}
	if xproto.SyncSend(xproto.REQ_FtpTransfer, ftp, nil, e.DeviceNo) != nil {
		return
	}
	alr := mnger.Alarm.Get(e.Session) // 从缓存中获取数据
	if alr == nil {
		return
	}
	data := &model.DevAlarmFile{}
	data.Guid = e.Session
	data.LinkType = model.AlarmLinkFtpFile
	data.DeviceNo = e.DeviceNo
	data.AlarmType = alr.AlarmType
	data.DTU = e.DTU
	data.Channel = v.Channel
	data.Size = v.Size
	data.Duration = v.Duration
	data.FileType = v.FileType
	data.Name = ftp.FileDst
	orm.DbCreate(data)
}

//
func ftpTimedCapture(e *xproto.Event) {
	v := e.Data.(xproto.EventFileTimedCapture)
	fpName := util.FilePicPath(v.FileName, e.DeviceNo)
	ftp := xproto.FtpTransfer{
		FtpURL:   configs.FtpAddr,
		FileSrc:  v.FileName,
		FileDst:  configs.Default.Public + "/" + fpName,
		Action:   xproto.ACTION_Download,
		FileType: xproto.FILE_NormalPic,
		Session:  e.Session,
	}
	if xproto.SyncSend(xproto.REQ_FtpTransfer, ftp, nil, e.DeviceNo) != nil {
		return
	}
	data := &model.DevCapture{
		DeviceNo:  e.DeviceNo,
		DTU:       e.DTU,
		Channel:   v.Channel,
		Latitude:  v.Latitude,
		Longitude: v.Longitude,
		Speed:     v.Speed,
		Name:      ftp.FileDst,
	}
	orm.DbCreate(data)
}
