package client

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/protocol"
)

func (c *Client) procLogin(reqData []byte) (err error) {
	reqLoginPB := &protocol.ReqLogin{}
	if err = proto.Unmarshal(reqData, reqLoginPB); err != nil {
		glog.Error(err)
		return
	}
	c.rpcCli.Manager.Login(reqLoginPB.UID, reqLoginPB.Token)
	return
}

func (c *Client) procSendP2PMsg(reqData []byte) (err error) {
	reqSendP2PMsgPB := &protocol.ReqSendP2PMsg{}
	if err = proto.Unmarshal(reqData, reqSendP2PMsgPB); err != nil {
		glog.Error(err)
		return
	}
	c.rpcCli.Router.SendMsgP2P(reqSendP2PMsgPB.TargetUID, reqSendP2PMsgPB.Msg)
	return
}
