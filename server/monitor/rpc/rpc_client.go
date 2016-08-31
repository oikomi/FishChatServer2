package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/monitor/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type MsgServerRPCCli struct {
}

func NewMsgServerRPCCli(address string) {
	glog.Info("NewMsgServerRPCCli")
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
}
