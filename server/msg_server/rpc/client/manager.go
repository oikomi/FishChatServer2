package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/msg_server/conf"
	"github.com/oikomi/FishChatServer2/server/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ManagerRPCCli struct {
	conn *grpc.ClientConn
}

func NewManagerRPCCli() (managerRPCCli *ManagerRPCCli, err error) {
	conn, err := grpc.Dial(conf.Conf.RPCClient.ManagerClient.Addr, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return
	}
	managerRPCCli = &ManagerRPCCli{
		conn: conn,
	}
	return
}

func (m *ManagerRPCCli) say() {
	c := pb.NewGreeterClient(m.conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "miaohong"})
	if err != nil {
		glog.Error(err)
	}
	glog.Info(r.Message)
}
