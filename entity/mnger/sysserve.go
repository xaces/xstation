package mnger

import (
	"sync"
	"time"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

// ServeMapper 服务
type Serve struct {
	ServeId     string    `json:"serveId"` // 服务id
	Token       string    `json:"token"`   //
	UpdatedTime time.Time `json:"updatedTime"`
	IsWorking   bool      `json:"isWorking"`
	Address     string    `json:"address"`
	model.SysServeOpt
}

func (s *Serve) Update(address string) error {
	s.Address = address
	s.UpdatedTime = time.Now()
	return nil
}

type ServeMapper struct {
	lock      sync.RWMutex
	lServeMap map[string]*Serve
}

var (
	Serves = &ServeMapper{lServeMap: make(map[string]*Serve)}
)

// 初始化
func (s *ServeMapper) Set(serves []model.SysServe) {
	for _, v := range serves {
		s.Add(&v)
	}
}

// GetByType 根据类型
func (s *ServeMapper) GetByType(ctype int) *Serve {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.lServeMap {
		if v.SysServeOpt.Role == ctype && v.IsWorking {
			return v
		}
	}
	return nil
}

// LoadAllLServe 获取实时状态
func (s *ServeMapper) GetAll() (lse []Serve) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.lServeMap {
		lse = append(lse, *v)
	}
	return
}

// Get 获取
func (s *ServeMapper) Get(serveId string) *Serve {
	s.lock.RLock()
	defer s.lock.RUnlock()
	v, ok := s.lServeMap[serveId]
	if !ok {
		return nil
	}
	return v
}

// Add 添加
func (s *ServeMapper) Add(srv *model.SysServe) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lServeMap[srv.Guid] = &Serve{
		ServeId:     srv.Guid,
		SysServeOpt: srv.SysServeOpt,
		IsWorking:   false,
	}
}

// UpdateStatus 更新状态
func (s *ServeMapper) UpdateStatus(serveIds []string, status int) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, id := range serveIds {
		v, ok := s.lServeMap[id]
		if !ok || status == v.Status {
			continue
		}
		v.Status = status
		orm.DbUpdateColsBy(&model.SysServe{}, orm.H{"status": status}, "guid like ?", s)
	}
	return nil
}

// UpdateStatus 更新状态
func (s *ServeMapper) Delete(devices []model.Device) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, d := range devices {
		delete(s.lServeMap, d.Guid)
	}
}
