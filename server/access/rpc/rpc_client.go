package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/access/rpc/client"
)

type RPCClient struct {
	MsgServer *client.MsgServerRPCCli
}

func NewRPCClient() (c *RPCClient, err error) {
	msgServer, err := client.NewMsgServerRPCCli()
	if err != nil {
		glog.Error(err)
		return
	}
	c = &RPCClient{
		MsgServer: msgServer,
	}
	return
}
