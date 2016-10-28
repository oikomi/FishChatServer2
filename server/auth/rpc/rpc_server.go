package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/auth/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
}

func (s *RPCServer) Login(ctx context.Context, in *rpc.AuthLoginReq) (res *rpc.AuthLoginRes, err error) {
	return
}

func RPCServerInit() {
	glog.Info("[auth] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpc.RegisterAuthServerRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
