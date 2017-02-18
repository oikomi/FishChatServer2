package rpc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	commmodel "github.com/oikomi/FishChatServer2/common/model"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/manager/conf"
	"github.com/oikomi/FishChatServer2/server/manager/dao"
	"github.com/oikomi/FishChatServer2/server/manager/model"
	sd "github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type RPCServer struct {
	dao *dao.Dao
}

// func (s *RPCServer) GetOfflineMsgs(ctx context.Context, in *rpc.MGOfflineMsgReq) (res *rpc.MGOfflineMsgRes, err error) {
// 	glog.Info("GetOfflineMsgs")
// 	tmpRes, err := s.dao.MongoDB.GetOfflineMsg(in.Uid)
// 	if err != nil {
// 		glog.Error(err)
// 		res = &rpc.MGOfflineMsgRes{
// 			ErrCode: ecode.ServerErr.Uint32(),
// 			ErrStr:  ecode.ServerErr.String(),
// 		}
// 		return
// 	}
// 	resMsgs := make([]*rpc.OfflineMsg, 0)
// 	for _, msg := range tmpRes {
// 		resMsg := &rpc.OfflineMsg{
// 			SourceUID: msg.SourceUID,
// 			TargetUID: msg.TargetUID,
// 			MsgID:     msg.MsgID,
// 			Msg:       msg.Msg,
// 		}
// 		resMsgs = append(resMsgs, resMsg)
// 	}
// 	res = &rpc.MGOfflineMsgRes{
// 		ErrCode: ecode.OK.Uint32(),
// 		ErrStr:  ecode.OK.String(),
// 		Msgs:    resMsgs,
// 	}
// 	return
// }

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
	if err = s.dao.Redis.SetExceptionMsg(ctx, in.MsgID, string(data)); err != nil {
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

func (s *RPCServer) Sync(ctx context.Context, in *rpc.MGSyncMsgReq) (res *rpc.MGSyncMsgRes, err error) {
	glog.Info("Sync")
	offsetMsgs := make([]*rpc.MGSyncMsgResOffsetMsg, 0)
	_, err = s.dao.Mysql.UpdateUserMsgID(ctx, in.UID, in.CurrentID)
	if err != nil {
		glog.Error(err)
		res = &rpc.MGSyncMsgRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	userMsgID, err := s.dao.Mysql.GetUserMsgID(ctx, in.UID)
	if err != nil {
		glog.Error(err)
		res = &rpc.MGSyncMsgRes{
			ErrCode: ecode.ServerErr.Uint32(),
			ErrStr:  ecode.ServerErr.String(),
		}
		return
	}
	if userMsgID.CurrentMsgID == userMsgID.TotalMsgID {
		res = &rpc.MGSyncMsgRes{
			ErrCode:   ecode.OK.Uint32(),
			ErrStr:    ecode.OK.String(),
			CurrentID: userMsgID.TotalMsgID,
			Msgs:      offsetMsgs,
		}
		return
	}
	for i := userMsgID.CurrentMsgID; i <= userMsgID.TotalMsgID; i++ {
		hRes, err := s.dao.HBase.GetMsgs(ctx, fmt.Sprintf("%d_%d", in.UID, i))
		if err != nil {
			glog.Error(err)
		}
		offsetMsg := &rpc.MGSyncMsgResOffsetMsg{}
		for _, c := range hRes.Cells {
			if c != nil {
				if bytes.Equal(c.Family, model.HbaseFamilyUser) {
					if bytes.Equal(c.Qualifier, model.HbaseColumnSourceUID) {
						offsetMsg.SourceUID = int64(binary.BigEndian.Uint64(c.Value))
					} else if bytes.Equal(c.Qualifier, model.HbaseColumnTargetUID) {
						offsetMsg.TargetUID = int64(binary.BigEndian.Uint64(c.Value))
					} else if bytes.Equal(c.Qualifier, model.HbaseColumnGroupID) {
						offsetMsg.GroupID = int64(binary.BigEndian.Uint64(c.Value))
					}
				} else if bytes.Equal(c.Family, model.HbaseFamilyMsg) {
					if bytes.Equal(c.Qualifier, model.HbaseColumnMsgType) {
						offsetMsg.MsgType = string(c.Value)
					} else if bytes.Equal(c.Qualifier, model.HbaseColumnMsgID) {
						offsetMsg.MsgID = string(c.Value)
					} else if bytes.Equal(c.Qualifier, model.HbaseColumnMsg) {
						offsetMsg.Msg = string(c.Value)
					}
				}
			}
		}
		offsetMsgs = append(offsetMsgs, offsetMsg)
	}
	// for _, o := range offsetMsgs {
	// 	glog.Info(o)
	// }
	res = &rpc.MGSyncMsgRes{
		ErrCode:   ecode.OK.Uint32(),
		ErrStr:    ecode.OK.String(),
		CurrentID: userMsgID.TotalMsgID,
		Msgs:      offsetMsgs,
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
	err = sd.Register(conf.Conf.ServiceDiscoveryServer.ServiceName, conf.Conf.ServiceDiscoveryServer.RPCAddr, conf.Conf.ServiceDiscoveryServer.EtcdAddr, conf.Conf.ServiceDiscoveryServer.Interval, conf.Conf.ServiceDiscoveryServer.TTL)
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
