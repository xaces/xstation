package mnger

import (
	"sync"
	"time"
	"xstation/internal"
	"xstation/models"

	"github.com/wlgd/xutils/orm"
	"github.com/wlgd/xutils/rpc"
)

// lServe 服务
type LServe struct {
	ServeId     string    `json:"serveId"` // 服务id
	Token       string    `json:"token"`   //
	UpdatedTime time.Time `json:"updatedTime"`
	IsWorking   bool      `json:"isWorking"`
	Address     string    `json:"address"`
	models.XServerOpt
}

type serve struct {
	lock      sync.RWMutex
	lServeMap map[string]*LServe
}

var (
	Serve = &serve{lServeMap: make(map[string]*LServe)}
)

// NewLServe 新建服务
func (o *serve) newLServe(serveId string, opt models.XServerOpt) {
	o.lServeMap[serveId] = &LServe{
		ServeId:    serveId,
		XServerOpt: opt,
		IsWorking:  false,
	}
}

// LoadOfDb load sub serve
func (o *serve) LoadOfDb() {
	var serves []models.XServer
	orm.DbFindBy(&serves, "role > ?", models.ServeTypeLocal)
	for _, v := range serves {
		o.newLServe(v.Guid, v.XServerOpt)
	}
}

// GetByType 根据类型
func (o *serve) GetByType(ctype int) *LServe {
	o.lock.RLock()
	defer o.lock.RUnlock()
	for _, v := range o.lServeMap {
		if v.XServerOpt.Role == ctype && v.IsWorking {
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

// UpdateAll 更新状态
func (o *serve) UpdateAll() {
	o.lock.Lock()
	defer o.lock.Unlock()
	for _, v := range o.lServeMap {
		if time.Since(v.UpdatedTime) < 60*time.Second {
			continue
		}
		v.IsWorking = false
		v.UpdatedTime = time.Time{}
	}
}

// Get 获取
func (o *serve) Get(serveId, address string) (*LServe, error) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	v, ok := o.lServeMap[serveId]
	if !ok {
		return nil, rpc.ErrServeNoExist
	}
	if v.Status == models.ServeStatusStoped {
		return nil, rpc.ErrServeStoped
	}
	// create token
	v.Token = internal.UUID()
	v.Address = address
	v.IsWorking = true
	v.UpdatedTime = time.Now()
	return v, nil
}

// Add 添加
func (o *serve) Add(srv *models.XServer) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.newLServe(srv.Guid, srv.XServerOpt)
}

// Refresh 刷新
func (o *serve) Refresh(serveId, token string) error {
	o.lock.Lock()
	defer o.lock.Unlock()
	v, ok := o.lServeMap[serveId]
	if !ok {
		return rpc.ErrServeNoExist
	}
	if v.Status == models.ServeStatusStoped {
		return rpc.ErrServeStoped
	}
	if v.Token != token {
		return rpc.ErrAuthority
	}
	v.UpdatedTime = time.Now()
	return nil
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
		orm.DbUpdateColsBy(&models.XServer{}, orm.H{"status": status}, "guid like ?", s)
	}
	return nil
}

// UpdateStatus 更新状态
func (o *serve) Delete(serveIds []string) error {
	o.lock.Lock()
	defer o.lock.Unlock()
	for _, s := range serveIds {
		delete(o.lServeMap, s)
	}
	return nil
}
