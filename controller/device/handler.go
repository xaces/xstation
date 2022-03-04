package device

import (
	"time"
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"

	"github.com/panjf2000/ants/v2"
	"github.com/wlgd/xutils/orm"
)

type handler struct {
	Status chan model.DevStatus
	Alarm  chan model.DevAlarm
}

var (
	Handler *handler = &handler{
		Status: make(chan model.DevStatus, 1),
		Alarm:  make(chan model.DevAlarm, 1),
	}
	nats    *natsHandler
	msgProc string = ""
)

func (h *handler) Run(msgproc string) {
	msgProc = msgproc
	if msgproc == "nats" {
		nats = &natsHandler{}
		nats.Run()
	}
	go h.StatusDispatch()
	go h.AlarmDispatch()
}

func (h *handler) AddStatus(v model.DevStatus) {
	h.Status <- v
}

func (h *handler) AddAlarm(v model.DevAlarm) {
	h.Alarm <- v
}

type statusObj struct {
	TableIdx int
	Data     []model.DevStatus
}

// DbStatusHandler 批量处理数据
func (h *handler) StatusDispatch() {
	stArray := make([][]model.DevStatus, model.DevStatusNum)
	statusFunc := func(v interface{}) {
		obj := v.(*statusObj)
		data := mnger.Device.StatusValue(obj.TableIdx, obj.Data)
		orm.DbCreate(data)
	}
	p, _ := ants.NewPoolWithFunc(5, statusFunc) // 协程池
	ticker := time.NewTicker(time.Second * 2)
	defer p.Release()
	defer close(h.Status)
	for {
		select {
		case v := <-h.Status:
			tabIdx := v.DeviceId % model.DevStatusNum
			stArray[tabIdx] = append(stArray[tabIdx], v)
		case <-ticker.C:
			for i := 0; i < model.DevStatusNum; i++ {
				size := len(stArray[i])
				if size < 1 {
					continue
				}
				o := &statusObj{}
				o.TableIdx = i
				o.Data = make([]model.DevStatus, size)
				copy(o.Data, stArray[i])
				if err := p.Invoke(o); err != nil {
					continue
				}
				stArray[i] = stArray[i][:0]
			}
		}
	}
}

// AlarmDispatch 批量处理数据
func (h *handler) AlarmDispatch() {
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
	defer close(h.Alarm)
	for {
		select {
		case v := <-h.Alarm:
			// 推送给第三放
			stArray = append(stArray, v)
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
