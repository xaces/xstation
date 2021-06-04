package models

const (
	ServeTypeUnKnow = iota // 未知
	ServeTypeLocal         // 本机服务
	ServeTypeStream        // 流媒体服务
)

const (
	ServeStatusIdel    = iota // ServeIdel 服务空闲
	ServeStatusRunning        // ServeRunning 运行
	ServeStatusStoped         // 停止
)

// ServerOpt 服务配置信息
type XServerOpt struct {
	Name       string `json:"name" gorm:"not null;unique;comment:名称;"`
	Role       int    `json:"role" gorm:"comment:角色;"`
	Port       uint16 `json:"port" gorm:"comment:端口号;"`
	AccessPort uint16 `json:"accessPort" gorm:"comment:设备接入端口号;"`
	RpcPort    uint16 `json:"rpcPort" gorm:"comment:rpc服务端口号;"`
	Status     int    `json:"status" gorm:"comment:服务状态 0-停止 1-启动;"`
	Address    string `json:"address" gorm:"comment:服务IP，由服务主动上报;"`
}

// Server 服务详细信息
type XServer struct {
	Model
	Guid string `json:"guid" gorm:"comment:唯一标识;"`
	XServerOpt
}

// TableName 表名
func (s *XServer) TableName() string {
	return "t_xserver"
}
