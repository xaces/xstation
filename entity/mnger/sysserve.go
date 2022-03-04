package mnger

import (
	"sync"
	"time"
	"xstation/model"

	"github.com/wlgd/xutils/orm"
)

type ServeInfo struct {
	ServeId     string    `json:"serveId"` // 服务id
	Token       string    `json:"token"`   //
	UpdatedTime time.Time `json:"updatedTime"`
	IsWorking   bool      `json:"isWorking"`
	Address     string    `json:"address"`
	model.SysServeOpt
}

func (s *ServeInfo) Update(address string) error {
	s.Address = address
	s.UpdatedTime = time.Now()
	return nil
}

// serveMapper 服务
type serveMapper struct {
	lock      sync.RWMutex
	lServeMap map[string]*ServeInfo
}

var (
	Serve = &serveMapper{lServeMap: make(map[string]*ServeInfo)}
)

// 初始化
func (s *serveMapper) Set(serves []model.SysServe) {
	for _, v := range serves {
		s.Add(&v)
	}
}

// GetByType 根据类型
func (s *serveMapper) GetByType(ctype int) *ServeInfo {
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
func (s *serveMapper) GetAll() (lse []ServeInfo) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.lServeMap {
		lse = append(lse, *v)
	}
	return
}

// Get 获取
func (s *serveMapper) Get(serveId string) *ServeInfo {
	s.lock.RLock()
	defer s.lock.RUnlock()
	v, ok := s.lServeMap[serveId]
	if !ok {
		return nil
	}
	return v
}

// Add 添加
func (s *serveMapper) Add(srv *model.SysServe) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lServeMap[srv.Guid] = &ServeInfo{
		ServeId:     srv.Guid,
		SysServeOpt: srv.SysServeOpt,
		IsWorking:   false,
	}
}

// UpdateStatus 更新状态
func (s *serveMapper) UpdateStatus(serveIds []string, status int) error {
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
func (s *serveMapper) Delete(devices []model.Device) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, d := range devices {
		delete(s.lServeMap, d.Guid)
	}
}