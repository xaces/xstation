package serve

import (
	"log"
	"time"
	"xstation/app/manager"
	"xstation/models"
	"xstation/service"

	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xproto/ho"
	"github.com/wlgd/xproto/jt"
	"github.com/wlgd/xproto/ttx"

	ants "github.com/panjf2000/ants/v2"
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
	go xdata.DbStatusHandler()
	return xdata
}

// AccessHandler 设备接入
func (o *XNotify) AccessHandler(data string, access *xproto.LinkAccess) error {
	log.Printf("%s\n", data)
	if access.LinkType == xproto.LINK_Signal {
		dev := manager.Dev.Get(access.DeviceNo)
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

// ToStatusModel 转化成Model数据格式
func ToStatusModel(st *xproto.Status) (o models.XStatus) {
	dev := manager.Dev.Get(st.DeviceNo)
	if dev == nil {
		return
	}
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
	return
}

// StatusHandler 接收状态数据
func (o *XNotify) StatusHandler(tag string, xst *xproto.Status) {
	xproto.LogStatus(tag, xst)
	st := ToStatusModel(xst)
	if st.Id <= 0 {
		return
	}
	o.Status <- st
}

// AlarmHandler 接收报警数据
func (o *XNotify) AlarmHandler(tag, data string, xalr *xproto.Alarm) {
	xproto.LogAlarm(tag, data, xalr)
	status := xalr.Status.Status
	xalr.Status.Status = 2
	st := ToStatusModel(xalr.Status)
	if st.Id <= 0 {
		return
	}
	o.Status <- st
	xalr.Status.Status = status
	o.Service.DbCreateAlarm(st.Id, xalr)
}

// DbStatusHandler 批量处理数据
func (o *XNotify) DbStatusHandler() {
	stArray := make([][]models.XStatus, models.KXStatusTabNumber)
	ticker := time.NewTicker(time.Second * 1)
	p, _ := ants.NewPoolWithFunc(20, service.DbTaskFunc) // 协程池
	defer p.Release()
	for {
		select {
		case d := <-o.Status:
			tabIdx := d.TableIdx
			stArray[tabIdx] = append(stArray[tabIdx], d)
		case <-ticker.C:
			for i := 0; i < models.KXStatusTabNumber; i++ {
				size := len(stArray[i])
				if size <= 0 {
					continue
				}
				task := &service.Task{}
				task.TableIdx = i
				task.Size = size
				task.Status = make([]models.XStatus, task.Size)
				copy(task.Status, stArray[i])
				p.Invoke(task)
				stArray[i] = stArray[i][:0]
			}
		}
	}
}

func protocolAdapter(b []byte) xproto.InterfaceProto {
	if _, err := ho.IsValidProto(b); err == nil {
		return &ho.ProtoImpl{
			SubAlarmStatus: 0xffff,
			SubStatus:      0xffff,
		}
	}
	if _, err := ttx.IsValidProto(b); err == nil {
		return &ttx.ProtoImpl{}
	}
	if _, err := jt.IsValidProto(b); err == nil {
		return &jt.ProtoImpl{}
	}
	return nil
}

var (
	_xproto *xproto.Serve = nil
)

// xprotoStart 启动access服务
func xprotoStart(port uint16) {
	_xproto = &xproto.Serve{
		RequestTimeOut: 50,
		RecvTimeout:    30,
		Adapter:        protocolAdapter,
	}
	xnotify := NewXNotify()
	_xproto.SetNotifyOfLinkAccess(xnotify.AccessHandler)
	_xproto.SetNotifyOfStatus(xnotify.StatusHandler)
	_xproto.SetNotifyOfAlarm(xnotify.AlarmHandler)
	_xproto.SetNotifyOfAVFrame(xproto.LogAVFrame)
	_xproto.SetNotifyOfEvent(xproto.LogEvent)
	_xproto.SetNotifyOfRaw(xproto.LogRawFrame)
	log.Printf("XProto Serve Start %d\n", port)
	if err := _xproto.ListenAndServe(port); err != nil {
		log.Fatalln("localAccess", err)
	}
}

// xprotoStop 停止access服务
func xprotoStop() {
	_xproto.Close()
}
