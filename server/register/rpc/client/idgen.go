package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"google.golang.org/grpc"
)

type IdgenRPCCli struct {
	conn *grpc.ClientConn
}

func NewIdgenRPCCli() (idgenRPCCli *IdgenRPCCli, err error) {
	r := sd.NewResolver(conf.Conf.RPCClient.IdgenClient.ServiceName)
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial(conf.Conf.RPCClient.IdgenClient.EtcdAddr, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	idgenRPCCli = &IdgenRPCCli{
		conn: conn,
	}
	return
}
