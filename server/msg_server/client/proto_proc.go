package client

import (
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/protocol"
)

func (c *Client) procLogin(reqData []byte) (err error) {
	reqLoginPB := &protocol.ReqLogin{}
	proto.Unmarshal(reqData, reqLoginPB)

	return
}
