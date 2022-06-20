package device

import (
	"strings"
	"time"
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
		status: make(chan *model.DevStatus, 50),
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
		if strings.Contains(v.Address, "nats://") {
			hooks = append(hooks, hook.NewNats(v))
		} else if strings.Contains(v.Address, "http://") {
			hooks = append(hooks, hook.NewHttp(v))
		}
	}
}

// dispatchStatus 批量处理数据
func (o *handler) dispatchStatus() {
	tableCount := model.DevStatus{}.TableCount()
	dataArr := make([][]model.DevStatus, tableCount)
	p, _ := ants.NewPoolWithFunc(5, func(v interface{}) {
		data := v.([]model.DevStatus)
		orm.Table(data[0]).Create(&data)
	}) // 协程池
	ticker := time.NewTicker(time.Second * 2)
	defer p.Release()
	defer close(o.status)
	for {
		select {
		case v := <-o.status:
			if v == nil {
				return
			}
			tabIdx := v.DeviceID % tableCount
			dataArr[tabIdx] = append(dataArr[tabIdx], *v)
		case <-ticker.C:
			for i := 0; i < int(tableCount); i++ {
				size := len(dataArr[i])
				if size < 1 {
					continue
				}
				data := make([]model.DevStatus, size)
				copy(data, dataArr[i])
				if err := p.Invoke(data); err != nil {
					continue
				}
				dataArr[i] = dataArr[i][:0]
			}
		}
	}
}

// dispatchAlarm 批量处理数据
func (o *handler) dispatchAlarm() {
	var dataArr []model.DevAlarmDetails
	ticker := time.NewTicker(time.Second * 2)
	p, _ := ants.NewPoolWithFunc(10, func(v interface{}) {
		data := v.([]model.DevAlarmDetails)
		orm.DbCreate(&data)
		for _, alr := range data {
			service.DevAlarmAdd(&alr)
		}
	}) // 协程池
	defer p.Release()
	defer close(o.alarm)
	for {
		select {
		case v := <-o.alarm:
			if v == nil {
				p.Invoke(dataArr)
				return
			}
			dataArr = append(dataArr, *v)
		case <-ticker.C:
			size := len(dataArr)
			if size < 1 {
				continue
			}
			data := make([]model.DevAlarmDetails, size)
			copy(data, dataArr)
			if err := p.Invoke(data); err != nil {
				continue
			}
			dataArr = dataArr[:0]
		}
	}
}
