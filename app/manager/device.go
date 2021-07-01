package manager

import (
	"sync"
	"xstation/models"
	"xstation/pkg/orm"

	"github.com/wlgd/xproto"
)

type devManager struct {
	lDevMap map[string]models.XDevice
	lock    sync.Mutex
}

var (
	Dev = &devManager{lDevMap: make(map[string]models.XDevice)}
)

func (o *devManager) LoadOfDb() {
	var devs []models.XDevice
	if err := orm.DbFind(&devs); err != nil {
		return
	}
	for _, dev := range devs {
		o.lDevMap[dev.DeviceNo] = dev
	}
}

// Get 获取
func (o *devManager) Get(deviceNo string) *models.XDevice {
	o.lock.Lock()
	defer o.lock.Unlock()
	if v, ok := o.lDevMap[deviceNo]; ok {
		return &v
	}
	return nil
}

// Add 添加
func (o *devManager) Add(dev *models.XDevice) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.lDevMap[dev.DeviceNo] = *dev
}

// Delete 删除
func (o *devManager) Delete(deviceNo string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.lDevMap, deviceNo)
	xproto.SyncStopConnection(deviceNo)
}
