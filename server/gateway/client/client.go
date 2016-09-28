package client

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol/external"
)

type Client struct {
	Session *libnet.Session
}

func New(session *libnet.Session) (c *Client) {
	c = &Client{
		Session: session,
	}
	return
}

func (c *Client) Parse(cmd uint32, reqData []byte) (err error) {
	switch cmd {
	case external.ReqAccessServerCMD:
		if err = c.procReqAccessServer(reqData); err != nil {
			glog.Error(err)
			return
		}
	}
	return
}
