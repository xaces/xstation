package mnger

import (
	"sync"
	"unsafe"
	"xstation/model"

	"github.com/wlgd/xproto"
)

type deviceMapper struct {
	lDevMap map[string]*model.Device
	lock    sync.RWMutex
}

var (
	Device = &deviceMapper{lDevMap: make(map[string]*model.Device)}
)

func (o *deviceMapper) Set(devs []model.Device) {
	for i := 0; i < len(devs); i++ {
		o.lDevMap[devs[i].DeviceNo] = &devs[i]
	}
}

// Get 获取
func (o *deviceMapper) Get(deviceNo string) *model.Device {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if v, ok := o.lDevMap[deviceNo]; ok {
		return v
	}
	return nil
}

// Add 添加
func (o *deviceMapper) Add(dev *model.Device) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.lDevMap[dev.DeviceNo] = dev
}

// Delete 删除
func (o *deviceMapper) Delete(deviceNo string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.lDevMap, deviceNo)
	xproto.SyncStop(deviceNo)
}

func (o *deviceMapper) StatusModel(deviceNo string) interface{} {
	dev := o.Get(deviceNo)
	if dev == nil {
		return nil
	}
	tabidx := int(dev.Id) % model.DevStatusNum
	switch tabidx {
	case 1:
		return &model.DevStatus1{}
	case 2:
		return &model.DevStatus1{}
	case 3:
		return model.DevStatus3{}
	case 4:
		return &model.DevStatus4{}
	}
	return &model.DevStatus{}
}

func (o *deviceMapper) StatusValue(tabidx int, v []model.DevStatus) interface{} {
	ptr := unsafe.Pointer(&v)
	switch tabidx {
	case 1:
		return (*[]model.DevStatus1)(ptr)
	case 2:
		return (*[]model.DevStatus2)(ptr)
	case 3:
		return (*[]model.DevStatus3)(ptr)
	case 4:
		return (*[]model.DevStatus4)(ptr)
	}
	return v
}
