package serve

import (
	"log"
	"time"
	"xstation/models"
	"xstation/service"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xproto/ho"
	"github.com/wlgd/xproto/jt"
	"github.com/wlgd/xproto/ttx"
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
	if access.LinkType == xproto.LINK_FileTransfer {
		if err := xproto.UploadFile(access, true); err != nil {
			return err
		}
	}
	return o.Service.DbUpdateAccess(access)
}

// StatusHandler 接收状态数据
func (o *XNotify) StatusHandler(tag string, xst *xproto.Status) {
	xproto.LogStatus(tag, xst)
	st := o.Service.ToStatusModel(xst)
	o.Status <- st
}

// AlarmHandler 接收报警数据
func (o *XNotify) AlarmHandler(tag, data string, xalr *xproto.Alarm) {
	xproto.LogAlarm(tag, data, xalr)
	status := xalr.Status.Status
	xalr.Status.Status = 2
	st := o.Service.ToStatusModel(xalr.Status)
	o.Status <- st
	xalr.Status.Status = status
	o.Service.DbCreateAlarm(st.Id, xalr)
}

// DbStatusHandler 批量处理数据
func (o *XNotify) DbStatusHandler() {
	var stArray []models.XStatus
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case d := <-o.Status:
			stArray = append(stArray, d)
		case <-ticker.C:
			o.Service.DbCreateStatus(stArray)
			stArray = stArray[:0]
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
