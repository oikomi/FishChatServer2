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

// SayHello implements helloworld.GreeterServer
func (s *RPCServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func RPCServerInit() {
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &RPCServer{})
	s.Serve(lis)
}
