package device

import "github.com/wlgd/xutils/orm"

// 搜索条件
type Where struct {
	orm.DbPage
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
	DeviceNo  string `form:"deviceNo"` //

	Guid      string `form:"guid"`     //
	LinkType  int    `form:"linkType"` //
	AlarmType int    `form:"alarmType"`

	Flag int    `form:"flag"` // 0-实时 1-补传
	Desc string `form:"desc"` //
}

// Where 初始化
func (s *Where) Online() *orm.DbWhere {
	where := s.DbWhere()
	where.String("device_no = ?", s.DeviceNo)
	where.String("online_time >= ?", s.StartTime)
	where.String("online_time <= ?", s.EndTime)
	where.Orders = append(where.Orders, "online_time desc")
	return where
}

// // Where 初始化
// func (s *AlarmPage) Where() string {
// 	tbname, _ := mnger.Devs.Model(s.DeviceNo)
// 	sql := "SELECT a.*, s.gps FROM t_devalarm a JOIN %s s" +
// 		" ON a.device_no like '%s' AND a.start_time >= '%s' AND a.start_time <= '%s' AND a.status_id = s.id ORDER BY a.start_time desc"
// 	sqlstr := fmt.Sprintf(sql, tbname, s.DeviceNo, s.StartTime, s.EndTime)
// 	return sqlstr
// }

func (s *Where) Alarm() *orm.DbWhere {
	where := s.DbWhere()
	where.String("start_time >= ?", s.StartTime)
	where.String("start_time <= ?", s.EndTime)
	where.String("device_no = ?", s.DeviceNo)
	where.Orders = append(where.Orders, "start_time desc")
	return where
}

func (s *Where) AlarmDetailsPage() *orm.DbWhere {
	where := s.DbWhere()
	where.String("dtu >= ?", s.StartTime)
	where.String("dtu <= ?", s.EndTime)
	where.String("guid = ?", s.Guid)
	where.String("device_no = ?", s.DeviceNo)
	where.Int("type = ?", s.AlarmType)
	where.Int("link_type", s.LinkType)
	where.Orders = append(where.Orders, "dtu desc")
	return where
}

// Where 初始化
func (s *Where) Status() *orm.DbWhere {
	where := s.DbWhere()
	where.String("device_no = ?", s.DeviceNo)
	where.String("dtu >= ?", s.StartTime)
	where.String("dtu <= ?", s.EndTime)
	where.Int("flag = ?", s.Flag)
	where.Orders = append(where.Orders, "dtu desc")
	return where
}
