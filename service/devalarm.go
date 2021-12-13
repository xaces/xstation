package service

import (
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// AlarmPage 分页
type AlarmPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
}

// // Where 初始化
// func (s *AlarmPage) Where() string {
// 	tbname, _ := mnger.Devs.Model(s.DeviceNo)
// 	sql := "SELECT a.*, s.gps FROM t_devalarm a JOIN %s s" +
// 		" ON a.device_no like '%s' AND a.start_time >= '%s' AND a.start_time <= '%s' AND a.status_id = s.id ORDER BY a.start_time desc"
// 	sqlstr := fmt.Sprintf(sql, tbname, s.DeviceNo, s.StartTime, s.EndTime)
// 	return sqlstr
// }

func (s *AlarmPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.String("start_time >= ?", s.StartTime)
	where.String("start_time <= ?", s.EndTime)
	where.String("device_no like ?", s.DeviceNo)
	where.Orders = append(where.Orders, "start_time desc")
	return &where
}

type AlarmLinkPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	DeviceNo  string `form:"deviceNo"` //
	Guid      string `form:"guid"`     //
	LinkType  int    `form:"linkType"` //
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	AlarmType int    `form:"alarmType"`
}

func (s *AlarmLinkPage) Where() *orm.DbWhere {
	var where orm.DbWhere
	where.String("dtu >= ?", s.StartTime)
	where.String("dtu <= ?", s.EndTime)
	where.String("guid = ?", s.Guid)
	where.String("device_no like ?", s.DeviceNo)
	where.Int("type = ?", s.AlarmType)
	where.Int("link_type", s.LinkType)
	where.Orders = append(where.Orders, "dtu desc")
	return &where
}

// type AlarmData struct {
// 	model.DevAlarm
// 	Gps model.JGps `json:"gps"`
// }

func alarmLinkCreate(alr *model.DevAlarm) error {
	l := &model.DevAlarmLink{}
	l.Guid = alr.Guid
	l.DevAlarmOpt = alr.DevAlarmOpt
	l.LinkType = model.AlarmLinkDev
	l.DevStatus = alr.DevStatus
	return orm.DbCreate(&l)
}

func AlarmDbAdd(alr *model.DevAlarm) error {
	alarmLinkCreate(alr)
	upfields := []string{"status", "dtu"}
	if alr.DTU != alr.StartTime {
		if alr.EndTime != "" {
			alr.Status = 1
			upfields = append(upfields, "end_time")
		} else {
			upfields = append(upfields, "data")
			alr.Status = 2
		}
		if orm.DbUpdateSelectWhere(alr, upfields, "guid = ?", alr.Guid) == nil {
			return nil
		}
	}
	return orm.DbCreate(alr)
}
