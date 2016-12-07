package rpc

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	commmodel "github.com/oikomi/FishChatServer2/common/model"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/manager/conf"
	"github.com/oikomi/FishChatServer2/server/manager/dao"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	dao *dao.Dao
}

func (s *RPCServer) ExceptionMsg(ctx context.Context, in *rpc.MGExceptionMsgReq) (res *rpc.MGExceptionMsgRes, err error) {
	glog.Info("ExceptionMsg")
	return
}

func (s *RPCServer) SetExceptionMsg(ctx context.Context, in *rpc.MGExceptionMsgReq) (res *rpc.MGExceptionMsgRes, err error) {
	glog.Info("SetExceptionMsg")
	exceptionMsg := &commmodel.ExceptionMsg{
		SourceUID: in.SourceUID,
		TargetUID: in.TargetUID,
		MsgID:     in.MsgID,
		Msg:       in.Msg,
	}
	data, err := json.Marshal(exceptionMsg)
	if err != nil {
		glog.Error(err)
		res = &rpc.MGExceptionMsgRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	if err = s.dao.SetExceptionMsg(ctx, in.MsgID, string(data)); err != nil {
		glog.Error(err)
		res = &rpc.MGExceptionMsgRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	res = &rpc.MGExceptionMsgRes{
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
	}
	return
}

func RPCServerInit() {
	glog.Info("[manager] rpc server init: ", conf.Conf.RPCServer.Addr)
	lis, err := net.Listen(conf.Conf.RPCServer.Proto, conf.Conf.RPCServer.Addr)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	s := grpc.NewServer()
	rpcServer := &RPCServer{
		dao: dao.NewDao(),
	}
	rpc.RegisterManagerServerRPCServer(s, rpcServer)
	s.Serve(lis)
}
