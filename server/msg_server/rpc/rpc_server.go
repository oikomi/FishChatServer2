package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/msg_server/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
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
	rpc.RegisterMsgServerRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
