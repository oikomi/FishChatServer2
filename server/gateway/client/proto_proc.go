package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol"
	"github.com/oikomi/FishChatServer2/server/gateway/job"
)

func (c *Client) procReqMsgServer(reqData []byte) (err error) {
	var addr string
	msgServerList := job.MsgServerList
	for _, v := range msgServerList {
		glog.Info(v.IP)
		addr = v.IP
	}
	err = c.Session.Send(&protocol.SelectMsgServerForClient{
		Cmd:  protocol.SelectMsgServerForClientCMD,
		Addr: addr,
	})
	return
}
