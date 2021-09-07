package mnger

import (
	"sync"
	"xstation/model"
	"xstation/service"

	"github.com/wlgd/xproto"
)

type devManager struct {
	lDevMap map[string]*model.Device
	lock    sync.RWMutex
}

var (
	Dev = &devManager{lDevMap: make(map[string]*model.Device)}
)

func (o *devManager) Set(devs []model.Device) {
	for _, dev := range devs {
		o.lDevMap[dev.No] = &dev
	}
}

// Get 获取
func (o *devManager) Get(deviceNo string) *model.Device {
	o.lock.RLock()
	defer o.lock.RUnlock()
	if v, ok := o.lDevMap[deviceNo]; ok {
		return v
	}
	return nil
}

// Add 添加
func (o *devManager) Add(dev *model.Device) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.lDevMap[dev.No] = dev
}

// Delete 删除
func (o *devManager) Delete(deviceNo string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.lDevMap, deviceNo)
	xproto.SyncStop(deviceNo)
}

func (o *devManager) GetModel(deviceNo string) (string, interface{}) {
	dev := o.Get(deviceNo)
	if dev == nil {
		return "", nil
	}
	return service.StatusModel(dev.Id)
}
