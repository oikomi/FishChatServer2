package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/user-api/conf"
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

func (registerRPCCli *RegisterRPCCli) Register(authReq *rpc.RGRegisterReq) (res *rpc.RGRegisterRes, err error) {
	r := rpc.NewRegisterServerRPCClient(registerRPCCli.conn)
	if res, err = r.Register(context.Background(), authReq); err != nil {
		glog.Error(err)
	}
	return
}

func (registerRPCCli *RegisterRPCCli) Auth(authReq *rpc.RGAuthReq) (res *rpc.RGAuthRes, err error) {
	r := rpc.NewRegisterServerRPCClient(registerRPCCli.conn)
	if res, err = r.Auth(context.Background(), authReq); err != nil {
		glog.Error(err)
	}
	return
}
