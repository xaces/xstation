package mnger

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
	Alarm = &alarmMapper{Cache: cachego.NewCache(cachego.WithAutoGC(60 * time.Minute))}
)

// Add 添加
func (o *alarmMapper) Add(a *model.DevAlarm) {
	o.Cache.Set(a.Guid, *a)
}

func (o *alarmMapper) Get(ss string) *model.DevAlarm {
	var alr model.DevAlarm
	if data, ok := o.Cache.Get(ss); ok {
		alr = data.(model.DevAlarm)
	}
	if err := orm.DbFirstBy(&alr, "guid = ?", ss); err == nil {
		o.Add(&alr)
		return &alr
	}
	return nil
}
