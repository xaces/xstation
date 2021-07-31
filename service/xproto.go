package service

import (
	"unsafe"
	"xstation/model"

	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xproto"
)

// DbUpdateOnline 更新链路信息
func DbUpdateOnline(reg *xproto.LinkAccess) error {
	ofline := &model.OnLine{
		Guid:          reg.Session,
		DeviceNo:      reg.DeviceNo,
		RemoteAddress: reg.RemoteAddress,
		NetType:       int(reg.NetType),
		Type:          int(reg.LinkType),
		UpTraffic:     reg.UpTraffic,
		DownTraffic:   reg.DownTraffic,
	}
	if reg.OnLine {
		ofline.OnTime = reg.DeviceTime
		return orm.DbCreate(ofline)
	}
	ofline.OffTime = reg.DeviceTime
	return orm.DbUpdateModelBy(ofline, "guid = ?", ofline.Guid)
}

type StatusTask struct {
	TableIdx int
	Data     []model.Status
	Size     int
}

// DbCreate 批量添加
func DbStatusTaskFunc(obj interface{}) {
	task := obj.(*StatusTask)
	// 映射
	var data interface{}
	ptr := unsafe.Pointer(&task.Data)
	switch task.TableIdx {
	case 1:
		data = (*[]model.Status1)(ptr)
	case 2:
		data = (*[]model.Status2)(ptr)
	case 3:
		data = (*[]model.Status3)(ptr)
	case 4:
		data = (*[]model.Status4)(ptr)
	default:
		data = &task.Data
	}
	orm.DbCreate(data)
}

// StatusModel 获取model模型
func StatusModel(devId uint64) (string, interface{}) {
	tabIdx := int(devId) % StatusTableNum
	switch tabIdx {
	case 1:
		m := &model.Status1{}
		return m.TableName(), m
	case 2:
		m := &model.Status2{}
		return m.TableName(), m
	case 3:
		m := &model.Status3{}
		return m.TableName(), m
	case 4:
		m := &model.Status4{}
		return m.TableName(), m
	default:
	}
	m := &model.Status{}
	return m.TableName(), m
}
