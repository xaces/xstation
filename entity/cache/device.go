package cache

import (
	"sync"

	"github.com/xaces/xproto"
)

type DeviceInfo struct {
	ID             uint   `json:"deviceId"`
	No             string `json:"deviceNo"`
	EffectiveTime  string `json:"effectiveTime"`
	LastOnlineTime string `json:"lastOnlineTime"`
}

type DeviceFtp struct {
	No     string `json:"deviceNo"`
	Alarms string `json:"alarms"`
}

type Device struct {
	DeviceInfo
	Online bool          `json:"online"`
	Status xproto.Status `json:"status"`
}

func (m *Device) Update(a *xproto.Access) {
	m.Online = a.Online
	m.LastOnlineTime = a.DeviceTime
}

var (
	gDevices = make(map[string]*Device)
	gDevlock sync.RWMutex
)

// 获取
func GetDevice(deviceNo string) *Device {
	gDevlock.RLock()
	defer gDevlock.RUnlock()
	if v, ok := gDevices[deviceNo]; ok {
		return v
	}
	return nil
}

// 新建
func NewDevice(info DeviceInfo) *Device {
	gDevlock.Lock()
	defer gDevlock.Unlock()
	v := &Device{DeviceInfo: info}
	gDevices[v.No] = v
	return v
}

// 删除
func DelDevice(deviceNo string) {
	xproto.SyncStop(deviceNo)
	gDevlock.Lock()
	defer gDevlock.Unlock()
	delete(gDevices, deviceNo)
}

func ListDevice() (devs []Device) {
	for _, v := range gDevices {
		devs = append(devs, *v)
	}
	return
}
