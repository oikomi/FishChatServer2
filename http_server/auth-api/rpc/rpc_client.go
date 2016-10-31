package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/auth-api/rpc/client"
)

type RPCClient struct {
	Auth *client.MsgServerRPCCli
}

func NewRPCClient() (c *RPCClient, err error) {
	auth, err := client.NewAuthRPCCli()
	if err != nil {
		glog.Error(err)
		return
	}
	c = &RPCClient{
		Auth: auth,
	}
	return
}
