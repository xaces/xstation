package task

import (
	"time"

	"github.com/go-co-op/gocron"
)

type timer struct {
	scheduler *gocron.Scheduler
}

var Timer = timer{
	scheduler: gocron.NewScheduler(time.UTC),
}

func (o *timer) AddDbPartFunc(handle func()) {
	o.scheduler.Every("1h").Do(handle)
}

func (o *timer) Run() {
	o.scheduler.StartAsync()
}
