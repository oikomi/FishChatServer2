package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/access/conf"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type LogicRPCCli struct {
	conn *grpc.ClientConn
}

func NewLogicRPCCli() (logicRPCCli *LogicRPCCli, err error) {
	r := sd.NewResolver(conf.Conf.RPCClient.LogicClient.ServiceName)
	b := grpc.RoundRobin(r)
	conn, err := grpc.Dial(conf.Conf.RPCClient.LogicClient.EtcdAddr, grpc.WithInsecure(), grpc.WithBalancer(b))
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	logicRPCCli = &LogicRPCCli{
		conn: conn,
	}
	return
}

func (logicRPCCli *LogicRPCCli) Login(loginReq *rpc.LoginReq) (res *rpc.LoginRes, err error) {
	l := rpc.NewLogicRPCClient(logicRPCCli.conn)
	if res, err = l.Login(context.Background(), loginReq); err != nil {
		glog.Error(err)
	}
	return
}

func (logicRPCCli *LogicRPCCli) Ping(pingReq *rpc.PingReq) (res *rpc.PingRes, err error) {
	l := rpc.NewLogicRPCClient(logicRPCCli.conn)
	if res, err = l.Ping(context.Background(), pingReq); err != nil {
		glog.Error(err)
	}
	return
}

func (logicRPCCli *LogicRPCCli) SendP2PMsg(sendP2PMsgReq *rpc.SendP2PMsgReq) (res *rpc.SendP2PMsgRes, err error) {
	l := rpc.NewLogicRPCClient(logicRPCCli.conn)
	if res, err = l.SendP2PMsg(context.Background(), sendP2PMsgReq); err != nil {
		glog.Error(err)
	}
	return
}

func (logicRPCCli *LogicRPCCli) AcceptP2PMsgAck(acceptP2PMsgAckReq *rpc.AcceptP2PMsgAckReq) (res *rpc.AcceptP2PMsgAckRes, err error) {
	l := rpc.NewLogicRPCClient(logicRPCCli.conn)
	if res, err = l.AcceptP2PMsgAck(context.Background(), acceptP2PMsgAckReq); err != nil {
		glog.Error(err)
	}
	return
}

func (logicRPCCli *LogicRPCCli) SendGroupMsg(sendGroupMsgReq *rpc.SendGroupMsgReq) (res *rpc.SendGroupMsgRes, err error) {
	l := rpc.NewLogicRPCClient(logicRPCCli.conn)
	if res, err = l.SendGroupMsg(context.Background(), sendGroupMsgReq); err != nil {
		glog.Error(err)
	}
	return
}

func (logicRPCCli *LogicRPCCli) SyncMsg(sendGroupMsgReq *rpc.SyncMsgReq) (res *rpc.SyncMsgRes, err error) {
	l := rpc.NewLogicRPCClient(logicRPCCli.conn)
	if res, err = l.SyncMsg(context.Background(), sendGroupMsgReq); err != nil {
		glog.Error(err)
	}
	return
}
