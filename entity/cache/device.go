package cache

import (
	"sync"
	"unsafe"
	"xstation/model"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils"
)

type Vehicle struct {
	DeviceId      uint   `json:"deviceId"`
	DeviceNo      string `json:"deviceNo"`
	EffectiveTime string `json:"effectiveTime"`
}

type VehicleFtp struct {
	DeviceNo string `json:"deviceNo"`
	Alarms   string `json:"alarms"`
}

type mdevice struct {
	Vehicle
	Online         bool           `json:"online"`
	LastOnlineTime string         `json:"lastOnlineTime"`
	Status         xproto.Status  `json:"status"`
	FtpAlarms      *xutils.BitMap `json:"-"`
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
	gDevlock.RLock()
	defer gDevlock.RUnlock()
	if v, ok := gDevices[deviceNo]; ok {
		return v
	}
	return nil
}

// 新建
func NewDevice(vehi Vehicle) *mdevice {
	gDevlock.Lock()
	defer gDevlock.Unlock()
	v := &mdevice{Vehicle: vehi}
	v.FtpAlarms = xutils.NewBitMapWithBase(0, 1000)
	gDevices[v.DeviceNo] = v
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
