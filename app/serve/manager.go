package serve

import (
	"time"
	"xstation/internal"
	"xstation/models"
	"xstation/pkg/orm"
	"xstation/pkg/rpc"
)

// lServe 服务
type LServe struct {
	ServeId     string    `json:"serveId"` // 服务id
	Token       string    `json:"token"`   //
	UpdatedTime time.Time `json:"updatedTime"`
	models.XServerOpt
}

var (
	lServeMap = make(map[string]*LServe)
)

// NewLServe 新建服务
func newLServe(serveId string, opt models.XServerOpt) {
	lServeMap[serveId] = &LServe{
		ServeId:    serveId,
		XServerOpt: opt,
	}
}

// dbLoadOtherServe load sub serve
func dbLoadOtherServe() {
	var serves []models.XServer
	orm.DbFindBy(&serves, "role > ?", models.ServeTypeLocal)
	for _, v := range serves {
		opt := v.XServerOpt
		if opt.Status != models.ServeStatusStoped {
			opt.Status = models.ServeStatusIdel
		}
		newLServe(v.Guid, opt)
	}
}

// GetServeOfType 根据类型获取服务信息
func GetServeOfType(ctype int) *models.XServerOpt {
	for _, v := range lServeMap {
		if v.XServerOpt.Role == ctype && v.XServerOpt.Status == models.ServeStatusRunning {
			return &v.XServerOpt
		}
	}
	return nil
}

// LoadAllLServe 获取所有服务状态
func LoadAllLServe() (lse []LServe) {
	for _, v := range lServeMap {
		lse = append(lse, *v)
	}
	return
}

// LoadLServe 获取服务信息
func loadLServe(serveId, address string) (*LServe, error) {
	v, ok := lServeMap[serveId]
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

// updateLServe 更新
func updateLServe(serveId, token string, update time.Time) error {
	v, ok := lServeMap[serveId]
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
func ChangeStatus(serveId string, status int) error {
	v, ok := lServeMap[serveId]
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
