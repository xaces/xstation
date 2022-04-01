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

func AccessHandler(b []byte, a *xproto.Access) (interface{}, error) {
	log.Printf("%s\n", b)
	m := mnger.Device.Get(a.DeviceNo)
	if m == nil {
		return nil, fmt.Errorf("[%s] invalid", a.DeviceNo)
	}
	switch a.LinkType {
	case xproto.Link_FileTransfer:
		filename, act := xproto.FileOfSess(a.Session)
		if act == xproto.Act_Upload {
			return nil, xproto.UploadFile(a, filename, true)
		}
		return xproto.DownloadFile(configs.Default.Public+"/"+filename, nil)
	case xproto.Link_Signal:
		for _, v := range hooks {
			v.Online(m.Model.Id, a)
		}
		service.DeviceUpdate(m.Model, a)
		o := devOnlineModel(a)
		o.OnlineTime = a.DeviceTime
		o.OnlineStatus = &m.Status
		return nil, orm.DbCreate(o)
	}
	return nil, xproto.ErrUnSupport

}

func DroppedHandler(v interface{}, a *xproto.Access, err error) {
	log.Println(err)
	m := mnger.Device.Get(a.DeviceNo)
	if m == nil {
		return
	}
	switch a.LinkType {
	case xproto.Link_FileTransfer:
		_, act := xproto.FileOfSess(a.Session)
		if act == xproto.Act_Upload {
			xproto.DownloadFile("", v)
		}
	case xproto.Link_Signal:
		for _, v := range hooks {
			v.Online(m.Model.Id, a)
		}
		o := devOnlineModel(a)
		o.OfflineTime = a.DeviceTime
		o.OfflineStatus = &m.Status
		orm.DbUpdateSelectWhere(o, []string{"offline_time, offline_status"}, "guid = ?", o.Guid)
	}
}

func StatusHandler(tag string, s *xproto.Status) {
	xproto.LogStatus(tag, s)
	m := mnger.Device.Get(s.DeviceNo)
	if m == nil {
		return
	}
	for _, v := range hooks {
		v.Status(m.Model.Id, s)
	}
	o := devStatusModel(s)
	o.DeviceId = m.Model.Id
	if s.Flag == 0 {
		m.Status = model.JDevStatus(*o) // TODO更新设备起始位置
	}
	m.Model.LastOnlineTime = s.DTU
	Handler.status <- o
}

func AlarmHandler(b []byte, a *xproto.Alarm) {
	xproto.LogAlarm(b, a)
	m := mnger.Device.Model(a.DeviceNo)
	if m == nil {
		return
	}
	for _, v := range hooks {
		v.Alarm(m.Id, a)
	}
	Handler.alarm <- a
}

func EventHandler(data []byte, e *xproto.Event) {
	xproto.LogEvent(data, e)
	m := mnger.Device.Model(e.DeviceNo)
	if m == nil {
		return
	}
	for _, v := range hooks {
		v.Event(m.Id, e)
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
	v := e.Data.(xproto.EventFileLittle)
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
