package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/user-api/rpc/client"
)

type RPCClient struct {
	Register *client.RegisterRPCCli
}

func NewRPCClient() (c *RPCClient, err error) {
	register, err := client.NewRegisterRPCCli()
	if err != nil {
		glog.Error(err)
		return
	}
	c = &RPCClient{
		Register: register,
	}
	return
}
