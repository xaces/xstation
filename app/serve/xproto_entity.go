package serve

import (
	"log"
	"time"
	"xstation/app/mnger"
	"xstation/internal"
	"xstation/models"
	"xstation/service"

	"github.com/panjf2000/ants/v2"
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

type XNotify struct {
	Status  chan models.XStatus
	Service *service.XData
}

// NewAccessData 实例化对象
func NewXNotify() *XNotify {
	xdata := &XNotify{
		Status:  make(chan models.XStatus, 1),
		Service: service.NewXData(),
	}
	go xdata.DbInsertHandler()
	return xdata
}

// AddDbStatus 转化成Model数据格式
func (x *XNotify) AddDbStatus(st *xproto.Status) int64 {
	dev := mnger.Dev.Get(st.DeviceNo)
	if dev == nil {
		return 0
	}
	o := models.XStatus{}
	o.Id = service.PrimaryKey()
	o.DeviceId = dev.Id
	o.DeviceNo = st.DeviceNo
	o.DTU = st.DTU
	o.Status = st.Status
	if st.Location.Speed < 1 {
		st.Location.Speed = 0
	}
	o.Acc = st.Acc
	o.Gps = models.JGps(st.Location)
	o.Tempers = models.JFloats(st.Tempers)
	o.Humiditys = models.JFloats(st.Humiditys)
	o.Mileage = models.JMileage(st.Mileage)
	o.Oils = models.JOil(st.Oils)
	o.Module = models.JModule(st.Module)
	o.Gsensor = models.JGsensor(st.Gsensor)
	o.Mobile = models.JMobile(st.Mobile)
	o.Disks = models.JDisks(st.Disks)
	o.People = models.JPeople(st.People)
	o.TableIdx = int(dev.Id) % models.KXStatusTabNumber
	x.Status <- o
	return o.Id
}

// AccessHandler 设备接入
func (o *XNotify) AccessHandler(data string, access *xproto.LinkAccess) error {
	log.Printf("%s\n", data)
	if access.LinkType == xproto.LINK_Signal {
		dev := mnger.Dev.Get(access.DeviceNo)
		if dev == nil {
			return xproto.ErrInvalidDevice
		}
		if dev.Version != access.Version || dev.Type != access.Type {
			dev.Version = access.Version
			dev.Type = access.Type
			orm.DbUpdateModel(dev)
		}
	}
	if access.LinkType == xproto.LINK_FileTransfer {
		if err := xproto.UploadFile(access, true); err != nil {
			return err
		}
	}
	return o.Service.DbUpdateAccess(access)
}

// StatusHandler 接收状态数据
func (x *XNotify) StatusHandler(tag string, xst *xproto.Status) {
	xproto.LogStatus(tag, xst)
	x.AddDbStatus(xst)
}

// AlarmHandler 接收报警数据
func (x *XNotify) AlarmHandler(tag, data string, xalr *xproto.Alarm) {
	xproto.LogAlarm(tag, data, xalr)
	status := xalr.Status.Status
	xalr.Status.Status = 2
	statusId := x.AddDbStatus(xalr.Status)
	if statusId <= 0 {
		return
	}
	alarm := models.XAlarm{
		DeviceNo:  xalr.DeviceNo,
		Guid:      internal.UUID(),
		UUID:      xalr.UUID,
		StatusId:  statusId,
		Status:    status,
		Type:      xalr.Type,
		StartTime: xalr.StartTime,
		EndTime:   xalr.EndTime,
		Data:      internal.ToJString(xalr.Data),
	}
	orm.DbCreate(&alarm)
}

// DbStatusHandler 批量处理数据
func (x *XNotify) DbInsertHandler() {
	stArray := make([][]models.XStatus, models.KXStatusTabNumber)
	ticker := time.NewTicker(time.Second * 2)
	p, _ := ants.NewPoolWithFunc(5, service.DbStatusTaskFunc) // 协程池
	defer p.Release()
	for {
		select {
		case d := <-x.Status:
			tabIdx := d.TableIdx
			stArray[tabIdx] = append(stArray[tabIdx], d)
		case <-ticker.C:
			for i := 0; i < models.KXStatusTabNumber; i++ {
				if err := x.DbInsertStatus(p, i, stArray[i]); err != nil {
					continue
				}
				stArray[i] = stArray[i][:0]
			}
		}
	}
}

// insertStatus
func (x *XNotify) DbInsertStatus(p *ants.PoolWithFunc, tabIdx int, data []models.XStatus) error {
	size := len(data)
	if size <= 0 {
		return nil
	}
	task := &service.StatusTask{}
	task.TableIdx = tabIdx
	task.Size = size
	task.Data = make([]models.XStatus, task.Size)
	copy(task.Data, data)
	return p.Invoke(task)
}
