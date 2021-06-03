package rpc

import (
	"errors"
	"time"
)

// 和中心服务通信
// func name <Register>

const (
	ServeStatusIdel     = iota // 空闲
	ServeStatusOk              // 正常
	ServeStatusDisabled        // 禁止
)

var (
	ErrServeNoExist = errors.New("serve no exist")                // ErrServeNoExist 服务不存在
	ErrServeStoped  = errors.New("the service has been disabled") // ErrServeStoped 服务已停止
	ErrAuthority    = errors.New("authority error")               // ErrAuthority 授权失败
)

// KeepAliveArgs 保活
type KeepAliveArgs struct {
	ServeId     string    `json:"serveId"`     // 服务ID
	Token       string    `json:"token"`       // 授权token
	UpdatedTime time.Time `json:"updatedTime"` // 更新时间
}

// LoginArgs 登录工作站
type LoginArgs struct {
	ServeId string `json:"serveId"` // 服务ID
	Address string `json:"address"` // 服务名称
}

// LoginReply 工作站响应
type LoginReply struct {
	Token string `json:"token"` // 授权token
}

// XLinkRegister 设备链路注册
// for external server
type XLinkRegister struct {
	Data interface{} `json:"data"`
}
