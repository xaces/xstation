package service

import (
	"xstation/model"

	"github.com/xaces/xutils/orm"
)

func DevAlarmAdd(a *model.DevAlarmDetails) error {
	o := &model.DevAlarm{
		GUID:      a.GUID,
		DeviceID:  a.DeviceID,
		DTU:       a.DTU,
		AlarmType: a.AlarmType,
	}
	if a.Status == 0 {
		o.StartTime = a.StartTime
		o.StartData = a.Data
		o.StartStatus = a.DevStatus
		return orm.DbCreate(o)
	}
	o.EndTime = a.EndTime
	o.EndStatus = a.DevStatus
	o.EndData = a.Data
	return orm.DbUpdatesBy(o, []string{"dtu", "end_time", "end_data", "end_status"}, "guid = ?", o.GUID)
}
