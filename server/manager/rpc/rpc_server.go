package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/manager/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
}

func (s *RPCServer) Login(ctx context.Context, in *rpc.LoginReq) (*rpc.LoginRes, error) {
	return &rpc.LoginRes{}, nil
}

func RPCServerInit() {
	glog.Info("[manager] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpc.RegisterManagerRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
