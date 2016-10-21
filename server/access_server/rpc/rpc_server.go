package rpc

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/access_server/conf"
	// "golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
}

func RPCServerInit() {
	glog.Info("[access_server] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpc.RegisterAccessServerRPCServer(s, &RPCServer{})
	s.Serve(lis)
}
