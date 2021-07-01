package serve

import (
	"context"
	"fmt"
	"log"
	"xstation/app/manager"
	"xstation/service"

	"github.com/wlgd/xutils/rpc"

	"github.com/smallnest/rpcx/server"
	"github.com/wlgd/xproto"
)

// Arith 标识
type Arith int

// Login 子服务登录
func (t *Arith) Login(cxt context.Context, args *rpc.LoginArgs, reply *rpc.LoginReply) error {
	s, err := manager.Serve.LoadLServe(args.ServeId, args.Address)
	if err != nil {
		return err
	}
	reply.Token = s.Token
	return nil
}

// KeepAlive 工作站保活
func (t *Arith) KeepAlive(cxt context.Context, args *rpc.KeepAliveArgs, reply *rpc.KeepAliveArgs) error {
	return manager.Serve.UpdateLServe(args.ServeId, args.Token, args.UpdatedTime)
}

// XLinkRegister 服务注册
func (t *Arith) XLinkRegister(cxt context.Context, args *rpc.XLinkRegister, reply *rpc.XLinkRegister) error {
	link, ok := args.Data.(xproto.LinkAccess)
	if !ok {
		return rpc.ErrParameter
	}
	return service.NewXData().DbUpdateAccess(&link)
}

var (
	_rpcx *server.Server = nil
)

// rpcxStart start rpc server
func rpcxStart(port uint16) {
	address := fmt.Sprintf(":%d", port)
	_rpcx = server.NewServer()
	_rpcx.RegisterName("xstation", new(Arith), "")
	log.Printf("loaclRpcServe Start %d\n", port)
	if err := _rpcx.Serve("tcp", address); err != nil {
		log.Fatal(err)
	}
}

// rpcxStop stop rpc server
func rpcxStop() {
	log.Printf("localRpcServe Close %v", _rpcx.Close())
}
