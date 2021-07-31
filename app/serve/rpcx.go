package serve

import (
	"context"
	"fmt"
	"log"
	"xstation/app/mnger"

	"github.com/wlgd/xutils/rpc"

	"github.com/smallnest/rpcx/server"
)

// Arith 标识
type Arith int

// Login 子服务登录
func (t *Arith) Login(cxt context.Context, args *rpc.LoginArgs, reply *rpc.LoginReply) error {
	if s := mnger.Serve.Get(args.ServeId); s == nil {
		return rpc.ErrServeNoExist
	} else {
	}
	return nil
}

// KeepAlive 工作站保活
func (t *Arith) KeepAlive(cxt context.Context, args *rpc.KeepAliveArgs, reply *rpc.KeepAliveArgs) error {
	return nil
}

var (
	_rpcx *server.Server = nil
)

// rpcxStart start rpc server
func rpcxStart(port uint16) error {
	address := fmt.Sprintf(":%d", port)
	_rpcx = server.NewServer()
	_rpcx.RegisterName("xstation", new(Arith), "")
	go _rpcx.Serve("tcp", address)
	return nil
}

// rpcxStop stop rpc server
func rpcxStop() {
	log.Printf("localRpcServe Close %v", _rpcx.Close())
}
