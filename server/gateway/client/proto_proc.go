package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol"
	"github.com/oikomi/FishChatServer2/server/gateway/job"
)

func (c *Client) procReqMsgServer(reqData []byte) (err error) {
	var addr string
	msgServerList := job.MsgServerList
	if len(msgServerList) == 0 {
		if err = c.Session.Send(&protocol.ResSelectMsgServerForClient{
			Cmd:     protocol.ResSelectMsgServerForClientCMD,
			ErrCode: ecode.NoMsgServer.Uint32(),
			ErrStr:  ecode.NoMsgServer.String(),
		}); err != nil {
			glog.Error(err)
		}
		return
	}
	for _, v := range msgServerList {
		addr = v.IP
	}
	if err = c.Session.Send(&protocol.ResSelectMsgServerForClient{
		Cmd:     protocol.ResSelectMsgServerForClientCMD,
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
		Addr:    addr,
	}); err != nil {
		glog.Error(err)
	}
	return
}
