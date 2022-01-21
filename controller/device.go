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

// 转换
func toOnlineModel(a *xproto.Access) *model.DevOnline {
	v := &model.DevOnline{
		Guid:          a.Session,
		DeviceNo:      a.DeviceNo,
		RemoteAddress: a.RemoteAddress,
		NetType:       int(a.NetType),
		Type:          int(a.LinkType),
		UpTraffic:     a.UpTraffic,
		DownTraffic:   a.DownTraffic,
	}
	return v
}

// AccessHandler 设备接入
func (x *XNotify) AccessHandler(b []byte, a *xproto.Access) (interface{}, error) {
	log.Printf("%s\n", b)
	m := mnger.Devs.Get(a.DeviceNo)
	if m == nil {
		return nil, fmt.Errorf("deviceNo<%s> invalid", a.DeviceNo)
	}
	if a.LinkType == xproto.LINK_Signal {
		m.Version = a.Version
		m.Type = a.Type
		m.Online = true
		orm.DbUpdates(m, []string{"version", "type", "last_time", "online"})
	}
	data := toOnlineModel(a)
	data.OnTime = a.DeviceTime
	service.OnlineUpdate(data)
	if a.LinkType == xproto.LINK_FileTransfer {
		filename, act := xproto.FileOfSess(a.Session)
		if act == xproto.ACTION_Upload {
			return nil, xproto.UploadFile(a, filename, true)
		}
		return xproto.DownloadFile(configs.Default.Public+"/"+filename, nil)
	}
	return nil, nil
}

func (x *XNotify) DroppedHandler(v interface{}, a *xproto.Access, err error) {
	log.Println(err)
	m := mnger.Devs.Get(a.DeviceNo)
	if m == nil {
		return
	}
	if a.LinkType == xproto.LINK_Signal {
		m.Online = false
		orm.DbUpdates(m, []string{"last_time", "online", "last_status"})
	} else if a.LinkType == xproto.LINK_FileTransfer {
		_, act := xproto.FileOfSess(a.Session)
		if act == xproto.ACTION_Upload {
			xproto.DownloadFile("", v)
		}
	}
	data := toOnlineModel(a)
	data.OffTime = a.DeviceTime
	service.OnlineUpdate(data)
}

// StatusHandler 接收状态数据
func (x *XNotify) StatusHandler(tag string, xst *xproto.Status) {
	xproto.LogStatus(tag, xst)
	x.AddDbStatus(xst)
}

// EventHandler 事件
func (x *XNotify) EventHandler(data []byte, e *xproto.Event) {
	xproto.LogEvent(data, e)
	m := mnger.Devs.Get(e.DeviceNo)
	if m == nil || !m.AutoFtp {
		return
	}
	switch e.Type {
	case xproto.EVENT_FtpTransfer:
	case xproto.EVENT_FileLittle:
		ftpLittleFileHandler(e)
	case xproto.EVENT_FileTimedCapture:
		ftpTimedCaptureHandler(e)
	}
}

// AlarmHandler 接收报警数据
func (x *XNotify) AlarmHandler(b []byte, xalr *xproto.Alarm) {
	xproto.LogAlarm(b, xalr)
	flag := xalr.Status.Flag
	xalr.Status.Flag = 2
	if xalr.EndTime != "" {
		xalr.Status.Flag = 3
	}
	status := x.AddDbStatus(xalr.Status)
	data := model.DevAlarm{
		Guid:      xalr.UUID,
		Flag:      flag,
		StartTime: xalr.StartTime,
		EndTime:   xalr.EndTime,
		DevStatus: model.JDevStatus(*status),
		StatusId:  status.Id,
	}
	data.DTU = xalr.DTU
	data.DeviceNo = xalr.DeviceNo
	data.AlarmType = xalr.Type
	data.Data = internal.ToJString(xalr.Data)
	mnger.Alarm.Add(data)
	service.AlarmDbAdd(&data)
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

// ftpLittleFileHandler 新文件通知
func ftpLittleFileHandler(e *xproto.Event) {
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
		dbFtpFile(e, &v)
	}
}

//
func ftpTimedCaptureHandler(e *xproto.Event) {
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

// dbFtpFile ftp结果通知
func dbFtpFile(e *xproto.Event, f *xproto.EventFileLittle) {
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
	data.Channel = f.Channel
	data.Size = f.Size
	data.Duration = f.Duration
	data.FileType = f.FileType
	data.Name = f.FileName
	orm.DbCreate(data)
}
