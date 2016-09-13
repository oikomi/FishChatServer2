package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol"
	"github.com/oikomi/FishChatServer2/server/msg_server/rpc"
)

type Client struct {
	Session *libnet.Session
	rpcCli  *rpc.RPCClient
}

func New(session *libnet.Session, rpcCli *rpc.RPCClient) (c *Client) {
	c = &Client{
		Session: session,
		rpcCli:  rpcCli,
	}
	return
}

func (c *Client) Parse(cmd uint32, reqData []byte) (err error) {
	switch cmd {
	case protocol.ReqLoginCMD:
		if err = c.procLogin(reqData); err != nil {
			glog.Error(err)
			return
		}
	}
	return
}
