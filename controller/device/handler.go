package device

import (
	"time"
	"xstation/entity/hook"
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"

	"github.com/panjf2000/ants/v2"
	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

type Hook interface {
	Online(deviceId uint, a *xproto.Access)
	Status(deviceId uint, s *xproto.Status)
	Alarm(deviceId uint, a *xproto.Alarm)
	Event(deviceId uint, e *xproto.Event)
	Stop()
}

type handler struct {
	status  chan *model.DevStatus
	alarm   chan *xproto.Alarm
	dsCache [][]model.DevStatus
}

var (
	Handler *handler = &handler{
		status: make(chan *model.DevStatus, 1),
		alarm:  make(chan *xproto.Alarm, 1),
	}
	hooks []Hook
)

func (h *handler) Disptah() {
	go h.dispatchStatus()
	go h.dispatchAlarm()
}

func (h *handler) Stop() {
	for _, v := range hooks {
		v.Stop()
	}
	h.alarm <- nil
	h.status <- nil
}

func Hooks(o []hook.Option) {
	for _, v := range o {
		switch v.Name {
		case "nats":
			hooks = append(hooks, hook.NewNats(v))
		case "http://":
			hooks = append(hooks, hook.NewHttp(v))
		}
	}
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
	var stArray []xproto.Alarm
	alarmFunc := func(v interface{}) {
		data := v.([]xproto.Alarm)
		for _, alr := range data {
			s := devStatusModel(alr.Status)
			if s != nil {
				h.status <- s
			}
			orm.DbCreate(devAlarmDetailsModel(&alr, s))
			o := devAlarmModel(&alr, s)
			service.DevAlarmAdd(o)
			mnger.Alarm.Add(o)
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
			data := make([]xproto.Alarm, size)
			copy(data, stArray)
			if err := p.Invoke(data); err != nil {
				continue
			}
			stArray = stArray[:0]
		}
	}
}
