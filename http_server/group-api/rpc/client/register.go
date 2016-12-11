package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/group-api/conf"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RegisterRPCCli struct {
	conn *grpc.ClientConn
}

func NewRegisterRPCCli() (registerRPCCli *RegisterRPCCli, err error) {
	conn, err := grpc.Dial(conf.Conf.RPCClient.RegisterClient.Addr, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return
	}
	registerRPCCli = &RegisterRPCCli{
		conn: conn,
	}
	return
}

func (registerRPCCli *RegisterRPCCli) CreateGroup(createGroupReq *rpc.RGCreateGroupReq) (res *rpc.RGCreateGroupRes, err error) {
	r := rpc.NewRegisterServerRPCClient(registerRPCCli.conn)
	if res, err = r.CreateGroup(context.Background(), createGroupReq); err != nil {
		glog.Error(err)
	}
	return
}
