package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/common/model"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/msg_server/conf"
	"github.com/oikomi/FishChatServer2/server/msg_server/job"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	kafkaJob *job.KafkaProducer
}

func (s *RPCServer) Login(ctx context.Context, in *rpc.LoginReq) (res *rpc.LoginRes, err error) {
	glog.Info("msg_server recive login")
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

func (s *RPCServer) SendP2PMsg(ctx context.Context, in *rpc.SendP2PMsgReq) (res *rpc.SendP2PMsgRes, err error) {
	glog.Info("msg_server recive SendP2PMsg")
	sendP2PMsgKafka := &model.SendP2PMsgKafka{
		UID:       in.SourceUID,
		TargetUID: in.TargetUID,
		Msg:       in.Msg,
	}
	s.kafkaJob.SendP2PMsg(sendP2PMsgKafka)
	res = &rpc.SendP2PMsgRes{}
	return
}

func RPCServerInit() {
	glog.Info("[msg_server] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpcServer := &RPCServer{
		kafkaJob: job.NewKafkaProducer(),
	}
	go rpcServer.kafkaJob.HandleError()
	go rpcServer.kafkaJob.HandleSuccess()
	go rpcServer.kafkaJob.Process()
	rpc.RegisterMsgServerRPCServer(s, rpcServer)
	s.Serve(lis)
}
