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

var (
	nats *middleware.Nats
)

func natsRun() {
	nats = middleware.NewNatsClient()
	nats.Subscribe(topicDevOnline, natsSubDevOnlineHandler)
	nats.Subscribe(topicDevStatus, natsSubDevStatusHandler)
	nats.Subscribe(topicDevAlarm, natsSubDevAlarmHandler)
	nats.Subscribe(topicDevEvent, natsSubDevEventHandler)
}

func natsSubDevOnlineHandler(b []byte) {
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

func natsSubDevStatusHandler(b []byte) {
	var p model.DevStatus
	jsoniter.Unmarshal(b, &p)
	m := mnger.Device.Get(p.DeviceNo)
	if m == nil {
		return
	}
	p.Id = util.PrimaryKey()
	p.DeviceId = m.Id
	devtask.AddStatus(p)
}

func natsSubDevAlarmHandler(b []byte) {
	var p model.DevAlarm
	jsoniter.Unmarshal(b, &p)
	m := mnger.Device.Get(p.DeviceNo)
	if m == nil {
		return
	}
	devtask.AddAlarm(p)
}

func natsSubDevEventHandler(b []byte) {
	var e xproto.Event
	jsoniter.Unmarshal(b, &e)
	devEventHandler(&e)
}
