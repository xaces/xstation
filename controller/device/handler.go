package device

import (
	"time"
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"

	"github.com/panjf2000/ants/v2"
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

type broker interface {
	Online(*xproto.Access)
	Status(*xproto.Status)
	Alarm(*xproto.Alarm)
	Event(*xproto.Event)
	Stop()
}

type handler struct {
	status  chan *model.DevStatus
	alarm   chan *model.DevAlarm
	dsCache [][]model.DevStatus
}

var (
	Handler *handler = &handler{
		status: make(chan *model.DevStatus, 1),
		alarm:  make(chan *model.DevAlarm, 1),
	}
	brokers []broker
)

func (h *handler) Run(msgproc string) error {
	switch msgproc {
	case "nats":
		if err := natsRun(); err != nil {
			return err
		}
	case "default":
		h.Handle(&defaultBroker{})
	}
	go h.dispatchStatus()
	go h.dispatchAlarm()
	return nil
}

func (h *handler) Stop() {
	for _, v := range brokers {
		v.Stop()
	}
	h.alarm <- nil
	h.status <- nil
}

func (h *handler) Handle(v broker) {
	brokers = append(brokers, v)
}

func (h *handler) addStatus(s *xproto.Status) {
	o := devStatusModel(s)
	h.status <- o
}

func (h *handler) addAlarm(a *xproto.Alarm) {
	o := devAlarmModel(a)
	status := devStatusModel(a.Status)
	if status != nil {
		o.DevStatus = model.JDevStatus(*status)
		h.status <- status
	}
	h.alarm <- o
}

type statusObj struct {
	TableIdx int
	Data     []model.DevStatus
}

func (h *handler) insertStatus(p *ants.PoolWithFunc) {
	for i := 0; i < model.DevStatusNum; i++ {
		size := len(h.dsCache[i])
		if size < 1 {
			continue
		}
		o := &statusObj{}
		o.TableIdx = i
		o.Data = make([]model.DevStatus, size)
		copy(o.Data, h.dsCache[i])
		if err := p.Invoke(o); err != nil {
			continue
		}
		h.dsCache[i] = h.dsCache[i][:0]
	}
}

// dispatchStatus 批量处理数据
func (h *handler) dispatchStatus() {
	h.dsCache = make([][]model.DevStatus, model.DevStatusNum)
	statusFunc := func(v interface{}) {
		obj := v.(*statusObj)
		data := mnger.Device.StatusValue(obj.TableIdx, obj.Data)
		orm.DbCreate(data)
	}
	p, _ := ants.NewPoolWithFunc(5, statusFunc) // 协程池
	ticker := time.NewTicker(time.Second * 2)
	defer p.Release()
	defer close(h.status)
	for {
		select {
		case v := <-h.status:
			if v == nil {
				h.insertStatus(p)
				return
			}
			tabIdx := v.DeviceId % model.DevStatusNum
			h.dsCache[tabIdx] = append(h.dsCache[tabIdx], *v)
		case <-ticker.C:
			h.insertStatus(p)
		}
	}
}

// dispatchAlarm 批量处理数据
func (h *handler) dispatchAlarm() {
	var stArray []model.DevAlarm
	alarmFunc := func(v interface{}) {
		data := v.([]model.DevAlarm)
		for k := range data {
			service.DevAlarmDbAdd(&data[k])
			mnger.Alarm.Add(&data[k])
		}
	}
	ticker := time.NewTicker(time.Second * 2)
	p, _ := ants.NewPoolWithFunc(2, alarmFunc) // 协程池
	defer p.Release()
	defer close(h.alarm)
	for {
		select {
		case v := <-h.alarm:
			if v == nil {
				p.Invoke(stArray)
				return
			}
			// 推送给第三放
			stArray = append(stArray, *v)
		case <-ticker.C:
			size := len(stArray)
			if size < 1 {
				continue
			}
			data := make([]model.DevAlarm, size)
			copy(data, stArray)
			if err := p.Invoke(data); err != nil {
				continue
			}
			stArray = stArray[:0]
		}
	}
}
