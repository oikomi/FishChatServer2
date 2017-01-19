package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/logic/conf"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type NotifyRPCCli struct {
	conn *grpc.ClientConn
}

func NewNotifyRPCCli() (notifyRPCCli *NotifyRPCCli, err error) {
	r := sd.NewResolver(conf.Conf.RPCClient.NotifyClient.ServiceName)
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial(conf.Conf.RPCClient.NotifyClient.EtcdAddr, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	notifyRPCCli = &NotifyRPCCli{
		conn: conn,
	}
	return
}

func (notifyRPCCli *NotifyRPCCli) Notify(ctx context.Context, targetUID, totalID int64, accessAddr string) (res *rpc.NFNotifyMsgRes, err error) {
	n := rpc.NewNotifyServerRPCClient(notifyRPCCli.conn)
	if res, err = n.Notify(ctx, &rpc.NFNotifyMsgReq{
		TargetUID:        targetUID,
		TotalID:          totalID,
		AccessServerAddr: accessAddr,
	}); err != nil {
		glog.Error(err)
	}
	return
}
