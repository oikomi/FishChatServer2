package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/access/conf"
	"github.com/oikomi/FishChatServer2/server/access/global"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
}

func (s *RPCServer) SendP2PMsgFromJob(ctx context.Context, in *rpc.ASSendP2PMsgFromJobReq) (res *rpc.ASSendP2PMsgFromJobRes, err error) {
	glog.Info("access recive SendP2PMsgFromJob")
	glog.Info(global.GSessions)
	if session, ok := global.GSessions[in.TargetUID]; ok {
		glog.Info("session is online")
		if err = session.Send(&external.ResSendP2PMsg{
			Cmd:       external.SendP2PMsgCMD,
			ErrCode:   ecode.OK.Uint32(),
			ErrStr:    ecode.OK.String(),
			SourceUID: in.SourceUID,
			TargetUID: in.TargetUID,
			MsgID:     in.MsgID,
			Msg:       in.Msg,
		}); err != nil {
			glog.Error(err)
			res = &rpc.ASSendP2PMsgFromJobRes{
				ErrCode: ecode.ServerErr.Uint32(),
				ErrStr:  ecode.ServerErr.String(),
			}
			return
		}
	} else {
		// offline msg
	}
	res = &rpc.ASSendP2PMsgFromJobRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func (s *RPCServer) SendNotify(ctx context.Context, in *rpc.ASSendNotifyReq) (res *rpc.ASSendNotifyRes, err error) {
	glog.Info("access recive SendNotify")
	glog.Info(global.GSessions)
	if session, ok := global.GSessions[in.UID]; ok {
		glog.Info("session is online")
		if err = session.Send(&external.ResNotify{
			Cmd:       external.NotifyCMD,
			ErrCode:   ecode.OK.Uint32(),
			ErrStr:    ecode.OK.String(),
			CurrentID: in.CurrentID,
		}); err != nil {
			glog.Error(err)
			res = &rpc.ASSendNotifyRes{
				ErrCode: ecode.ServerErr.Uint32(),
				ErrStr:  ecode.ServerErr.String(),
			}
			return
		}
	} else {
		// offline
	}
	res = &rpc.ASSendNotifyRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

// func (s *RPCServer) SendGroupMsg(ctx context.Context, in *rpc.ASSendGroupMsgReq) (res *rpc.ASSendGroupMsgRes, err error) {
// 	glog.Info("access recive SendGroupMsg")
// 	return
// }

func RPCServerInit() {
	glog.Info("[access] rpc server init at " + conf.Conf.RPCServer.Addr)
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
	rpc.RegisterAccessServerRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
