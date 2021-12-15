package mnger

import (
	"time"
	"xstation/model"

	"github.com/FishGoddess/cachego"
	"github.com/wlgd/xutils/orm"
)

type AlarmMapper struct {
	Cache *cachego.Cache
}

var (
	Alarm = &AlarmMapper{Cache: cachego.NewCache(cachego.WithAutoGC(60 * time.Minute))}
)

// Add 添加
func (o *AlarmMapper) Add(p model.DevAlarm) {
	o.Cache.Set(p.Guid, p)
}

func (o *AlarmMapper) Get(ss string) *model.DevAlarm {
	var alr model.DevAlarm
	if data, ok := o.Cache.Get(ss); ok {
		alr = data.(model.DevAlarm)
	}
	if err := orm.DbFirstBy(&alr, "guid = ?", ss); err == nil {
		o.Add(alr)
		return &alr
	}
	return nil
}
