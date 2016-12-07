package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/logic/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type ManagerRPCCli struct {
	conn *grpc.ClientConn
}

func NewManagerRPCCli() (managerRPCCli *ManagerRPCCli, err error) {
	conn, err := grpc.Dial(conf.Conf.RPCClient.ManagerClient.Addr, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return
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
