package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
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

func (idgenRPCCli *IdgenRPCCli) GetUUID(getUUIDReq *rpc.Snowflake_NullRequest) (res *rpc.Snowflake_UUID, err error) {
	i := rpc.NewIDGenServerRPCClient(idgenRPCCli.conn)
	if res, err = i.GetUUID(context.Background(), getUUIDReq); err != nil {
		glog.Error(err)
	}
	return
}
