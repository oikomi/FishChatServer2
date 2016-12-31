package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/group-api/conf"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RegisterRPCCli struct {
	conn *grpc.ClientConn
}

func NewRegisterRPCCli() (registerRPCCli *RegisterRPCCli, err error) {
	r := sd.NewResolver(conf.Conf.RPCClient.RegisterClient.ServiceName)
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial(conf.Conf.RPCClient.RegisterClient.EtcdAddr, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		glog.Error(err)
		panic(err)
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

func (registerRPCCli *RegisterRPCCli) JoinGroup(joinGroupReq *rpc.RGJoinGroupReq) (res *rpc.RGJoinGroupRes, err error) {
	r := rpc.NewRegisterServerRPCClient(registerRPCCli.conn)
	if res, err = r.JoinGroup(context.Background(), joinGroupReq); err != nil {
		glog.Error(err)
	}
	return
}

func (registerRPCCli *RegisterRPCCli) QuitGroup(quitGroupReq *rpc.RGQuitGroupReq) (res *rpc.RGQuitGroupRes, err error) {
	r := rpc.NewRegisterServerRPCClient(registerRPCCli.conn)
	if res, err = r.QuitGroup(context.Background(), quitGroupReq); err != nil {
		glog.Error(err)
	}
	return
}
