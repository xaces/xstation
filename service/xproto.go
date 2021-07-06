package service

import (
	"unsafe"
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

type StatusTask struct {
	TableIdx int
	Data     []models.XStatus
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
		data = (*[]models.XStatus1)(ptr)
	case 2:
		data = (*[]models.XStatus2)(ptr)
	case 3:
		data = (*[]models.XStatus3)(ptr)
	case 4:
		data = (*[]models.XStatus4)(ptr)
	default:
		data = &task.Data
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
