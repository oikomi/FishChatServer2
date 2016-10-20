package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/access_server/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type MsgServerRPCCli struct {
	conn *grpc.ClientConn
}

func NewMsgServerRPCCli() (msgServerRPCCli *MsgServerRPCCli, err error) {
	conn, err := grpc.Dial(conf.Conf.RPCClient.MsgServerClient.Addr, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return
	}
	msgServerRPCCli = &MsgServerRPCCli{
		conn: conn,
	}
	return
}

func (msgServerRPCCli *MsgServerRPCCli) Login(loginReq *rpc.LoginReq) (res *rpc.LoginRes, err error) {
	m := rpc.NewMsgServerRPCClient(msgServerRPCCli.conn)
	if res, err = m.Login(context.Background(), loginReq); err != nil {
		glog.Error(err)
	}
	return
}

func (msgServerRPCCli *MsgServerRPCCli) SendP2PMsg(sendP2PMsgReq *rpc.SendP2PMsgReq) (res *rpc.SendP2PMsgRes, err error) {
	m := rpc.NewMsgServerRPCClient(msgServerRPCCli.conn)
	if res, err = m.SendP2PMsg(context.Background(), sendP2PMsgReq); err != nil {
		glog.Error(err)
	}
	return
}
