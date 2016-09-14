package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/pb"
	"github.com/oikomi/FishChatServer2/server/router/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
}

func (s *RPCServer) SendMsgP2P(ctx context.Context, in *pb.SendMsgP2PReq) (*pb.SendMsgP2PRes, error) {
	return &pb.SendMsgP2PRes{}, nil
}

func RPCServerInit() {
	glog.Info("[router] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterRouterRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
