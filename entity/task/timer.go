package task

import (
	"time"

	"github.com/go-co-op/gocron"
)

type timer struct {
	gScheduler *gocron.Scheduler
}

var Timer = timer{
	gScheduler: gocron.NewScheduler(time.UTC),
}

func (o *timer) AddDbPartFunc(handle func()) {
	o.gScheduler.Every("1h").Do(handle)
}

func (o *timer) Run() {
	o.gScheduler.StartAsync()
}
