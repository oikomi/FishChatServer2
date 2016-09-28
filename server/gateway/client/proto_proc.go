package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/server/gateway/job"
)

func (c *Client) procReqAccessServer(reqData []byte) (err error) {
	var addr string
	msgServerList := job.AccessServerList
	if len(msgServerList) == 0 {
		if err = c.Session.Send(&external.ResSelectAccessServerForClient{
			Cmd:     external.ResSelectAccessServerForClientCMD,
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
	if err = c.Session.Send(&external.ResSelectAccessServerForClient{
		Cmd:     external.ResSelectAccessServerForClientCMD,
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
		Addr:    addr,
	}); err != nil {
		glog.Error(err)
	}
	return
}
