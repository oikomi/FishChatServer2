package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/msg-api/rpc/client"
)

type RPCClient struct {
	Manager *client.ManagerRPCCli
}

func NewRPCClient() (c *RPCClient, err error) {
	manager, err := client.NewManagerRPCCli()
	if err != nil {
		glog.Error(err)
		return
	}
	c = &RPCClient{
		Manager: manager,
	}
	return
}
