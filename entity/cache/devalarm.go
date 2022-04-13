package cache

import (
	"time"
	"xstation/model"

	"github.com/FishGoddess/cachego"
	"github.com/wlgd/xutils/orm"
)

type alarmMapper struct {
	Cache *cachego.Cache
}

var (
	gAlarm = cachego.NewCache(cachego.WithAutoGC(60 * time.Minute))
)

// DevAlarm 添加
func NewDevAlarm(a *model.DevAlarmDetails) {
	gAlarm.Set(a.Guid, *a)
}

func DevAlarm(ss string) *model.DevAlarmDetails {
	var alr model.DevAlarmDetails
	if data, ok := gAlarm.Get(ss); ok {
		alr = data.(model.DevAlarmDetails)
	}
	if err := orm.DbFirstBy(&alr, "guid = ?", ss); err == nil {
		NewDevAlarm(&alr)
		return &alr
	}
	return nil
}
