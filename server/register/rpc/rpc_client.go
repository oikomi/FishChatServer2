package rpc

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/register/rpc/client"
)

type RPCClient struct {
	Idgen *client.IdgenRPCCli
}

func NewRPCClient() (c *RPCClient, err error) {
	idgen, err := client.NewIdgenRPCCli()
	if err != nil {
		glog.Error(err)
		return
	}
	c = &RPCClient{
		Idgen: idgen,
	}
	return
}
