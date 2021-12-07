package service

import (
	"unsafe"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// StatusPage 分页
type StatusPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
	Flag      int    `form:"flag"`     // 0-实时 1-补传
	Desc      string `form:"desc"`     //
}

// Where 初始化
func (s *StatusPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.String("device_no like ?", s.DeviceNo)
	where.String("dtu >= ?", s.StartTime+" 00:00:00")
	where.String("dtu <= ?", s.EndTime+" 23:59:59")
	where.Int("flag = ?", s.Flag)
	where.Orders = append(where.Orders, "dtu desc")
	return &where
}

// StatusGet 获取
type StatusGet struct {
	DeviceNo string `form:"deviceNo"` //
	StatusId uint64 `form:"statusId"` //
}

type StatusTask struct {
	TableIdx int
	Data     []model.DevStatus
	Size     int
}

// StatusCreates 批量添加
func StatusCreates(obj interface{}) {
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
