package cache

import (
	"sync"
	"xstation/configs"
	"xstation/model"

	"github.com/wlgd/xproto"
)

type Vehicle struct {
	ID             uint   `json:"deviceId"`
	No             string `json:"deviceNo"`
	EffectiveTime  string `json:"effectiveTime"`
	LastOnlineTime string `json:"lastOnlineTime"`
}

type VehicleFtp struct {
	DeviceNo string `json:"deviceNo"`
	Alarms   string `json:"alarms"`
}

type mdevice struct {
	Vehicle
	Online bool          `json:"online"`
	Status xproto.Status `json:"status"`
}

func (m *mdevice) Update(a *xproto.Access) {
	m.Online = a.Online
	m.LastOnlineTime = a.DeviceTime
}

func (m *mdevice) Model() interface{} {
	if configs.MsgProc > 0 {
		return model.DevStatusVal(m.ID)
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
	gDevices[v.No] = v
	return v
}

// 删除
func DeviceDel(deviceNo string) {
	xproto.SyncStop(deviceNo)
	gDevlock.Lock()
	defer gDevlock.Unlock()
	delete(gDevices, deviceNo)
}

func DeviceList() (devs []mdevice) {
	for _, v := range gDevices {
		devs = append(devs, *v)
	}
	return
}
