package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/server/access/rpc"
)

type Client struct {
	Session   *libnet.Session
	RPCClient *rpc.RPCClient
}

func New(session *libnet.Session, rpcClient *rpc.RPCClient) (c *Client) {
	c = &Client{
		Session:   session,
		RPCClient: rpcClient,
	}
	return
}

func (c *Client) Parse(cmd uint32, reqData []byte) (err error) {
	switch cmd {
	case external.LoginCMD:
		if err = c.procReqLogin(reqData); err != nil {
			glog.Error(err)
			return
		}
	case external.PingCMD:
		if err = c.procReqPing(reqData); err != nil {
			glog.Error(err)
			return
		}
	case external.SendP2PMsgCMD:
		if err = c.procSendP2PMsg(reqData); err != nil {
			glog.Error(err)
			return
		}
	// case external.AcceptP2PMsgAckCMD:
	// 	if err = c.procAcceptP2PMsgAck(reqData); err != nil {
	// 		glog.Error(err)
	// 		return
	// 	}
	case external.SendGroupMsgCMD:
		if err = c.procSendGroupMsg(reqData); err != nil {
			glog.Error(err)
			return
		}
	case external.SyncMsgCMD:
		if err = c.procSyncMsg(reqData); err != nil {
			glog.Error(err)
			return
		}
	}
	return
}
