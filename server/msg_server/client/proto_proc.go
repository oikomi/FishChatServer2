package client

import (
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/protocol"
)

func (c *Client) procSendClientID(reqData []byte) (err error) {
	sendClientIDPB := &protocol.SendClientID{}
	proto.Unmarshal(reqData, sendClientIDPB)
	if err = c.Session.Send(&protocol.SelectMsgServerForClient{
		Cmd:  protocol.SelectMsgServerForClientCMD,
		Addr: addr,
	}); err != nil {
		glog.Error(err)
	}
	return
}
