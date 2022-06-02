package device

import (
	"fmt"
	"log"
	"xstation/configs"
	"xstation/entity/cache"
	"xstation/model"
	"xstation/util"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

func AccessHandler(b []byte, a *xproto.Access) (interface{}, error) {
	switch a.LinkType {
	case xproto.Link_Signal:
		m := cache.Device(a.DeviceNo)
		log.Printf("%s", b)
		if m == nil {
			return nil, fmt.Errorf("[%s] invalid", a.DeviceNo)
		}
		m.Update(a)
		for _, v := range hooks {
			v.Online(m.ID, a)
		}
		if configs.MsgProc == 0 {
			return nil, deviceUpdate(m.ID, a)
		}
		return nil, devOnlineUpdate(a, &m.Status)
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
		m := cache.Device(a.DeviceNo)
		if m == nil {
			return
		}
		m.Update(a)
		for _, v := range hooks {
			v.Online(m.ID, a)
		}
		if configs.MsgProc == 0 {
			deviceUpdate(m.ID, a)
			return
		}
		devOnlineUpdate(a, &m.Status)
	}
}

func StatusHandler(tag string, s *xproto.Status) {
	xproto.LogStatus(tag, s)
	m := cache.Device(s.DeviceNo)
	if m == nil {
		return
	}
	if s.Flag == 0 {
		m.LastOnlineTime = s.DTU
		m.Status = *s
	}
	for _, v := range hooks {
		v.Status(m.ID, s)
	}
	o := devStatusModel(s)
	o.DeviceID = m.ID
	Handler.status <- o
}

func AlarmHandler(b []byte, a *xproto.Alarm) {
	xproto.LogAlarm(b, a)
	m := cache.Device(a.DeviceNo)
	if m == nil {
		return
	}
	for _, v := range hooks {
		v.Alarm(m.ID, a)
	}
	o := devAlarmDetailsModel(a)
	if o.Status == 0 {
		cache.NewDevAlarm(o)
	}
	Handler.alarm <- o
}

func EventHandler(data []byte, e *xproto.Event) {
	xproto.LogEvent(data, e)
	if configs.MsgProc == 0 {
		return
	}
	m := cache.Device(e.DeviceNo)
	if m == nil {
		return
	}
	for _, v := range hooks {
		v.Event(m.ID, e)
	}
	switch e.Type {
	case xproto.Event_FtpTransfer:
	case xproto.Event_FileLittle:
		alr := cache.DevAlarm(e.Session) // 从缓存中获取数据
		if alr == nil {
			return
		}
		ftpLittleFile(e, alr.AlarmType)
	case xproto.Event_FileTimedCapture:
		ftpTimedCapture(e)
	}
}

// ftpLittleFile 新文件通知
func ftpLittleFile(e *xproto.Event, alrType int) {
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
		DeviceNo:  e.DeviceNo,
		AlarmType: alrType,
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
