package service

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/user-api/model"
	authRpc "github.com/oikomi/FishChatServer2/http_server/user-api/rpc"
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

func (s *Service) Auth(uid int64, pw string) (loginModel *model.Login, err error) {
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

func (s *Service) Register(uid int64, userName, pw string) (err error) {
	rgRegisterReq := &rpc.RGRegisterReq{
		UID:      uid,
		Name:     userName,
		Password: pw,
	}
	_, err = s.rpcClient.Register.Register(rgRegisterReq)
	if err != nil {
		glog.Error(err)
		return
	}
	return
}
