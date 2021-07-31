package model

// ServeOpt 服务配置信息
type ServeOpt struct {
	Name       string `json:"name" gorm:"not null;unique;comment:名称;"`
	Role       int    `json:"role" gorm:"comment:角色;"`
	HttpPort   uint16 `json:"httpPort" gorm:"comment:端口号;"`
	AccessPort uint16 `json:"accessPort" gorm:"comment:设备接入端口号;"`
	RpcPort    uint16 `json:"rpcPort" gorm:"comment:rpc服务端口号;"`
	Status     int    `json:"status" gorm:"comment:服务状态 0-停止 1-启动;"`
}

// Serve 服务详细信息
type Serve struct {
	Model
	Guid string `json:"guid" gorm:"comment:唯一标识;"`
	ServeOpt
}

// TableName 表名
func (s *Serve) TableName() string {
	return "t_serve"
}
