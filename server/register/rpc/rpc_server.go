package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	"github.com/oikomi/FishChatServer2/server/register/dao"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	dao *dao.Dao
}

func (s *RPCServer) Login(ctx context.Context, in *rpc.RGLoginReq) (res *rpc.RGLoginRes, err error) {
	glog.Info("register recive login")
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(fmt.Sprintf("%d", in.UID)))
	md5Ctx.Write([]byte(conf.Conf.Auth.Salt))
	cipherStr := md5Ctx.Sum(nil)
	calcToken := hex.EncodeToString(cipherStr)
	glog.Info(calcToken)
	if calcToken != in.Token {
		res = &rpc.RGLoginRes{
			ErrCode: ecode.CalcTokenFailed.Uint32(),
			ErrStr:  ecode.CalcTokenFailed.String(),
		}
		return
	}
	if _, err = s.dao.Token(ctx, in.UID); err != nil {
		res = &rpc.RGLoginRes{
			ErrCode: ecode.CalcTokenFailed.Uint32(),
			ErrStr:  ecode.CalcTokenFailed.String(),
		}
		return
	}
	// if err = s.dao.SetToken(ctx, in.UID, calcToken); err != nil {
	// 	res = &rpc.RGLoginRes{
	// 		ErrCode: ecode.ServerErr.Uint32(),
	// 		ErrStr:  ecode.ServerErr.String(),
	// 	}
	// 	return
	// }
	// regster
	glog.Info(in.AccessAddr)
	if err = s.dao.RegisterAccess(ctx, in.UID, in.AccessAddr); err != nil {
		res = &rpc.RGLoginRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		glog.Error(err)
		return
	}
	if err = s.dao.SetOnline(ctx, in.UID); err != nil {
		glog.Error(err)
		res = &rpc.RGLoginRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	res = &rpc.RGLoginRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
		Token:   calcToken,
	}
	return
}

func (s *RPCServer) RouterAccess(ctx context.Context, in *rpc.RGAccessReq) (res *rpc.RGAccessRes, err error) {
	glog.Info("register recive RouterAccess")
	var accessAddr string
	if accessAddr, err = s.dao.RouterAccess(ctx, in.UID); err != nil {
		res = &rpc.RGAccessRes{
			ErrCode: ecode.CalcTokenFailed.Uint32(),
			ErrStr:  ecode.CalcTokenFailed.String(),
		}
		glog.Error(err)
		return
	}
	res = &rpc.RGAccessRes{
		ErrCode:    ecode.OK.Uint32(),
		ErrStr:     ecode.OK.String(),
		AccessAddr: accessAddr,
	}
	return
}

func (s *RPCServer) Auth(ctx context.Context, in *rpc.RGAuthReq) (res *rpc.RGAuthRes, err error) {
	glog.Info("register recive auth")
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(fmt.Sprintf("%d", in.UID)))
	md5Ctx.Write([]byte(conf.Conf.Auth.Salt))
	cipherStr := md5Ctx.Sum(nil)
	calcToken := hex.EncodeToString(cipherStr)
	glog.Info(calcToken)
	if err = s.dao.SetToken(ctx, in.UID, calcToken); err != nil {
		res = &rpc.RGAuthRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	res = &rpc.RGAuthRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
		Token:   calcToken,
	}
	return
}

func (s *RPCServer) CreateGroup(ctx context.Context, in *rpc.RGCreateGroupReq) (res *rpc.RGCreateGroupRes, err error) {
	return
}

func (s *RPCServer) Online(ctx context.Context, in *rpc.RGOnlineReq) (res *rpc.RGOnlineRes, err error) {
	glog.Info("Online")
	if _, err = s.dao.GetOnline(ctx, in.UID); err != nil {
		glog.Error(err)
		// not found
		err = nil
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
	glog.Info("Ping")
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
	glog.Info("[register] rpc server init: ", conf.Conf.RPCServer.Addr)
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
	rpcServer := &RPCServer{
		dao: dao.NewDao(),
	}
	rpc.RegisterRegisterServerRPCServer(s, rpcServer)
	s.Serve(lis)
}
