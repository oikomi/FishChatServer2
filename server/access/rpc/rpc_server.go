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

func (s *RPCServer) SendP2PMsg(ctx context.Context, in *rpc.ASSendP2PMsgReq) (res *rpc.ASSendP2PMsgRes, err error) {
	glog.Info("access recive SendP2PMsg")
	glog.Info(global.GSessions)
	if session, ok := global.GSessions[in.TargetUID]; ok {
		glog.Info("session is online")
		if err = session.Send(&external.ResSendP2PMsg{
			Cmd:       external.SendP2PMsgCMD,
			SourceUID: in.SourceUID,
			TargetUID: in.TargetUID,
			MsgID:     in.MsgID,
			Msg:       in.Msg,
		}); err != nil {
			glog.Error(err)
			res = &rpc.ASSendP2PMsgRes{
				ErrCode: ecode.ServerErr.Uint32(),
				ErrStr:  ecode.ServerErr.String(),
			}
			return
		}
	} else {
		// offline msg
	}
	res = &rpc.ASSendP2PMsgRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func RPCServerInit() {
	glog.Info("[access] rpc server init")
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
