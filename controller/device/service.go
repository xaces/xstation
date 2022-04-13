package device

import (
	"xstation/model"
	"xstation/util"

	"github.com/wlgd/xproto"
	"github.com/wlgd/xutils/orm"
)

// 转换
func devOnlineUpdate(a *xproto.Access, s *xproto.Status) error {
	o := &model.DevOnline{
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
		o.OnlineTime = a.DeviceTime
		o.OnlineStatus = devStatusModel(s)
		return orm.DbCreate(o)
	}
	o.OfflineTime = a.DeviceTime
	o.OfflineStatus = devStatusModel(s)
	return orm.DbUpdatesBy(o, []string{"offline_time, offline_status"}, "guid = ?", o.Guid)
}

func devStatusModel(s *xproto.Status) *model.DevStatus {
	if s == nil {
		return nil
	}
	o := &model.DevStatus{}
	o.Id = util.PrimaryKey()
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

func devAlarmModel(a *xproto.Alarm, s *model.DevStatus) *model.DevAlarm {
	o := &model.DevAlarm{
		Guid:      a.UUID,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		DTU:       a.DTU,
		DeviceNo:  a.DeviceNo,
		AlarmType: a.Type,
	}
	if a.DTU > a.StartTime {
		o.EndStatus = s
		o.EndData = util.JString(a.Data)
	} else {
		o.StartData = util.JString(a.Data)
		o.StartStatus = s
	}
	return o
}

func devAlarmDetailsModel(a *xproto.Alarm) *model.DevAlarmDetails {
	o := &model.DevAlarmDetails{}
	o.DeviceNo = a.DeviceNo
	o.DTU = a.DTU
	o.Guid = a.UUID
	o.LinkType = model.AlarmLinkDev
	o.Data = util.JString(a.Data)
	o.DevStatus = devStatusModel(a.Status)
	o.Flag = a.Status.Flag
	o.AlarmType = a.Type
	o.Status = 0
	if a.DTU > a.StartTime {
		o.Status = 2
		if a.EndTime != "" {
			o.Status = 1
		}
	}
	return o
}
