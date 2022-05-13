package device

import (
	"time"
	"xstation/configs"
	"xstation/entity/hook"
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
	status chan *model.DevStatus
	alarm  chan *model.DevAlarmDetails
}

var (
	Handler *handler = &handler{
		status: make(chan *model.DevStatus, 1),
		alarm:  make(chan *model.DevAlarmDetails, 1),
	}
	hooks []Hook
)

func (o *handler) Disptah() {
	go o.dispatchStatus()
	go o.dispatchAlarm()
}

func (o *handler) Stop() {
	for _, v := range hooks {
		v.Stop()
	}
	o.alarm <- nil
	o.status <- nil
}

func NewHooks(o []hook.Option) {
	for _, v := range o {
		switch v.Name {
		case "nats":
			hooks = append(hooks, hook.NewNats(v))
		case "http://":
			hooks = append(hooks, hook.NewHttp(v))
		}
	}
}

// dispatchStatus 批量处理数据
func (o *handler) dispatchStatus() {
	dataArr := make([][]model.DevStatus, model.DevStatusTabCount)
	p, _ := ants.NewPoolWithFunc(5, func(v interface{}) {
		orm.DbCreate(v)
	}) // 协程池
	ticker := time.NewTicker(time.Second * 2)
	var (
		tabIdx uint = 0
		err    error
	)
	defer p.Release()
	defer close(o.status)
	for {
		select {
		case v := <-o.status:
			if v == nil {
				return
			}
			tabIdx = 0
			if configs.MsgProc > 0 {
				tabIdx = v.DeviceId % model.DevStatusTabCount
			}
			dataArr[tabIdx] = append(dataArr[tabIdx], *v)
		case <-ticker.C:
			for i := 0; i < model.DevStatusTabCount; i++ {
				size := len(dataArr[i])
				if size < 1 {
					continue
				}
				data := make([]model.DevStatus, size)
				copy(data, dataArr[i])
				if configs.MsgProc > 0 {
					err = p.Invoke(model.DevStatusTabVal(i, data))
				} else {
					err = p.Invoke(data)
				}
				if err != nil {
					continue
				}
				dataArr[i] = dataArr[i][:0]
			}
		}
	}
}

// dispatchAlarm 批量处理数据
func (o *handler) dispatchAlarm() {
	var stArray []model.DevAlarmDetails
	ticker := time.NewTicker(time.Second * 2)
	p, _ := ants.NewPoolWithFunc(2, func(v interface{}) {
		data := v.([]model.DevAlarmDetails)
		for _, alr := range data {
			orm.DbCreate(alr.DevStatus) // 报警
			orm.DbCreate(&alr)
			service.DevAlarmAdd(&alr)
		}
	}) // 协程池
	defer p.Release()
	defer close(o.alarm)
	for {
		select {
		case v := <-o.alarm:
			if v == nil {
				p.Invoke(stArray)
				return
			}
			if configs.MsgProc > 0 {
				stArray = append(stArray, *v)
			}
		case <-ticker.C:
			size := len(stArray)
			if size < 1 {
				continue
			}
			data := make([]model.DevAlarmDetails, size)
			copy(data, stArray)
			if err := p.Invoke(data); err != nil {
				continue
			}
			stArray = stArray[:0]
		}
	}
}
