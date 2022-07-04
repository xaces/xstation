package hook

import "github.com/xaces/xproto"

type Interface interface {
	PublishOnline(*xproto.Access)
	PublishStatus(*xproto.Status)
	PublishAlarm(*xproto.Alarm)
	PublishEvent(*xproto.Event)
	Stop()
}

type Option struct {
	Address string
	Online  string
	Status  string
	Alarm   string
	Event   string
}

type online struct {
	Code int `json:"msgCode"`
	*xproto.Access
}

type status struct {
	Code int `json:"msgCode"`
	*xproto.Status
}

type alarm struct {
	Code int `json:"msgCode"`
	*xproto.Alarm
}

type event struct {
	Code int `json:"msgCode"`
	*xproto.Event
}
