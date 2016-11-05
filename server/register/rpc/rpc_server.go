package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	"github.com/oikomi/FishChatServer2/server/register/dao"
	"github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	dao *dao.Dao
}

func (s *RPCServer) Online(ctx context.Context, in *rpc.RGOnlineReq) (res *rpc.RGOnlineRes, err error) {
	if _, err = s.dao.GetOnline(ctx, in.UID); err != nil {
		glog.Error(err)
		res = &rpc.RGOnlineRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
			Online:  false,
		}
		return
	}
	res = &rpc.RGOnlineRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
		Online:  true,
	}
	return
}

func (s *RPCServer) Ping(ctx context.Context, in *rpc.RGPingReq) (res *rpc.RGPingRes, err error) {
	if err = s.dao.SetOnline(ctx, in.UID); err != nil {
		glog.Error(err)
		res = &rpc.RGPingRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	res = &rpc.RGPingRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func RPCServerInit() {
	glog.Info("[register] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpcServer := &RPCServer{
		dao: dao.NewDao(),
	}
	rpc.RegisterRegisterServerRPCServer(s, rpcServer)
	s.Serve(lis)
}

func SDHeart() {
	work := etcd.NewWorker(conf.Conf.Etcd.Name, conf.Conf.RPCServer.Addr, conf.Conf.Etcd.Root, conf.Conf.Etcd.Addrs)
	go work.HeartBeat()
}
