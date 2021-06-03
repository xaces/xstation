package serve

import (
	"log"
	"time"
	"xstation/internal"
	"xstation/models"
	"xstation/pkg/orm"
	"xstation/service"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xproto/ho"
	"github.com/wlgd/xproto/jt"
	"github.com/wlgd/xproto/ttx"
)

// updateAccessFlow 更新流量信息
func updateAccessFlow(reg *xproto.LinkAccess) error {
	ofline := &models.XOFLine{
		Guid:          reg.Session,
		DeviceId:      reg.DeviceId,
		RemoteAddress: reg.RemoteAddress,
		AccessType:    int(reg.AccessNet),
		Type:          int(reg.LinkType),
		UpFlow:        reg.UpFlow,
		DownFlow:      reg.DownFlow,
		Version:       reg.Version,
	}
	if reg.OnLine {
		ofline.OnTime = reg.DeviceTime
	} else {
		ofline.OffTime = reg.DeviceTime
	}
	return orm.DbSave(ofline)
}

type XNotify struct {
	Status chan models.XStatus
}

// NewAccessData 实例化对象
func NewXNotify() *XNotify {
	xdata := &XNotify{
		Status: make(chan models.XStatus, 1),
	}
	go xdata.DbStatusHandler()
	return xdata
}

// AccessHandler 设备接入
func (o *XNotify) AccessHandler(data string, access *xproto.LinkAccess) error {
	log.Printf("%s\n", data)
	if access.LinkType == xproto.LINK_FileTransfer {
		if err := xproto.UploadFile(access, true); err != nil {
			return err
		}
	}
	return updateAccessFlow(access)
}

// statusToModel 转化成Model数据格式
func statusToModel(st *xproto.Status) (o models.XStatus) {
	o.Id = service.PrimaryKey()
	o.DeviceId = st.DeviceId
	o.DTU = st.DTU
	o.Status = st.Status
	if st.Location.Speed < 1 {
		st.Location.Speed = 0
	}
	o.Gps = internal.ToJString(st.Location)
	o.Tempers = internal.ToJString(st.Tempers)
	o.Humiditys = internal.ToJString(st.Humiditys)
	o.Mileage = internal.ToJString(st.Mileage)
	o.Oils = internal.ToJString(st.Oils)
	o.Module = internal.ToJString(st.Module)
	o.Gsensor = internal.ToJString(st.Gsensor)
	o.Mobile = internal.ToJString(st.Mobile)
	o.Disks = internal.ToJString(st.Disks)
	o.People = internal.ToJString(st.People)
	return
}

// StatusHandler 接收状态数据
func (o *XNotify) StatusHandler(tag string, xst *xproto.Status) {
	xproto.LogStatus(tag, xst)
	st := statusToModel(xst)
	o.Status <- st
}

// AlarmHandler 接收报警数据
func (o *XNotify) AlarmHandler(tag, data string, xalr *xproto.Alarm) {
	xproto.LogAlarm(tag, data, xalr)
	status := xalr.Status.Status
	xalr.Status.Status = 2
	st := statusToModel(xalr.Status)
	o.Status <- st
	alr := &models.XAlarm{
		DeviceId:  xalr.DeviceId,
		UUID:      xalr.UUID,
		StatusId:  st.Id,
		Status:    status,
		Type:      xalr.Type,
		StartTime: xalr.StartTime,
		EndTime:   xalr.EndTime,
		Data:      internal.ToJString(xalr.Data),
	}
	orm.DbCreate(alr)
}

// dbCreateXStatus 批量添加
func dbCreateXStatus(stlist []models.XStatus) {
	stsize := len(stlist)
	if stsize <= 0 {
		return
	}
	status := make([]models.XStatus, stsize)
	copy(status, stlist)
	go func() {
		orm.DbCreate(&status)
	}()
}

// DbStatusHandler 批量处理数据
func (o *XNotify) DbStatusHandler() {
	var stlist []models.XStatus
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case d := <-o.Status:
			stlist = append(stlist, d)
		case <-ticker.C:
			dbCreateXStatus(stlist)
			stlist = stlist[:0]
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
		RecvTimeout:    60,
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
