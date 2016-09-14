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

func (m *ManagerRPCCli) Login(uid int64, token string) (err error) {
	c := pb.NewManagerRPCClient(m.conn)
	r, err := c.Login(context.Background(), &pb.LoginReq{})
	if err != nil {
		glog.Error(err)
		return
	}
	glog.Info(r.ErrCode)
	return
}
