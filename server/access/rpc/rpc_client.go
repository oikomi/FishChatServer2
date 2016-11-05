package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/access/rpc/client"
)

type RPCClient struct {
	Logic *client.LogicRPCCli
}

func NewRPCClient() (c *RPCClient, err error) {
	logic, err := client.NewLogicRPCCli()
	if err != nil {
		glog.Error(err)
		return
	}
	c = &RPCClient{
		Logic: logic,
	}
	return
}
