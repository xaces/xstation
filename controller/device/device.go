package device

import (
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
		if m == nil {
			m = cache.NewDevice(cache.Vehicle{DeviceNo: a.DeviceNo})
		}
		log.Printf("%s", b)
		if m.DeviceId == 0 {
			return nil, nil
			// return nil, fmt.Errorf("[%s] invalid", a.DeviceNo)
		}
		m.Update(a)
		for _, v := range hooks {
			v.Online(m.DeviceId, a)
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
		if m.DeviceId == 0 {
			return
		}
		for _, v := range hooks {
			v.Online(m.DeviceId, a)
		}
		m.Update(a)
		devOnlineUpdate(a, &m.Status)
	}
}

func StatusHandler(tag string, s *xproto.Status) {
	xproto.LogStatus(tag, s)
	m := cache.Device(s.DeviceNo)
	if s.Flag == 0 {
		m.LastOnlineTime = s.DTU
		m.Status = *s
	}
	if m.DeviceId == 0 {
		return
	}
	for _, v := range hooks {
		v.Status(m.DeviceId, s)
	}
	o := devStatusModel(s)
	o.DeviceId = m.DeviceId
	Handler.status <- o
}

func AlarmHandler(b []byte, a *xproto.Alarm) {
	xproto.LogAlarm(b, a)
	m := cache.Device(a.DeviceNo)
	if m.DeviceId == 0 {
		return
	}
	for _, v := range hooks {
		v.Alarm(m.DeviceId, a)
	}
	Handler.alarm <- a
}

func EventHandler(data []byte, e *xproto.Event) {
	xproto.LogEvent(data, e)
	m := cache.Device(e.DeviceNo)
	if m.DeviceId == 0 {
		return
	}
	for _, v := range hooks {
		v.Event(m.DeviceId, e)
	}
	switch e.Type {
	case xproto.Event_FtpTransfer:
	case xproto.Event_FileLittle:
		alr := cache.DevAlarm(e.Session) // 从缓存中获取数据
		if alr == nil || !m.FtpAlarms.Include(alr.AlarmType) {
			return
		}
		ftpLittleFile(e, alr.AlarmType)
	case xproto.Event_FileTimedCapture:
		ftpTimedCapture(e)
	}
}

// ftpLittleFile 新文件通知
func ftpLittleFile(e *xproto.Event, alrType int) {
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
	data := &model.DevAlarmFile{}
	data.Guid = e.Session
	data.LinkType = model.AlarmLinkFtpFile
	data.DeviceNo = e.DeviceNo
	data.AlarmType = alrType
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
