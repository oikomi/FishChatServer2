package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/logic/conf"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ManagerRPCCli struct {
	conn *grpc.ClientConn
}

func NewManagerRPCCli() (managerRPCCli *ManagerRPCCli, err error) {
	r := sd.NewResolver(conf.Conf.RPCClient.ManagerClient.ServiceName)
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial(conf.Conf.RPCClient.ManagerClient.EtcdAddr, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	managerRPCCli = &ManagerRPCCli{
		conn: conn,
	}
	return
}

func (managerRPCCli *ManagerRPCCli) SetExceptionMsg(ctx context.Context, sourceUID, targetUID int64, msgID, msg string) (res *rpc.MGExceptionMsgRes, err error) {
	m := rpc.NewManagerServerRPCClient(managerRPCCli.conn)
	if res, err = m.SetExceptionMsg(ctx, &rpc.MGExceptionMsgReq{SourceUID: sourceUID, TargetUID: targetUID, MsgID: msgID, Msg: msg}); err != nil {
		glog.Error(err)
	}
	return
}

func (managerRPCCli *ManagerRPCCli) SyncMsg(ctx context.Context, uid, currentID, totalID int64) (res *rpc.MGSyncMsgRes, err error) {
	m := rpc.NewManagerServerRPCClient(managerRPCCli.conn)
	if res, err = m.Sync(ctx, &rpc.MGSyncMsgReq{UID: uid, CurrentID: currentID, TotalID: totalID}); err != nil {
		glog.Error(err)
	}
	return
}
