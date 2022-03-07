package device

import (
	"xstation/entity/mnger"
	"xstation/model"
	"xstation/service"
	"xstation/util"

	"github.com/wlgd/xproto"
)

// 转换
func devOnlineModel(a *xproto.Access) *model.DevOnline {
	v := &model.DevOnline{
		Guid:          a.Session,
		DeviceNo:      a.DeviceNo,
		RemoteAddress: a.RemoteAddress,
		NetType:       int(a.NetType),
		Type:          int(a.LinkType),
		UpTraffic:     a.UpTraffic,
		DownTraffic:   a.DownTraffic,
		Version:       a.Version,
		DevType:       a.DevType,
	}
	if a.Online {
		v.OnTime = a.DeviceTime
	} else {
		v.OffTime = a.DeviceTime
	}
	return v
}

func devStatusModel(s *xproto.Status) *model.DevStatus {
	m := mnger.Device.Get(s.DeviceNo)
	if m == nil {
		return nil
	}
	o := &model.DevStatus{}
	o.Id = util.PrimaryKey()
	o.DeviceId = m.Id
	o.DeviceNo = s.DeviceNo
	o.DTU = s.DTU
	o.Flag = s.Flag
	if s.Location.Speed < 1 {
		s.Location.Speed = 0
	}
	o.Acc = s.Acc
	o.Location = model.JLocation(s.Location)
	o.Tempers = model.JFloats(s.Tempers)
	o.Humiditys = model.JFloats(s.Humiditys)
	o.Mileage = model.JMileage(s.Mileage)
	o.Oils = model.JOil(s.Oils)
	o.Module = model.JModule(s.Module)
	o.Gsensor = model.JGsensor(s.Gsensor)
	o.Mobile = model.JMobile(s.Mobile)
	o.Disks = model.JDisks(s.Disks)
	o.People = model.JPeople(s.People)
	o.Obds = model.JObds(s.Obds)
	o.Vols = model.JFloats(s.Vol)
	return o
}

func devAlarmModel(a *xproto.Alarm) *model.DevAlarm {
	flag := a.Status.Flag
	a.Status.Flag = 2
	if a.EndTime != "" {
		a.Status.Flag = 3
	}
	o := &model.DevAlarm{
		Guid:      a.UUID,
		Flag:      flag,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
	}
	o.DTU = a.DTU
	o.DeviceNo = a.DeviceNo
	o.AlarmType = a.Type
	o.Data = util.ToJString(a.Data)
	return o
}

func updateDevOnline(a *xproto.Access) error {
	m := mnger.Device.Get(a.DeviceNo)
	o := devOnlineModel(a)
	service.DeviceUpdate(m, a.Online, o.Version, o.DevType)
	service.DevOnlineUpdate(o)
	return nil
}
