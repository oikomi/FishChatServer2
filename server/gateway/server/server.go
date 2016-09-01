package server

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol"
	"github.com/oikomi/FishChatServer2/server/gateway/client"
	"github.com/oikomi/FishChatServer2/server/gateway/conf"
	"github.com/oikomi/FishChatServer2/service_discovery/etcd"
)

type Server struct {
	Config        *conf.Config
	Server        *libnet.Server
	Master        *etcd.Master
	MsgServerList []*etcd.Member
}

func New(config *conf.Config) (s *Server) {
	s = &Server{
		Config: config,
	}
	return
}

func (s *Server) sessionLoop(client *client.Client) {
	for {
		reqData, err := client.Session.Receive()
		if err != nil {
			glog.Error(err)
		}
		if reqData != nil {
			baseCMD := &protocol.Base{}
			if err = proto.Unmarshal(reqData, baseCMD); err != nil {
				err = client.Session.Send(&protocol.Error{
					ErrCode: ecode.ServerErr.Uint32(),
					ErrStr:  ecode.ServerErr.String(),
				})
			}
			client.Parse(baseCMD.Cmd, reqData)
		}
	}
}

func (s *Server) Loop() {
	for {
		session, err := s.Server.Accept()
		if err != nil {
			glog.Error(err)
		}
		go s.sessionLoop(client.New(session))
	}
}
