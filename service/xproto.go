package service

import (
	"unsafe"
	"xstation/internal"
	"xstation/models"
	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xproto"
)

type XData struct {
}

func NewXData() *XData {
	return new(XData)
}

// DbUpdateAccess 更新链路信息
func (x *XData) DbUpdateAccess(reg *xproto.LinkAccess) error {
	ofline := &models.XLink{
		Guid:          reg.Session,
		DeviceNo:      reg.DeviceNo,
		RemoteAddress: reg.RemoteAddress,
		NetType:       int(reg.NetType),
		Type:          int(reg.LinkType),
		UpTraffic:     reg.UpTraffic,
		DownTraffic:   reg.DownTraffic,
		Version:       reg.Version,
	}
	if reg.OnLine {
		ofline.OnTime = reg.DeviceTime
		return orm.DbCreate(ofline)
	}
	ofline.OffTime = reg.DeviceTime
	return orm.DbUpdateModelBy(ofline, "guid = ?", ofline.Guid)
}

// ToAlarmModel 转化成Model数据格式
func (x *XData) DbCreateAlarm(stId int64, xalr *xproto.Alarm) error {
	alr := &models.XAlarm{
		DeviceNo:  xalr.DeviceNo,
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

type Task struct {
	TableIdx int
	Status   []models.XStatus
	Size     int
}

// DbCreate 批量添加
func DbTaskFunc(obj interface{}) {
	task := obj.(*Task)
	// 映射
	ptr := unsafe.Pointer(&task.Status)
	var data interface{}
	switch task.TableIdx {
	case 1:
		data = (*[]models.XStatus1)(ptr)
	case 2:
		data = (*[]models.XStatus2)(ptr)
	case 3:
		data = (*[]models.XStatus3)(ptr)
	case 4:
		data = (*[]models.XStatus4)(ptr)
	default:
		data = &task.Status
	}
	orm.DbCreate(data)
}

// GetXStatusModel 获取model模型
func GetXStatusModel(devId uint64) interface{} {
	tabIdx := int(devId) % models.KXStatusTabNumber
	switch tabIdx {
	case 1:
		return &models.XStatus1{}
	case 2:
		return &models.XStatus2{}
	case 3:
		return &models.XStatus3{}
	case 4:
		return &models.XStatus4{}
	default:
	}
	return &models.XStatus{}
}
