package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	commmodel "github.com/oikomi/FishChatServer2/common/model"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/logic/conf"
	"github.com/oikomi/FishChatServer2/server/logic/dao"
	"github.com/oikomi/FishChatServer2/server/logic/model"
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
	// success
	res = &rpc.LoginRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func (s *RPCServer) Ping(ctx context.Context, in *rpc.PingReq) (res *rpc.PingRes, err error) {
	glog.Info("logic recive login")
	// FIXME
	if in.UID < 0 {
		res = &rpc.PingRes{
			ErrCode: ecode.NoToken.Uint32(),
			ErrStr:  ecode.NoToken.String(),
		}
		return
	}
	// check

	// success
	res = &rpc.PingRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func (s *RPCServer) SendP2PMsg(ctx context.Context, in *rpc.SendP2PMsgReq) (res *rpc.SendP2PMsgRes, err error) {
	glog.Info("msg_server recive SendP2PMsg")
	sendP2PMsgKafka := &commmodel.SendP2PMsgKafka{
		UID:       in.SourceUID,
		TargetUID: in.TargetUID,
		Msg:       in.Msg,
	}
	// Online
	if _, err = s.rpcClient.Register.Online(in.TargetUID); err != nil {
		glog.Info(in.TargetUID, " is offline")
		// set offline msg
		offlineMsg := &model.OfflineMsg{
			MsgID:     122,
			SourceUID: in.SourceUID,
			TargetUID: in.TargetUID,
			Msg:       in.Msg,
		}
		if err = s.dao.MongoDB.StoreOfflineMsg(offlineMsg); err != nil {
			glog.Error(err)
		}
		return
	}
	s.dao.KafkaProducer.SendP2PMsg(sendP2PMsgKafka)
	res = &rpc.SendP2PMsgRes{}
	return
}

func RPCServerInit(rpcClient *RPCClient) {
	glog.Info("[msg_server] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
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
