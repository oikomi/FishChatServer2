package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/msg_server/rpc/client"
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

// func (rc *RPCClient) init() {

// }
