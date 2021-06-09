package serve

import (
	"sync"
	"xstation/models"
	"xstation/pkg/orm"
)

var (
	lDevMap = make(map[string]models.XDevice)
	lock    sync.Mutex
)

func loadAllDevices() {
	var devs []models.XDevice
	if err := orm.DbFind(&devs); err != nil {
		return
	}
	for _, dev := range devs {
		lDevMap[dev.VehiNo] = dev
	}
}

func getDeivce(vehiNo string) *models.XDevice {
	lock.Lock()
	defer lock.Unlock()
	if v, ok := lDevMap[vehiNo]; ok {
		return &v
	}
	return nil
}
