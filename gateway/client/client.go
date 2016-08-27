package client

import (
	"github.com/oikomi/FishChatServer2/libnet"
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
