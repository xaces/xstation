package device

import (
	"fmt"
	"log"
	"xstation/configs"
	"xstation/entity/cache"
	"xstation/model"
	"xstation/util"

	"github.com/xaces/xproto"
	"github.com/xaces/xutils/orm"
)

func AccessHandler(b []byte, a *xproto.Access) (interface{}, error) {
	switch a.LinkType {
	case xproto.Link_Signal:
		m := cache.GetDevice(a.DeviceNo)
		log.Printf("%s", b)
		if m == nil {
			return nil, fmt.Errorf("[%s] invalid", a.DeviceNo)
		}
		a.DeviceID = m.ID
		m.Update(a)
		for _, v := range hooks {
			v.PublishOnline(a)
		}
		if configs.MsgProc == 0 {
			return m, deviceUpdate(a)
		}
		return m, devOnlineUpdate(a, &m.Status)
	case xproto.Link_FileTransfer:
		filename, act := xproto.FileOfSess(a.Session)
		if act == xproto.Act_Upload {
			return nil, xproto.UploadFile(a, filename, true)
		}
		return xproto.DownloadFile(configs.Default.Public+"/"+filename, nil)
	default:
		return nil, xproto.ErrUnSupport
	}
}

func DroppedHandler(v interface{}, a *xproto.Access, err error) {
	switch a.LinkType {
	case xproto.Link_FileTransfer:
		_, act := xproto.FileOfSess(a.Session)
		if act == xproto.Act_Upload {
			xproto.DownloadFile("", v)
		}
	case xproto.Link_Signal:
		log.Println(err)
		m, ok := v.(*cache.Device)
		if !ok || m == nil {
			break
		}
		m.Update(a)
		for _, v := range hooks {
			v.PublishOnline(a)
		}
		if configs.MsgProc == 0 {
			deviceUpdate(a)
			return
		}
		devOnlineUpdate(a, &m.Status)
	}
}

func StatusHandler(tag string, arg interface{}, s *xproto.Status) {
	// xproto.LogStatus(tag, s)
	m, ok := arg.(*cache.Device)
	if !ok || m == nil {
		return
	}
	if s.Flag == 0 {
		m.LastOnlineTime = s.DTU
		m.Status = *s
	}
	for _, v := range hooks {
		v.PublishStatus(s)
	}
	o := devStatusModel(s)
	o.DeviceID = m.ID
	Handler.status <- o
}

func AlarmHandler(b []byte, arg interface{}, a *xproto.Alarm) {
	xproto.LogAlarm(b, arg, a)
	for _, v := range hooks {
		v.PublishAlarm(a)
	}
	o := devAlarmDetailsModel(a)
	if o.Status == 0 {
		cache.NewDevAlarm(o)
	}
	Handler.status <- o.DevStatus
	Handler.alarm <- o
}

func EventHandler(data []byte, arg interface{}, e *xproto.Event) {
	xproto.LogEvent(data, arg, e)
	if configs.MsgProc == 0 {
		return
	}
	for _, v := range hooks {
		v.PublishEvent(e)
	}
	switch e.Type {
	case xproto.Event_FtpTransfer:
	case xproto.Event_FileLittle:
		ftpLittleFile(e)
	case xproto.Event_FileTimedCapture:
		ftpTimedCapture(e)
	}
}

// ftpLittleFile 新文件通知
func ftpLittleFile(e *xproto.Event) {
	alr := cache.DevAlarm(e.Session) // 从缓存中获取数据
	if alr == nil {
		return
	}
	v := e.Data.(xproto.FileLittle)
	fpName := util.FilePath(v.FileName, e.DeviceNo)
	ftp := xproto.FtpTransfer{
		FtpURL:   configs.FtpAddr,
		FileSrc:  v.FileName,
		FileDst:  fpName,
		Action:   xproto.Act_Download,
		FileType: v.FileType,
		Session:  e.Session,
	}
	if xproto.SyncSend(xproto.Req_FtpTransfer, ftp, nil, e.DeviceNo) != nil {
		return
	}
	data := &model.DevAlarmFile{
		GUID:      e.Session,
		LinkType:  model.AlarmLinkFtpFile,
		DeviceID:  e.DeviceID,
		AlarmType: alr.AlarmType,
		DTU:       e.DTU,
		Channel:   v.Channel,
		Size:      v.Size,
		Duration:  v.Duration,
		FileType:  v.FileType,
		Name:      ftp.FileDst,
	}
	orm.DbCreate(data)
}

//
func ftpTimedCapture(e *xproto.Event) {
	v := e.Data.(xproto.FileCapture)
	fpName := util.FilePicPath(v.FileName, e.DeviceNo)
	ftp := xproto.FtpTransfer{
		FtpURL:   configs.FtpAddr,
		FileSrc:  v.FileName,
		FileDst:  fpName,
		Action:   xproto.Act_Download,
		FileType: xproto.File_NormalPic,
		Session:  e.Session,
	}
	if xproto.SyncSend(xproto.Req_FtpTransfer, ftp, nil, e.DeviceNo) != nil {
		return
	}
	data := &model.DevCapture{
		DeviceID:  e.DeviceID,
		DTU:       e.DTU,
		Channel:   v.Channel,
		Latitude:  v.Latitude,
		Longitude: v.Longitude,
		Speed:     v.Speed,
		Name:      ftp.FileDst,
	}
	orm.DbCreate(data)
}
