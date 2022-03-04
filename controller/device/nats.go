package device

import (
	"xstation/entity/mnger"
	"xstation/middleware"
	"xstation/model"
	"xstation/service"
	"xstation/util"

	jsoniter "github.com/json-iterator/go"
	"github.com/wlgd/xproto"
)

const (
	topicDevOnline = "station.device.online"
	topicDevStatus = "station.device.status"
	topicDevAlarm  = "station.device.alarm"
	topicDevEvent  = "station.device.event"
)

type natsHandler struct {
	client *middleware.Nats
}

func (n *natsHandler) Run() {
	n.client = middleware.NewNatsClient()
	n.client.Subscribe(topicDevOnline, n.SubDevOnlineHandler)
	n.client.Subscribe(topicDevStatus, n.SubDevStatusHandler)
	n.client.Subscribe(topicDevAlarm, n.SubDevAlarmHandler)
	n.client.Subscribe(topicDevEvent, n.SubDevEventHandler)
}

func (n *natsHandler) Notify(topic string, v interface{}) {
	if n.client == nil {
		return
	}
	n.client.Notify(topic, v)
}

func (n *natsHandler) SubDevOnlineHandler(b []byte) {
	var p model.DevOnline
	jsoniter.Unmarshal(b, &p)
	m := mnger.Device.Get(p.DeviceNo)
	if m == nil {
		return
	}
	online := true
	if p.OffTime != "" {
		online = false
	}
	service.DeviceUpdate(m, online, p.Version, p.DevType)
	service.DevOnlineUpdate(&p)
}

func (n *natsHandler) SubDevStatusHandler(b []byte) {
	var p model.DevStatus
	jsoniter.Unmarshal(b, &p)
	m := mnger.Device.Get(p.DeviceNo)
	if m == nil {
		return
	}
	p.Id = util.PrimaryKey()
	p.DeviceId = m.Id
	Handler.AddStatus(p)
}

func (n *natsHandler) SubDevAlarmHandler(b []byte) {
	var p model.DevAlarm
	jsoniter.Unmarshal(b, &p)
	m := mnger.Device.Get(p.DeviceNo)
	if m == nil {
		return
	}
	Handler.AddAlarm(p)
}

func (n *natsHandler) SubDevEventHandler(b []byte) {
	var e xproto.Event
	jsoniter.Unmarshal(b, &e)
	devEventHandler(&e)
}
