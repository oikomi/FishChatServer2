package rpc

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/common/ecode"
	"encoding/json"
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
	exceptionMsg := &commmodel.ExceptionMsg{
		SourceUID: in.SourceUID,
		TargetUID: in.TargetUID,
		MsgID:     in.MsgID,
		Msg:       in.Msg,
	}
	data, err := json.Marshal(exceptionMsg)
	if err != nil {
		glog.Error(err)
		return
	}
	if err = s.dao.SetExceptionMsg(ctx, in.MsgID, string(data)); err != nil {
		glog.Error(err)
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
