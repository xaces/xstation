package service

import (
	"unsafe"
	"xstation/model"

	"github.com/wlgd/xutils/orm"

	"github.com/wlgd/xproto"
)

// DbUpdateOnline 更新链路信息
func DbUpdateOnline(reg *xproto.Access) error {
	ofline := &model.DevOnline{
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
	Data     []model.DevStatus
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
		data = (*[]model.DevStatus1)(ptr)
	case 2:
		data = (*[]model.DevStatus2)(ptr)
	case 3:
		data = (*[]model.DevStatus3)(ptr)
	case 4:
		data = (*[]model.DevStatus4)(ptr)
	default:
		data = &task.Data
	}
	orm.DbCreate(data)
}
