package serve

import (
	"sync"
	"xstation/models"
	"xstation/pkg/orm"
)

type devManager struct {
	lDevMap map[string]models.XDevice
	lock    sync.Mutex
}

var (
	DefaultDevsManager = &devManager{lDevMap: make(map[string]models.XDevice)}
)

func (o *devManager) LoadOfDb() {
	var devs []models.XDevice
	if err := orm.DbFind(&devs); err != nil {
		return
	}
	for _, dev := range devs {
		o.lDevMap[dev.VehiNo] = dev
	}
}

// Get 获取
func (o *devManager) Get(vehiNo string) *models.XDevice {
	o.lock.Lock()
	defer o.lock.Unlock()
	if v, ok := o.lDevMap[vehiNo]; ok {
		return &v
	}
	return nil
}

// Add 添加
func (o *devManager) Add(dev *models.XDevice) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.lDevMap[dev.VehiNo] = *dev
}

// Delete 删除
func (o *devManager) Delete(vehiNo string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	delete(o.lDevMap, vehiNo)
}
