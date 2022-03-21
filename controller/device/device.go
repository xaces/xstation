package device

import (
	"fmt"
	"log"
	"xstation/configs"
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/util"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

func onlineHandler(a *xproto.Access) error {
	for _, v := range brokers {
		v.Online(a)
	}
	return nil
}

func AccessHandler(b []byte, a *xproto.Access) (interface{}, error) {
	log.Printf("%s\n", b)
	m := mnger.Device.Model(a.DeviceNo)
	if m == nil {
		return nil, fmt.Errorf("[%s] invalid", a.DeviceNo)
	}
	switch a.LinkType {
	case xproto.LINK_FileTransfer:
		filename, act := xproto.FileOfSess(a.Session)
		if act == xproto.ACTION_Upload {
			return nil, xproto.UploadFile(a, filename, true)
		}
		return xproto.DownloadFile(configs.Default.Public+"/"+filename, nil)
	case xproto.LINK_Signal:
		return nil, onlineHandler(a)
	}
	return nil, xproto.ErrUnSupport

}

func DroppedHandler(v interface{}, a *xproto.Access, err error) {
	log.Println(err)
	switch a.LinkType {
	case xproto.LINK_FileTransfer:
		_, act := xproto.FileOfSess(a.Session)
		if act == xproto.ACTION_Upload {
			xproto.DownloadFile("", v)
		}
	case xproto.LINK_Signal:
		onlineHandler(a)
	}
}
func StatusHandler(tag string, s *xproto.Status) {
	xproto.LogStatus(tag, s)
	for _, v := range brokers {
		v.Status(s)
	}
}
func AlarmHandler(b []byte, a *xproto.Alarm) {
	xproto.LogAlarm(b, a)
	for _, v := range brokers {
		v.Alarm(a)
	}
}
func EventHandler(data []byte, e *xproto.Event) {
	xproto.LogEvent(data, e)
	for _, v := range brokers {
		v.Event(e)
	}
}

func devEventHandler(e *xproto.Event) {
	m := mnger.Device.Model(e.DeviceNo)
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
