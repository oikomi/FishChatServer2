package server

import (
	"github.com/oikomi/FishChatServer2/server/gateway/client"
	// "github.com/oikomi/FishChatServer2/server/gateway/server"
	"github.com/oikomi/FishChatServer2/server/gateway/conf"
)

type ProtoProc struct {
	Config *conf.Config
}

func NewProtoProc(config *conf.Config) (pp *ProtoProc) {
	pp = &ProtoProc{
		Config: config,
	}
	return
}

func (pp *ProtoProc) Parse(cmd uint32, reqData []byte, client *client.Client) {

}
