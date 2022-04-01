package service

import (
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// AlarmPage 分页
type AlarmPage struct {
	orm.DbPage
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
	where := s.DbWhere()
	where.String("start_time >= ?", s.StartTime)
	where.String("start_time <= ?", s.EndTime)
	where.String("device_no like ?", s.DeviceNo)
	where.Orders = append(where.Orders, "dtu desc")
	return where
}

type AlarmDetailsPage struct {
	orm.DbPage
	DeviceNo  string `form:"deviceNo"` //
	Guid      string `form:"guid"`     //
	LinkType  int    `form:"linkType"` //
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	AlarmType int    `form:"alarmType"`
}

func (s *AlarmDetailsPage) Where() *orm.DbWhere {
	where := s.DbWhere()
	where.String("dtu >= ?", s.StartTime)
	where.String("dtu <= ?", s.EndTime)
	where.String("guid = ?", s.Guid)
	where.String("device_no like ?", s.DeviceNo)
	where.Int("type = ?", s.AlarmType)
	where.Int("link_type", s.LinkType)
	where.Orders = append(where.Orders, "dtu desc")
	return where
}

func DevAlarmAdd(alr *model.DevAlarm) error {
	if alr.Status == 0 {
		return orm.DbCreate(alr)
	}
	return orm.DbUpdateSelectWhere(alr, []string{"dtu, end_data, end_time, end_status"}, "guid = ?", alr.Guid)
}
