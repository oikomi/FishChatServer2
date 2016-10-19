package client

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/msg_server/conf"
	// "golang.org/x/net/context"
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

// func (m *ManagerRPCCli) Login(uid int64, token string) (err error) {
// 	c := rpc.NewManagerRPCClient(m.conn)
// 	r, err := c.Login(context.Background(), &rpc.LoginReq{})
// 	if err != nil {
// 		glog.Error(err)
// 		return
// 	}
// 	glog.Info(r.ErrCode)
// 	return
// }
