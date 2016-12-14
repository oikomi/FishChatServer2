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
	"time"
)

type RPCServer struct {
}

func (s *RPCServer) SendP2PMsg(ctx context.Context, in *rpc.ASSendP2PMsgReq) (res *rpc.ASSendP2PMsgRes, err error) {
	glog.Info("access recive SendP2PMsg")
	if session, ok := global.GSessions[in.TargetUID]; ok {
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
	err := sd.Register("access_server", "127.0.0.1", 20000, conf.Conf.Etcd.Addrs[0], time.Second*3, 5)
	if err != nil {
		panic(err)
	}
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpc.RegisterAccessServerRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
