package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/manager/conf"
	"github.com/oikomi/FishChatServer2/server/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
}

// SayHello implements helloworld.GreeterServer
func (s *RPCServer) Login(ctx context.Context, in *pb.LoginReq) (*pb.LoginRes, error) {
	return &pb.LoginRes{}, nil
}

func RPCServerInit() {
	glog.Info("[manager] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterManagerRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
