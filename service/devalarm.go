package service

import "github.com/wlgd/xutils/orm"

// AlarmPage 分页
type AlarmPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
}

// type AlarmData struct {
// 	model.DevAlarm
// 	Gps model.JGps `json:"gps"`
// }

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
