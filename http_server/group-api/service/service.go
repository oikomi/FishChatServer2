package service

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/http_server/group-api/model"
	groupRpc "github.com/oikomi/FishChatServer2/http_server/group-api/rpc"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
)

type Service struct {
	rpcClient *groupRpc.RPCClient
}

func New() (service *Service, err error) {
	rpcClient, err := groupRpc.NewRPCClient()
	if err != nil {
		glog.Error(err)
		return
	}
	service = &Service{
		rpcClient: rpcClient,
	}
	return
}

func (s *Service) CreateGroup(uid int64, groupName string) (err error) {
	rgCreateGroupReq := &rpc.RGCreateGroupReq{
		Uid:       uid,
		GroupName: groupName,
	}
	_, err = s.rpcClient.Register.CreateGroup(rgCreateGroupReq)
	if err != nil {
		glog.Error(err)
	}
	return
}

func (s *Service) JoinGroup(uid, gid int64) (err error) {
	rgJoinGroupReq := &rpc.RGJoinGroupReq{
		Uid: uid,
		Gid: gid,
	}
	_, err = s.rpcClient.Register.JoinGroup(rgJoinGroupReq)
	if err != nil {
		glog.Error(err)
	}
	return
}

func (s *Service) QuitGroup(uid, gid int64) (err error) {
	rgQuitGroupReq := &rpc.RGQuitGroupReq{
		Uid: uid,
		Gid: gid,
	}
	_, err = s.rpcClient.Register.QuitGroup(rgQuitGroupReq)
	if err != nil {
		glog.Error(err)
	}
	return
}
