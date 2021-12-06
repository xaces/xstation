package mnger

import (
	"sync"
	"xstation/model"

	"github.com/wlgd/xproto"
)

type DevMapper struct {
	lDevMap map[string]*model.Device
	lock    sync.RWMutex
}

var (
	Devs = &DevMapper{lDevMap: make(map[string]*model.Device)}
)

func (o *DevMapper) Set(devs []model.Device) {
	for i := 0; i < len(devs); i++ {
		o.lDevMap[devs[i].No] = &devs[i]
	}
}

// Get 获取
func (o *DevMapper) Get(deviceNo string) *model.Device {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if v, ok := o.lDevMap[deviceNo]; ok {
		return v
	}
	return nil
}

// Add 添加
func (o *DevMapper) Add(dev *model.Device) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.lDevMap[dev.No] = dev
}

// Delete 删除
func (o *DevMapper) Delete(deviceNo string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.lDevMap, deviceNo)
	xproto.SyncStop(deviceNo)
}

func (o *DevMapper) Model(deviceNo string) (string, interface{}) {
	dev := o.Get(deviceNo)
	if dev == nil {
		return "", nil
	}
	return model.Status(dev.Id)
}
