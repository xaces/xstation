package mnger

import (
	"sync"
	"time"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// lServe 服务
type LServe struct {
	ServeId     string    `json:"serveId"` // 服务id
	Token       string    `json:"token"`   //
	UpdatedTime time.Time `json:"updatedTime"`
	IsWorking   bool      `json:"isWorking"`
	Address     string    `json:"address"`
	model.ServeOpt
}

func (s *LServe) Update(address string) error {
	s.Address = address
	s.UpdatedTime = time.Now()
	return nil
}

type serve struct {
	lock      sync.RWMutex
	lServeMap map[string]*LServe
}

var (
	Serve = &serve{lServeMap: make(map[string]*LServe)}
)

// NewLServe 新建服务
func (o *serve) newLServe(serveId string, opt model.ServeOpt) {
	o.lServeMap[serveId] = &LServe{
		ServeId:   serveId,
		ServeOpt:  opt,
		IsWorking: false,
	}
}

// LoadOfDb load sub serve
func (o *serve) LoadOfDb() {
	var serves []model.Serve
	orm.DbFind(&serves)
	for _, v := range serves {
		o.newLServe(v.Guid, v.ServeOpt)
	}
}

// GetByType 根据类型
func (o *serve) GetByType(ctype int) *LServe {
	o.lock.RLock()
	defer o.lock.RUnlock()
	for _, v := range o.lServeMap {
		if v.ServeOpt.Role == ctype && v.IsWorking {
			return v
		}
	}
	return nil
}

// LoadAllLServe 获取实时状态
func (o *serve) GetAll() (lse []LServe) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	for _, v := range o.lServeMap {
		lse = append(lse, *v)
	}
	return
}

// Get 获取
func (o *serve) Get(serveId string) *LServe {
	o.lock.RLock()
	defer o.lock.RUnlock()
	v, ok := o.lServeMap[serveId]
	if !ok {
		return nil
	}
	return v
}

// Add 添加
func (o *serve) Add(srv *model.Serve) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.newLServe(srv.Guid, srv.ServeOpt)
}

// UpdateStatus 更新状态
func (o *serve) UpdateStatus(serveIds []string, status int) error {
	o.lock.Lock()
	defer o.lock.Unlock()
	for _, s := range serveIds {
		v, ok := o.lServeMap[s]
		if !ok || status == v.Status {
			continue
		}
		v.Status = status
		orm.DbUpdateColsBy(&model.Serve{}, orm.H{"status": status}, "guid like ?", s)
	}
	return nil
}

// UpdateStatus 更新状态
func (o *serve) Delete(serveIds []string) {
	o.lock.Lock()
	defer o.lock.Unlock()
	for _, s := range serveIds {
		delete(o.lServeMap, s)
	}
}
