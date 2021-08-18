package middleware

import (
	"encoding/json"
	"log"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/mq"
)

var (
	_nclient *mq.NatsClient
)

func PublishAlarm(data interface{}) {
	if _nclient == nil {
		return
	}
	_nclient.Publish("xalarm", data)
}

func PublishStatus(data interface{}) {
	if _nclient == nil {
		return
	}
	_nclient.Publish("xstatus", data)
}

func realtimeAlarmHandler(data []byte) {
	var alr xproto.Alarm
	json.Unmarshal(data, &alr)
	log.Printf("%v\n", alr)
}

func realtimeStatusHandler(data []byte) {
	var alr xproto.Status
	json.Unmarshal(data, &alr)
	log.Printf("%v\n", alr)
}

func MqttStart() error {
	client, err := mq.NewNatsClient(mq.DefaultURL, false)
	if err != nil {
		return err
	}
	_nclient = client
	_nclient.Subscribe("xalarm", realtimeAlarmHandler)
	_nclient.Subscribe("xstatus", realtimeStatusHandler)
	return nil
}

func MqttStop() {
	if _nclient == nil {
		return
	}
	_nclient.Release()
}
