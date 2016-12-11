package service

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/group-api/model"
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

func (s *Service) CreateGroup(uid int64, groupID int64) (loginModel *model.Login, err error) {
	rgCreateGroupReq := &rpc.RGCreateGroupReq{
		UID: uid,
	}
	res, err := s.rpcClient.Register.CreateGroup(rgCreateGroupReq)
	if err != nil {
		glog.Error(err)
	}
	return
}
