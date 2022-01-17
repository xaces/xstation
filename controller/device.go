package controller

import (
	"fmt"
	"log"
	"time"
	"xstation/app/mnger"
	"xstation/configs"
	"xstation/internal"
	"xstation/model"
	"xstation/service"

	"github.com/panjf2000/ants/v2"
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

type XNotify struct {
	Status chan model.DevStatus
}

// NewXNotify 实例化对象
func NewXNotify() *XNotify {
	xdata := &XNotify{
		Status: make(chan model.DevStatus),
	}
	go xdata.DbStatusHandler()
	return xdata
}

// AddDbStatus 转化成Model数据格式
func (x *XNotify) AddDbStatus(st *xproto.Status) *model.DevStatus {
	dev := mnger.Devs.Get(st.DeviceNo)
	if dev == nil {
		return nil
	}
	o := model.DevStatus{}
	o.Id = service.PrimaryKey()
	o.DeviceId = dev.Id
	o.DeviceNo = st.DeviceNo
	o.DTU = st.DTU
	o.Flag = st.Flag
	if st.Location.Speed < 1 {
		st.Location.Speed = 0
	}
	o.Acc = st.Acc
	o.Location = model.JLocation(st.Location)
	o.Tempers = model.JFloats(st.Tempers)
	o.Humiditys = model.JFloats(st.Humiditys)
	o.Mileage = model.JMileage(st.Mileage)
	o.Oils = model.JOil(st.Oils)
	o.Module = model.JModule(st.Module)
	o.Gsensor = model.JGsensor(st.Gsensor)
	o.Mobile = model.JMobile(st.Mobile)
	o.Disks = model.JDisks(st.Disks)
	o.People = model.JPeople(st.People)
	x.Status <- o
	dev.LastStatus = model.JDevStatus(o)
	return &o
}

// 上线下线处理
func onlineHandler(x *xproto.Access) error {
	ofline := &model.DevOnline{
		Guid:          x.Session,
		DeviceNo:      x.DeviceNo,
		RemoteAddress: x.RemoteAddress,
		NetType:       int(x.NetType),
		Type:          int(x.LinkType),
		UpTraffic:     x.UpTraffic,
		DownTraffic:   x.DownTraffic,
	}
	if x.Online {
		ofline.OnTime = x.DeviceTime
	} else {
		ofline.OffTime = x.DeviceTime
	}
	return service.OnlineUpdate(ofline)
}

// AccessHandler 设备接入
func (o *XNotify) AccessHandler(data []byte, arg *interface{}, x *xproto.Access) error {
	log.Printf("%s\n", string(data))
	m := mnger.Devs.Get(x.DeviceNo)
	if m == nil {
		return fmt.Errorf("deviceNo:%s invalid", x.DeviceNo)
	}
	if x.LinkType == xproto.LINK_Signal {
		m.Version = x.Version
		m.Type = x.Type
		m.Online = x.Online
		fields := []string{"version", "type", "last_time", "online"}
		if !m.Online {
			fields = append(fields, "last_status")
		}
		orm.DbUpdates(m, fields)
	} else if x.LinkType == xproto.LINK_FileTransfer {
		filename, act := xproto.FileOfSess(x.Session)
		if act == xproto.ACTION_Upload {
			xproto.UploadFile(x, filename, true)
		} else {
			xproto.DownloadFile(x, filename, arg)
		}
	}
	return onlineHandler(x)
}

// StatusHandler 接收状态数据
func (x *XNotify) StatusHandler(tag string, xst *xproto.Status) {
	xproto.LogStatus(tag, xst)
	x.AddDbStatus(xst)
}

// EventHandler 事件
func (x *XNotify) EventHandler(data []byte, e *xproto.Event) {
	xproto.LogEvent(data, e)
	switch e.Type {
	case xproto.EVENT_FtpTransfer:
	case xproto.EVENT_FileLittle:
		littleFileHandler(e)
	case xproto.EVENT_FileTimedCapture:
		timedCaptureHandler(e)
	}
}

// AlarmHandler 接收报警数据
func (x *XNotify) AlarmHandler(data []byte, xalr *xproto.Alarm) {
	xproto.LogAlarm(data, xalr)
	flag := xalr.Status.Flag
	xalr.Status.Flag = 2
	if xalr.EndTime != "" {
		xalr.Status.Flag = 3
	}
	status := x.AddDbStatus(xalr.Status)
	alarm := model.DevAlarm{
		Guid:      xalr.UUID,
		Flag:      flag,
		StartTime: xalr.StartTime,
		EndTime:   xalr.EndTime,
		DevStatus: model.JDevStatus(*status),
		StatusId:  status.Id,
	}
	alarm.DTU = xalr.DTU
	alarm.DeviceNo = xalr.DeviceNo
	alarm.AlarmType = xalr.Type
	alarm.Data = internal.ToJString(xalr.Data)
	mnger.Alarm.Add(alarm)
	service.AlarmDbAdd(&alarm)
}

// DbStatusHandler 批量处理数据
func (x *XNotify) DbStatusHandler() {
	stArray := make([][]model.DevStatus, model.DevStatusTabs)
	ticker := time.NewTicker(time.Second * 2)
	p, _ := ants.NewPoolWithFunc(5, service.StatusCreates) // 协程池
	defer p.Release()
	defer close(x.Status)
	for {
		select {
		case d := <-x.Status:
			tabIdx := d.DeviceId % model.DevStatusTabs
			stArray[tabIdx] = append(stArray[tabIdx], d)
		case <-ticker.C:
			for i := 0; i < model.DevStatusTabs; i++ {
				if err := x.DbInsertStatus(p, i, stArray[i]); err != nil {
					continue
				}
				stArray[i] = stArray[i][:0]
			}
		}
	}
}

// DbInsertStatus
func (x *XNotify) DbInsertStatus(p *ants.PoolWithFunc, tabIdx int, data []model.DevStatus) error {
	size := len(data)
	if size <= 0 {
		return nil
	}
	task := &service.StatusTask{}
	task.TableIdx = tabIdx
	task.Size = size
	task.Data = make([]model.DevStatus, task.Size)
	copy(task.Data, data)
	return p.Invoke(task)
}

// littleFileHandler 新文件通知
func littleFileHandler(e *xproto.Event) {
	m := mnger.Devs.Get(e.DeviceNo)
	if m == nil || !m.AutoFtp {
		return
	}
	v := e.Data.(xproto.EventFileLittle)
	dst := internal.StringIndex(v.FileName, "/", 3)
	fpName := fmt.Sprintf("%s/%s", e.DeviceNo, dst)
	ftp := xproto.FtpTransfer{
		FtpURL:   configs.Default.Ftp.Address,
		FileSrc:  v.FileName,
		FileDst:  fpName,
		Action:   xproto.ACTION_Download,
		FileType: v.FileType,
		Session:  e.Session,
	}
	if xproto.SyncSend(xproto.REQ_FtpTransfer, ftp, nil, e.DeviceNo) == nil {
		v.FileName = fpName
		ftpFileHandler(e, &v)
	}
}

//
func timedCaptureHandler(e *xproto.Event) {
	m := mnger.Devs.Get(e.DeviceNo)
	if m == nil || !m.AutoFtp {
		return
	}
	v := e.Data.(xproto.EventFileTimedCapture)
	dst := internal.StringIndex(v.FileName, "/", 3)
	fpName := fmt.Sprintf("%s/%s", e.DeviceNo, dst)
	ftp := xproto.FtpTransfer{
		FtpURL:   configs.Default.Ftp.Address,
		FileSrc:  v.FileName,
		FileDst:  fpName,
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
		Name:      fpName,
	}
	orm.DbCreate(data)
}

// ftpFileHandler ftp结果通知
func ftpFileHandler(e *xproto.Event, f *xproto.EventFileLittle) {
	alr := mnger.Alarm.Get(e.Session) // 从缓存中获取数据
	if alr == nil {
		return
	}
	fl := &model.DevAlarmFile{}
	fl.Guid = e.Session
	fl.LinkType = model.AlarmLinkFtpFile
	fl.DeviceNo = e.DeviceNo
	fl.AlarmType = alr.AlarmType
	fl.DTU = e.DTU
	fl.Channel = f.Channel
	fl.Size = f.Size
	fl.Duration = f.Duration
	fl.FileType = f.FileType
	fl.Name = f.FileName
	orm.DbCreate(fl)
}
