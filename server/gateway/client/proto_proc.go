package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/server/gateway/job"
	"math/rand"
)

func (c *Client) procReqAccessServer(reqData []byte) (err error) {
	var addr string
	var accessServerList []string
	for _, v := range job.AccessServerList {
		accessServerList = append(accessServerList, v.IP)
	}
	if len(accessServerList) == 0 {
		if err = c.Session.Send(&external.ResSelectAccessServerForClient{
			Cmd:     external.ReqAccessServerCMD,
			ErrCode: ecode.NoAccessServer.Uint32(),
			ErrStr:  ecode.NoAccessServer.String(),
		}); err != nil {
			glog.Error(err)
		}
		return
	}
	addr = accessServerList[rand.Intn(len(accessServerList))]
	if err = c.Session.Send(&external.ResSelectAccessServerForClient{
		Cmd:     external.ReqAccessServerCMD,
		ErrCode: ecode.OK.Uint32(),
		ErrStr:  ecode.OK.String(),
		Addr:    addr,
	}); err != nil {
		glog.Error(err)
	}
	return
}
