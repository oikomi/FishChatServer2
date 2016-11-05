package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/router/conf"
	"github.com/oikomi/FishChatServer2/server/router/dao"
	"github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	dao *dao.Dao
}

func (s *RPCServer) RouterAccess(ctx context.Context, in *rpc.RTAccessReq) (res *rpc.RTAccessRes, err error) {
	return
}

func RPCServerInit() {
	glog.Info("[router] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpcServer := &RPCServer{
		dao: dao.NewDao(),
	}
	rpc.RegisterRouterRPCServer(s, rpcServer)
	s.Serve(lis)
}

func SDHeart() {
	work := etcd.NewWorker(conf.Conf.Etcd.Name, conf.Conf.RPCServer.Addr, conf.Conf.Etcd.Root, conf.Conf.Etcd.Addrs)
	go work.HeartBeat()
}
