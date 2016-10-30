package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/auth/conf"
	"github.com/oikomi/FishChatServer2/server/auth/dao"
	"github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	dao *dao.Dao
}

func (s *RPCServer) Login(ctx context.Context, in *rpc.AuthLoginReq) (res *rpc.AuthLoginRes, err error) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(fmt.Sprintf("%d", in.UID)))
	md5Ctx.Write([]byte(conf.Conf.Auth.Salt))
	cipherStr := md5Ctx.Sum(nil)
	calcToken := hex.EncodeToString(cipherStr)
	if err = s.dao.SetToken(ctx, in.UID, calcToken); err != nil {
		res = &rpc.AuthLoginRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	res = &rpc.AuthLoginRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func (s *RPCServer) Auth(ctx context.Context, in *rpc.AuthAuthReq) (res *rpc.AuthAuthRes, err error) {
	var token string
	if token, err = s.dao.Token(ctx, in.UID); err != nil {
		res = &rpc.AuthAuthRes{
			ErrCode: ecode.CalcTokenFailed.Uint32(),
			ErrStr:  ecode.CalcTokenFailed.String(),
		}
		return
	}
	if token != in.Token {
		res = &rpc.AuthAuthRes{
			ErrCode: ecode.CalcTokenFailed.Uint32(),
			ErrStr:  ecode.CalcTokenFailed.String(),
		}
		return
	}
	res = &rpc.AuthAuthRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func RPCServerInit() {
	glog.Info("[auth] rpc server init")
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpcServer := &RPCServer{
		dao: dao.NewDao(),
	}
	rpc.RegisterAuthServerRPCServer(s, rpcServer)
	s.Serve(lis)
}

func SDHeart() {
	work := etcd.NewWorker(conf.Conf.Etcd.Name, conf.Conf.RPCServer.Addr, conf.Conf.Etcd.Root, conf.Conf.Etcd.Addrs)
	go work.HeartBeat()
}
