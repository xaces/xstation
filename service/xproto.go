package service

import (
	"xstation/internal"
	"xstation/models"
	"xstation/pkg/orm"

	"github.com/wlgd/xproto"
)

type XData struct {
}

func NewXData() *XData {
	return new(XData)
}

// DbUpdateAccess 更新链路信息
func (x *XData) DbUpdateAccess(reg *xproto.LinkAccess) error {
	ofline := &models.XOFLine{
		Guid:          reg.Session,
		DeviceId:      reg.DeviceId,
		RemoteAddress: reg.RemoteAddress,
		AccessType:    int(reg.AccessNet),
		Type:          int(reg.LinkType),
		UpFlow:        reg.UpFlow,
		DownFlow:      reg.DownFlow,
		Version:       reg.Version,
	}
	if reg.OnLine {
		ofline.OnTime = reg.DeviceTime
		return orm.DbCreate(ofline)
	}
	ofline.OffTime = reg.DeviceTime
	return orm.DbUpdateModelBy(ofline, "guid = ?", ofline.Guid)
}

// ToStatusModel 转化成Model数据格式
func (x *XData) ToStatusModel(st *xproto.Status) (o models.XStatus) {
	o.Id = PrimaryKey()
	o.DeviceId = st.DeviceId
	o.DTU = st.DTU
	o.Status = st.Status
	if st.Location.Speed < 1 {
		st.Location.Speed = 0
	}
	o.Gps = internal.ToJString(st.Location)
	o.Tempers = internal.ToJString(st.Tempers)
	o.Humiditys = internal.ToJString(st.Humiditys)
	o.Mileage = internal.ToJString(st.Mileage)
	o.Oils = internal.ToJString(st.Oils)
	o.Module = internal.ToJString(st.Module)
	o.Gsensor = internal.ToJString(st.Gsensor)
	o.Mobile = internal.ToJString(st.Mobile)
	o.Disks = internal.ToJString(st.Disks)
	o.People = internal.ToJString(st.People)
	return
}

// ToAlarmModel 转化成Model数据格式
func (x *XData) DbCreateAlarm(stId int64, xalr *xproto.Alarm) error {
	alr := &models.XAlarm{
		DeviceId:  xalr.DeviceId,
		UUID:      xalr.UUID,
		StatusId:  stId,
		Status:    xalr.Status.Status,
		Type:      xalr.Type,
		StartTime: xalr.StartTime,
		EndTime:   xalr.EndTime,
		Data:      internal.ToJString(xalr.Data),
	}
	return orm.DbCreate(alr)
}

// DbCreate 批量添加
func (x *XData) DbCreateStatus(stArray []models.XStatus) error {
	size := len(stArray)
	if size <= 0 {
		return nil
	}
	if size <= 10 {
		return orm.DbCreate(&stArray)
	}
	status := make([]models.XStatus, size)
	copy(status, stArray)
	go func() {
		orm.DbCreate(&status)
	}()
	return nil
}
