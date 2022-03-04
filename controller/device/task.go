package device

import (
	"time"
	"unsafe"
	"xstation/configs"
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"
	"xstation/util"

	"github.com/panjf2000/ants/v2"
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

var (
	devtask *devTask
)

type statusTask struct {
	TableIdx int
	Data     []model.DevStatus
}

// devStatusBatch 批量添加
func devStatusBatch(obj interface{}) {
	task := obj.(*statusTask)
	// 映射
	var data interface{}
	ptr := unsafe.Pointer(&task.Data)
	switch task.TableIdx {
	case 1:
		data = (*[]model.DevStatus1)(ptr)
	case 2:
		data = (*[]model.DevStatus2)(ptr)
	case 3:
		data = (*[]model.DevStatus3)(ptr)
	case 4:
		data = (*[]model.DevStatus4)(ptr)
	default:
		data = &task.Data
	}
	orm.DbCreate(data)
}

type devTask struct {
	Status chan model.DevStatus
	Alarm  chan model.DevAlarm
	Online chan model.DevOnline
}

func devTaskRun() {
	devtask = &devTask{
		Status: make(chan model.DevStatus, 1),
		Alarm:  make(chan model.DevAlarm, 1),
		Online: make(chan model.DevOnline, 1),
	}
	go devtask.StatusDispatch()
	go devtask.AlarmDispatch()
}

func (d *devTask) AddStatus(v model.DevStatus) {
	d.Status <- v
}

func (d *devTask) AddAlarm(v model.DevAlarm) {
	d.Alarm <- v
}

func (d *devTask) AddOnline(v model.DevOnline) {
	d.Online <- v
}

// DbStatusHandler 批量处理数据
func (d *devTask) StatusDispatch() {
	stArray := make([][]model.DevStatus, model.DevStatusNum)
	ticker := time.NewTicker(time.Second * 2)
	p, _ := ants.NewPoolWithFunc(5, devStatusBatch) // 协程池
	defer p.Release()
	defer close(d.Status)
	for {
		select {
		case v := <-d.Status:
			tabIdx := v.DeviceId % model.DevStatusNum
			stArray[tabIdx] = append(stArray[tabIdx], v)
		case <-ticker.C:
			for i := 0; i < model.DevStatusNum; i++ {
				size := len(stArray[i])
				if size < 1 {
					continue
				}
				task := &statusTask{}
				task.TableIdx = i
				task.Data = make([]model.DevStatus, size)
				copy(task.Data, stArray[i])
				if err := p.Invoke(task); err != nil {
					continue
				}
				stArray[i] = stArray[i][:0]
			}
		}
	}
}

// AlarmDispatch 批量处理数据
func (d *devTask) AlarmDispatch() {
	var stArray []model.DevAlarm
	alarmFunc := func(v interface{}) {
		data := v.([]model.DevAlarm)
		for k := range data {
			service.DevAlarmDbAdd(&data[k])
			mnger.Alarm.Add(&data[k])
		}
	}
	ticker := time.NewTicker(time.Second * 2)
	p, _ := ants.NewPoolWithFunc(2, alarmFunc) // 协程池
	defer p.Release()
	defer close(d.Alarm)
	for {
		select {
		case v := <-d.Alarm:
			// 推送给第三放
			stArray = append(stArray, v)
		case <-ticker.C:
			size := len(stArray)
			if size < 1 {
				continue
			}
			data := make([]model.DevAlarm, size)
			copy(data, stArray)
			if err := p.Invoke(data); err != nil {
				continue
			}
			stArray = stArray[:0]
		}
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
	fpName := util.FilePath(v.FileName, e.DeviceNo)
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
