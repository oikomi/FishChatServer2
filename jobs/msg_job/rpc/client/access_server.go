package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/conf"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type AccessServerRPCCli struct {
	conn *grpc.ClientConn
}

func NewAccessServerRPCCli() (accessServerRPCCli *AccessServerRPCCli, err error) {
	r := sd.NewResolver(conf.Conf.RPCClient.AccessClient.ServiceName)
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial(conf.Conf.RPCClient.AccessClient.EtcdAddr, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		glog.Error(err)
		panic(err)
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
