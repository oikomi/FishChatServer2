package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	commmodel "github.com/oikomi/FishChatServer2/common/model"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/logic/conf"
	"github.com/oikomi/FishChatServer2/server/logic/dao"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	rpcClient *RPCClient
	dao       *dao.Dao
}

func (s *RPCServer) Login(ctx context.Context, in *rpc.LoginReq) (res *rpc.LoginRes, err error) {
	glog.Info("logic recive login")
	// FIXME
	if in.Token == "" || in.UID < 0 {
		res = &rpc.LoginRes{
			ErrCode: ecode.NoToken.Uint32(),
			ErrStr:  ecode.NoToken.String(),
		}
		return
	}
	rgRes, err := s.rpcClient.Register.Login(ctx, in.UID, in.Token, in.AccessAddr)
	if err != nil {
		res = &rpc.LoginRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		glog.Error(err)
		return
	}
	// success
	res = &rpc.LoginRes{
		ErrCode: rgRes.ErrCode,
		ErrStr:  rgRes.ErrStr,
	}
	return
}

func (s *RPCServer) Ping(ctx context.Context, in *rpc.PingReq) (res *rpc.PingRes, err error) {
	glog.Info("logic recive Ping")
	// FIXME
	if in.UID < 0 {
		res = &rpc.PingRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	// check
	rgPingRes, err := s.rpcClient.Register.Ping(ctx, in.UID)
	if err != nil {
		res = &rpc.PingRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		glog.Error(err)
		return
	}
	// success
	res = &rpc.PingRes{
		ErrCode: rgPingRes.ErrCode,
		ErrStr:  rgPingRes.ErrStr,
	}
	return
}

func (s *RPCServer) SendP2PMsg(ctx context.Context, in *rpc.SendP2PMsgReq) (res *rpc.SendP2PMsgRes, err error) {
	glog.Info("logic recive SendP2PMsg")
	sendP2PMsgKafka := &commmodel.SendP2PMsgKafka{
		SourceUID: in.SourceUID,
		TargetUID: in.TargetUID,
		MsgID:     in.MsgID,
		MsgType:   in.MsgType,
		Msg:       in.Msg,
	}
	// idgen
	idgenRes, err := s.rpcClient.Idgen.Next(ctx, in.TargetUID)
	if err != nil {
		glog.Error(err)
		return
	}
	sendP2PMsgKafka.IncrementID = idgenRes.Value
	// Online
	onlineRes, err := s.rpcClient.Register.Online(ctx, in.TargetUID)
	if err != nil {
		glog.Error(err)
		return
	}
	if !onlineRes.Online {
		glog.Info(in.TargetUID, " is offline")
		sendP2PMsgKafka.Online = false
	} else {
		sendP2PMsgKafka.Online = true
		// send notify to client
		notifyRes, err := s.rpcClient.Notify.Notify(ctx, in.TargetUID, idgenRes.Value)
		if err != nil {
			glog.Error(err)
		}
		res = &rpc.SendP2PMsgRes{
			ErrCode: notifyRes.ErrCode,
			ErrStr:  notifyRes.ErrStr,
		}
	}
	// get access server Addr
	routerAccessRes, err := s.rpcClient.Register.RouterAccess(ctx, in.TargetUID)
	if err != nil {
		glog.Error(err)
		return
	}
	sendP2PMsgKafka.AccessServerAddr = routerAccessRes.AccessAddr
	// send to kafka
	glog.Info("send to kafka", sendP2PMsgKafka)
	s.dao.KafkaProducer.SendP2PMsg(sendP2PMsgKafka)
	res = &rpc.SendP2PMsgRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func (s *RPCServer) AcceptP2PMsgAck(ctx context.Context, in *rpc.AcceptP2PMsgAckReq) (res *rpc.AcceptP2PMsgAckRes, err error) {
	glog.Info("logic recive AcceptP2PMsgAck")
	return
}

func (s *RPCServer) SendGroupMsg(ctx context.Context, in *rpc.SendGroupMsgReq) (res *rpc.SendGroupMsgRes, err error) {
	glog.Info("logic recive SendGroupMsg")
	sendGroupMsgKafka := &commmodel.SendGroupMsgKafka{
		SourceUID: in.SourceUID,
		GroupID:   in.GroupID,
		MsgID:     in.MsgID,
		Msg:       in.Msg,
	}
	s.dao.KafkaProducer.SendGroupMsg(sendGroupMsgKafka)
	res = &rpc.SendGroupMsgRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func (s *RPCServer) SyncMsg(ctx context.Context, in *rpc.SyncMsgReq) (res *rpc.SyncMsgRes, err error) {
	glog.Info("logic recive SyncMsg")
	res = &rpc.SyncMsgRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func RPCServerInit(rpcClient *RPCClient) {
	glog.Info("[logic] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	err = sd.Register(conf.Conf.ServiceDiscoveryServer.ServiceName, conf.Conf.ServiceDiscoveryServer.RPCAddr, conf.Conf.ServiceDiscoveryServer.EtcdAddr, conf.Conf.ServiceDiscoveryServer.Interval, conf.Conf.ServiceDiscoveryServer.TTL)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	dao, err := dao.NewDao()
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	rpcServer := &RPCServer{
		rpcClient: rpcClient,
		dao:       dao,
	}
	rpc.RegisterLogicRPCServer(s, rpcServer)
	s.Serve(lis)
}
