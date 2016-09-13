package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RouterRPCCli struct {
}

func NewRouterRPCCli(address string) (err error) {
	glog.Info("NewRouterRPCCli")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "miaohong"})
	if err != nil {
		glog.Error(err)
	}
	glog.Info(r.Message)
	return
}
