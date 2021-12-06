package service

import (
	"fmt"
	"xstation/app/mnger"
	"xstation/model"
)

// AlarmPage 分页
type AlarmPage struct {
	PageNum   int    `form:"pageNum"`  // 当前页码
	PageSize  int    `form:"pageSize"` // 每页数
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //
}

type AlarmData struct {
	model.DevAlarm
	StartGps model.JGps `json:"startGps"`
	EndGps   model.JGps `json:"endGps"`
}

// Where 初始化
func (s *AlarmPage) Where() string {
	tbname, _ := mnger.Devs.Model(s.DeviceNo)
	sql := "SELECT a.*, s.gps FROM t_devalarm a JOIN %s s" +
		" ON a.device_no like '%s' AND a.start_time >= '%s' AND a.start_time <= '%s' AND a.status_id = s.id ORDER BY a.start_time desc"
	sqlstr := fmt.Sprintf(sql, tbname, s.DeviceNo, s.StartTime, s.EndTime)
	return sqlstr
}
