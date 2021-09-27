package device

import (
	"fmt"
	"xstation/mnger"
	"xstation/model"
)

// alarmPage 分页
type alarmPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
}

type alarmData struct {
	model.DevAlarm
	Gps model.JGps `json:"gps"`
}

// Where 初始化
func (s *alarmPage) Where() string {
	tbname, _ := mnger.Dev.Model(s.DeviceNo)
	sql := "SELECT a.*, s.gps FROM t_devalarm a JOIN %s s" +
		" ON a.device_no like '%s' AND a.start_time >= '%s' AND a.start_time <= '%s' AND a.status_id = s.id ORDER BY a.start_time desc"
	sqlstr := fmt.Sprintf(sql, tbname, s.DeviceNo, s.StartTime, s.EndTime)
	return sqlstr
}
