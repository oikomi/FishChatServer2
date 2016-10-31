package service

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/auth-api/model"
	authRpc "github.com/oikomi/FishChatServer2/http_server/auth-api/rpc"
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

func (s *Service) Login(uid int64, pw string) (loginModel *model.Login, err error) {
	// check uid pw
	authLoginReq := &rpc.AuthLoginReq{
		UID: uid,
	}
	res, err := s.rpcClient.Auth.Login(authLoginReq)
	if err != nil {
		glog.Error(err)
	}
	loginModel = &model.Login{
		Token: res.Token,
	}
	return
}
