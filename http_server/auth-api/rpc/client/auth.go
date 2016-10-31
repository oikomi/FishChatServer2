package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/auth-api/conf"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type AuthRPCCli struct {
	conn *grpc.ClientConn
}

func NewAuthRPCCli() (authRPCCli *AuthRPCCli, err error) {
	conn, err := grpc.Dial(conf.Conf.RPCClient.AuthClient.Addr, grpc.WithInsecure())
	if err != nil {
		glog.Error(err)
		return
	}
	authRPCCli = &AuthRPCCli{
		conn: conn,
	}
	return
}

func (authRPCCli *AuthRPCCli) Login(loginReq *rpc.AuthLoginReq) (res *rpc.AuthLoginRes, err error) {
	m := rpc.NewAuthServerRPCClient(authRPCCli.conn)
	if res, err = m.Login(context.Background(), loginReq); err != nil {
		glog.Error(err)
	}
	return
}
