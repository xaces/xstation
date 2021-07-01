package manager

import (
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
	models.XServerOpt
}

type serveManager struct {
	lServeMap map[string]*LServe
}

var (
	Serve = &serveManager{lServeMap: make(map[string]*LServe)}
)

// NewLServe 新建服务
func (o *serveManager) newLServe(serveId string, opt models.XServerOpt) {
	o.lServeMap[serveId] = &LServe{
		ServeId:    serveId,
		XServerOpt: opt,
	}
}

// LoadOfDb load sub serve
func (o *serveManager) LoadOfDb() {
	var serves []models.XServer
	orm.DbFindBy(&serves, "role > ?", models.ServeTypeLocal)
	for _, v := range serves {
		opt := v.XServerOpt
		if opt.Status != models.ServeStatusStoped {
			opt.Status = models.ServeStatusIdel
		}
		o.newLServe(v.Guid, opt)
	}
}

// GetServeByType 根据类型获取服务信息
func (o *serveManager) GetServeByType(ctype int) *models.XServerOpt {
	for _, v := range o.lServeMap {
		if v.XServerOpt.Role == ctype && v.XServerOpt.Status == models.ServeStatusRunning {
			return &v.XServerOpt
		}
	}
	return nil
}

// LoadAllLServe 获取所有服务状态
func (o *serveManager) LoadAllLServe() (lse []LServe) {
	for _, v := range o.lServeMap {
		lse = append(lse, *v)
	}
	return
}

// LoadLServe 获取服务信息
func (o *serveManager) LoadLServe(serveId, address string) (*LServe, error) {
	v, ok := o.lServeMap[serveId]
	if !ok {
		return v, rpc.ErrServeNoExist
	}

	if v.Status == models.ServeStatusStoped {
		return v, rpc.ErrServeStoped
	}
	// create token
	v.Token = internal.UUID()
	if v.Address != address {
		orm.DbUpdateColsBy(&models.XServer{}, orm.H{"address": address}, "guid like ?", serveId)
	}
	v.Status = models.ServeStatusRunning
	return v, nil
}

// UpdateLServe 更新
func (o *serveManager) UpdateLServe(serveId, token string, update time.Time) error {
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
	v.UpdatedTime = update
	return nil
}

// Shutdown 停止服务
func (o *serveManager) ChangeStatus(serveId string, status int) error {
	v, ok := o.lServeMap[serveId]
	if !ok {
		return rpc.ErrServeNoExist
	}
	if status == models.ServeStatusStoped {
		v.Status = models.ServeStatusStoped
	} else {
		v.Status = models.ServeStatusIdel
	}
	orm.DbUpdateColsBy(&models.XServer{}, orm.H{"status": status}, "guid like ?", serveId)
	return nil
}
