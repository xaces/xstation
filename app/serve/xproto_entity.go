package serve

import (
	"log"
	"time"
	"xstation/app/mnger"
	"xstation/internal"
	"xstation/model"
	"xstation/service"

	"github.com/panjf2000/ants/v2"
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

type XNotify struct {
	Status chan model.Status
}

// NewAccessData 实例化对象
func NewXNotify() *XNotify {
	xdata := &XNotify{
		Status: make(chan model.Status),
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
	o := model.Status{}
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
	return o.Id
}

// AccessHandler 设备接入
func (o *XNotify) AccessHandler(data string, access *xproto.LinkAccess) error {
	log.Printf("%s\n", data)
	m := mnger.Dev.Get(access.DeviceNo)
	if m == nil {
		return xproto.ErrInvalidDevice
	}
	if access.LinkType == xproto.LINK_Signal {
		if m.Version != access.Version || m.Type != access.Type {
			m.Version = access.Version
			m.Type = access.Type
			orm.DbUpdateModel(m)
		}
	} else if access.LinkType == xproto.LINK_FileTransfer {
		if err := xproto.UploadFile(access, true); err != nil {
			return err
		}
	}
	return service.DbUpdateOnline(access)
}

// StatusHandler 接收状态数据
func (x *XNotify) StatusHandler(tag string, xst *xproto.Status) {
	xproto.LogStatus(tag, xst)
	x.AddDbStatus(xst)
}

// AlarmHandler 接收报警数据
func (x *XNotify) AlarmHandler(tag, data string, xalr *xproto.Alarm) {
	xproto.LogAlarm(tag, data, xalr)
	flag := xalr.Status.Flag
	xalr.Status.Flag = 2
	statusId := x.AddDbStatus(xalr.Status)
	if statusId <= 0 {
		return
	}
	alarm := model.Alarm{
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
	stArray := make([][]model.Status, service.StatusTableNum)
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
func (x *XNotify) DbInsertStatus(p *ants.PoolWithFunc, tabIdx int, data []model.Status) error {
	size := len(data)
	if size <= 0 {
		return nil
	}
	task := &service.StatusTask{}
	task.TableIdx = tabIdx
	task.Size = size
	task.Data = make([]model.Status, task.Size)
	copy(task.Data, data)
	return p.Invoke(task)
}
