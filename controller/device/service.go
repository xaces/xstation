package device

import (
	"xstation/model"
	"xstation/util"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

// 转换
func devOnlineUpdate(a *xproto.Access, s *xproto.Status, deviceId uint) error {
	o := &model.DevOnline{
		GUID:          a.Session,
		DeviceID:      deviceId,
		RemoteAddress: a.RemoteAddress,
		NetType:       int(a.NetType),
		Type:          int(a.LinkType),
		UpTraffic:     a.UpTraffic,
		DownTraffic:   a.DownTraffic,
		Version:       a.Version,
		DevType:       a.DevType,
	}
	if a.Online {
		o.OnlineTime = a.DeviceTime
		o.OnlineStatus = devStatusModel(s)
		return orm.DbCreate(o)
	}
	o.OfflineTime = a.DeviceTime
	o.OfflineStatus = devStatusModel(s)
	return orm.DbUpdatesBy(o, []string{"offline_time, offline_status"}, "guid = ?", o.GUID)
}

func deviceUpdate(deviceId uint, a *xproto.Access) error {
	o := &model.Device{
		Type:           a.DevType,
		Version:        a.Version,
		Online:         a.Online,
		LastOnlineTime: a.DeviceTime,
	}
	o.ID = deviceId
	return orm.DbUpdateModel(o)
}

func devStatusModel(s *xproto.Status) *model.DevStatus {
	o := &model.DevStatus{}
	if s == nil {
		return o
	}
	o.DTU = s.DTU
	o.Flag = s.Flag
	if s.Location.Speed < 1 {
		s.Location.Speed = 0
	}
	o.Acc = s.Acc
	o.Location = model.JLocation(s.Location)
	o.Mileage = model.JMileage(s.Mileage)
	o.Oils = model.JOil(s.Oils)
	o.P1 = model.JParam1{
		Obds:      s.Obds,
		Tempers:   s.Tempers,
		Humiditys: s.Humiditys,
		Module:    s.Module,
		Disks:     s.Disks,
	}
	o.P2 = model.JParam2{
		Gsensor: s.Gsensor,
		People:  s.People,
		Vols:    s.Vol,
	}
	return o
}

func devAlarmDetailsModel(a *xproto.Alarm) *model.DevAlarmDetails {
	o := &model.DevAlarmDetails{
		DTU:       a.DTU,
		AlarmType: a.Type,
		GUID:      a.UUID,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		Data:      util.JString(a.Data),
		DevStatus: devStatusModel(a.Status),
	}
	o.Flag = a.Status.Flag
	o.DevStatus.Flag = 2
	o.Status = 0
	if a.DTU > a.StartTime {
		o.Status = 2
		if a.EndTime != "" {
			o.Status = 1
		}
		o.DevStatus.Flag = 3
	}
	return o
}
