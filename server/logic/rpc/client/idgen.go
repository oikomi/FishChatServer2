package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/logic/conf"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strconv"
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

func (idgenRPCCli *IdgenRPCCli) Next(ctx context.Context, targetUID int64) (res *rpc.Snowflake_Value, err error) {
	i := rpc.NewIDGenServerRPCClient(idgenRPCCli.conn)
	if res, err = i.Next(ctx, &rpc.Snowflake_Key{Name: strconv.FormatInt(targetUID, 10)}); err != nil {
		glog.Error(err)
	}
	return
}
