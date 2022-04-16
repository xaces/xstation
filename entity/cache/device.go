package cache

import (
	"sync"
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
	return model.DevStatusVal(m.DeviceId)
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
