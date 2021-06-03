package rpc

import (
	"context"
	"fmt"

	"github.com/smallnest/rpcx/client"
)

var xclient client.XClient

// var xchannel chan *protocol.Message

// Connect 连接
func Connect(servePath, ipAddr string, port uint16) error {
	addr := fmt.Sprintf("%s:%d", ipAddr, port)
	d, err := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	if err != nil {
		return err
	}
	// xchannel = make(chan *protocol.Message)
	xclient = client.NewXClient(servePath, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	return nil
}

// Call 发送消息
func Call(method string, args interface{}, reply interface{}) error {
	return xclient.Call(context.Background(), method, args, reply)
}

// Login 登录
func Login(args *LoginArgs, reply *LoginReply) error {
	return Call("Login", args, reply)
}

// KeepAlive 保持心跳
func KeepAlive(args *KeepAliveArgs) error {
	return Call("KeepAlive", args, nil)
}
