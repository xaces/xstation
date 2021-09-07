package app

import (
	"log"
	"time"
	"xstation/internal"
	"xstation/mnger"
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
	go xdata.DbInsertHandler()
	return xdata
}

// AddDbStatus 转化成Model数据格式
func (x *XNotify) AddDbStatus(st *xproto.Status) uint64 {
	dev := mnger.Dev.Get(st.DeviceNo)
	if dev == nil {
		return 0
	}
	o := model.DevStatus{}
	o.Id = service.PrimaryKey()
	o.DeviceId = dev.Id
	o.DeviceNo = st.DeviceNo
	o.DTU = st.DTU
	o.Flag = st.Flag
	if st.Gps.Speed < 1 {
		st.Gps.Speed = 0
	}
	o.Acc = st.Acc
	o.Gps = model.JGps(st.Gps)
	o.Tempers = model.JFloats(st.Tempers)
	o.Humiditys = model.JFloats(st.Humiditys)
	o.Mileage = model.JMileage(st.Mileage)
	o.Oils = model.JOil(st.Oils)
	o.Module = model.JModule(st.Module)
	o.Gsensor = model.JGsensor(st.Gsensor)
	o.Mobile = model.JMobile(st.Mobile)
	o.Disks = model.JDisks(st.Disks)
	o.People = model.JPeople(st.People)
	o.TableIdx = int(dev.Id) % service.StatusTableNum
	x.Status <- o
	dev.LastStatus = model.JDevStatus(o)
	return o.Id
}

// AccessHandler 设备接入
func (o *XNotify) AccessHandler(data []byte, arg *interface{}, x *xproto.Access) error {
	log.Printf("%s\n", string(data))
	m := mnger.Dev.Get(x.DeviceNo)
	if m == nil {
		return xproto.ErrInvalidDevice
	}
	if x.LinkType == xproto.LINK_Signal {
		m.Version = x.Version
		m.Type = x.Type
		m.Online = x.OnLine
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
	return service.DbUpdateOnline(x)
}

// StatusHandler 接收状态数据
func (x *XNotify) StatusHandler(tag string, xst *xproto.Status) {
	xproto.LogStatus(tag, xst)
	x.AddDbStatus(xst)
}

// AlarmHandler 接收报警数据
func (x *XNotify) AlarmHandler(data []byte, xalr *xproto.Alarm) {
	xproto.LogAlarm(data, xalr)
	flag := xalr.Status.Flag
	xalr.Status.Flag = 2
	statusId := x.AddDbStatus(xalr.Status)
	if statusId <= 0 {
		return
	}
	alarm := model.DevAlarm{
		DeviceNo:  xalr.DeviceNo,
		Guid:      internal.UUID(),
		UUID:      xalr.UUID,
		StatusId:  statusId,
		Flag:      flag,
		Type:      xalr.Type,
		StartTime: xalr.StartTime,
		EndTime:   xalr.EndTime,
		Data:      internal.ToJString(xalr.Data),
	}
	orm.DbCreate(&alarm)
}

// DbStatusHandler 批量处理数据
func (x *XNotify) DbInsertHandler() {
	stArray := make([][]model.DevStatus, service.StatusTableNum)
	ticker := time.NewTicker(time.Second * 2)
	p, _ := ants.NewPoolWithFunc(5, service.DbStatusTaskFunc) // 协程池
	defer p.Release()
	defer close(x.Status)
	for {
		select {
		case d := <-x.Status:
			tabIdx := d.TableIdx
			stArray[tabIdx] = append(stArray[tabIdx], d)
		case <-ticker.C:
			for i := 0; i < service.StatusTableNum; i++ {
				if err := x.DbInsertStatus(p, i, stArray[i]); err != nil {
					continue
				}
				stArray[i] = stArray[i][:0]
			}
		}
	}
}

// insertStatus
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
