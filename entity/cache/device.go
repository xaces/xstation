package cache

import (
	"sync"
	"unsafe"
	"xstation/model"

	"github.com/wlgd/xproto"
)

type mdevice struct {
	DeviceId       uint          `json:"deviceId"` // deviceId为0，表示不合法设备
	DeviceNo       string        `json:"deviceNo"`
	Online         bool          `json:"online"`
	LastOnlineTime string        `json:"lastOnlineTime"`
	Status         xproto.Status `json:"status"`
}

func (m *mdevice) Update(a *xproto.Access) {
	m.Online = a.Online
	m.LastOnlineTime = a.DeviceTime
}

func (m *mdevice) Model() interface{} {
	tabIdex := m.DeviceId % model.DevStatusNum
	switch tabIdex {
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

var (
	gDevices = make(map[string]*mdevice)
	gDevlock sync.RWMutex
)

// 获取
func Device(deviceNo string) *mdevice {
	gDevlock.Lock()
	defer gDevlock.Unlock()
	v, ok := gDevices[deviceNo]
	if !ok {
		v = &mdevice{DeviceId: 0, DeviceNo: deviceNo}
		gDevices[deviceNo] = v
	}
	return v
}

// 新建
func NewDevice(deviceId uint, deviceNo string) *mdevice {
	gDevlock.Lock()
	defer gDevlock.Unlock()
	v := &mdevice{DeviceId: deviceId, DeviceNo: deviceNo}
	gDevices[deviceNo] = v
	return v
}

// 删除
func DeviceDel(deviceNo string) {
	gDevlock.Lock()
	defer gDevlock.Unlock()
	delete(gDevices, deviceNo)
	xproto.SyncStop(deviceNo)
}

func DeviceList() (devs []mdevice) {
	for _, v := range gDevices {
		devs = append(devs, *v)
	}
	return
}

func DevStatus(tabidx int, v []model.DevStatus) interface{} {
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
