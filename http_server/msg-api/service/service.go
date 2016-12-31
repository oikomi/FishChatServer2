package service

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/msg-api/model"
	authRpc "github.com/oikomi/FishChatServer2/http_server/msg-api/rpc"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
)

type Service struct {
	rpcClient *authRpc.RPCClient
}

func New() (service *Service, err error) {
	rpcClient, err := authRpc.NewRPCClient()
	if err != nil {
		glog.Error(err)
		return
	}
	service = &Service{
		rpcClient: rpcClient,
	}
	return
}

func (s *Service) GetOfflineMsgs(uid int64, pw string) (loginModel *model.Login, err error) {
	// check uid pw
	rgAuthReq := &rpc.RGAuthReq{
		UID: uid,
	}
	res, err := s.rpcClient.Register.Auth(rgAuthReq)
	if err != nil {
		glog.Error(err)
		return
	}
	loginModel = &model.Login{
		Token: res.Token,
	}
	return
}
