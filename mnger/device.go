package mnger

import (
	"sync"
	"xstation/model"
	"xstation/service"

	"github.com/wlgd/xproto"
)

type devMapper struct {
	lDevMap map[string]*model.Device
	lock    sync.RWMutex
}

var (
	Dev = &devMapper{lDevMap: make(map[string]*model.Device)}
)

func (o *devMapper) Set(devs []model.Device) {
	for i := 0; i < len(devs); i++ {
		o.lDevMap[devs[i].No] = &devs[i]
	}
}

// Get 获取
func (o *devMapper) Get(deviceNo string) *model.Device {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if v, ok := o.lDevMap[deviceNo]; ok {
		return v
	}
	return nil
}

// Add 添加
func (o *devMapper) Add(dev *model.Device) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.lDevMap[dev.No] = dev
}

// Delete 删除
func (o *devMapper) Delete(deviceNo string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.lDevMap, deviceNo)
	xproto.SyncStop(deviceNo)
}

func (o *devMapper) Model(deviceNo string) (string, interface{}) {
	dev := o.Get(deviceNo)
	if dev == nil {
		return "", nil
	}
	return service.StatusModel(dev.Id)
}
