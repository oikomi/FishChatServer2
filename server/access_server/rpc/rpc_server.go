package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/access_server/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
}

func (s *RPCServer) SendP2PMsg(ctx context.Context, in *rpc.ASSendP2PMsgReq) (res *rpc.ASSendP2PMsgRes, err error) {
	glog.Info("access_server recive SendP2PMsg")
	res = &rpc.ASSendP2PMsgRes{}
	return
}

func RPCServerInit() {
	glog.Info("[access_server] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpc.RegisterAccessServerRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
