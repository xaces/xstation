package device

import (
	"xstation/middleware"

	jsoniter "github.com/json-iterator/go"
	"github.com/wlgd/xproto"
)

const (
	topicDevOnline = "station.device.online"
	topicDevStatus = "station.device.status"
	topicDevAlarm  = "station.device.alarm"
	topicDevEvent  = "station.device.event"
)

func natsRun() error {
	nats := &middleware.Nats{}
	if err := nats.Start(middleware.NatsOption{
		TopicOnline: topicDevOnline,
		TopicStatus: topicDevStatus,
		TopicAlarm:  topicDevAlarm,
		TopicEvent:  topicDevEvent,
	}); err != nil {
		return err
	} // 推送中间件
	Handler.Handle(nats)
	nats.Conn.Subscribe(topicDevOnline, natsOnlineHandler)
	nats.Conn.Subscribe(topicDevStatus, natsStatusHandler)
	nats.Conn.Subscribe(topicDevAlarm, natsAlarmHandler)
	nats.Conn.Subscribe(topicDevEvent, natsEventHandler)
	return nil
}

func natsOnlineHandler(b []byte) {
	var p xproto.Access
	jsoniter.Unmarshal(b, &p)
	updateDevOnline(&p)
}

func natsStatusHandler(b []byte) {
	var p xproto.Status
	jsoniter.Unmarshal(b, &p)
	Handler.addStatus(&p)
}

func natsAlarmHandler(b []byte) {
	var p xproto.Alarm
	jsoniter.Unmarshal(b, &p)
	Handler.addAlarm(&p)
}

func natsEventHandler(b []byte) {
	var e xproto.Event
	jsoniter.Unmarshal(b, &e)
	devEventHandler(&e)
}
