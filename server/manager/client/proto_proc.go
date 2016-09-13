package client

import (
	"github.com/oikomi/FishChatServer2/protocol"
)

func (c *Client) procReqMsgServer(reqData []byte) (err error) {

	err = c.Session.Send(&protocol.SelectMsgServerForClient{
		Cmd: protocol.SelectMsgServerForClientCMD,
	})
	return
}
