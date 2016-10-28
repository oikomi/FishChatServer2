package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/jobs/msg_job/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type AccessServerRPCCli struct {
	conn *grpc.ClientConn
}

func NewAccessServerRPCCli() (accessServerRPCCli *AccessServerRPCCli, err error) {
	conn, err := grpc.Dial(conf.Conf.RPCClient.AccessServerClient.Addr, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return
	}
	accessServerRPCCli = &AccessServerRPCCli{
		conn: conn,
	}
	return
}

func (accessServerRPCCli *AccessServerRPCCli) SendP2PMsg(sendP2PMsgReq *rpc.ASSendP2PMsgReq) (res *rpc.ASSendP2PMsgRes, err error) {
	a := rpc.NewAccessServerRPCClient(accessServerRPCCli.conn)
	if res, err = a.SendP2PMsg(context.Background(), sendP2PMsgReq); err != nil {
		glog.Error(err)
	}
	return
}
