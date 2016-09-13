package server

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol"
	"github.com/oikomi/FishChatServer2/server/router/client"
	"github.com/oikomi/FishChatServer2/server/router/conf"
	"github.com/oikomi/FishChatServer2/service_discovery/etcd"
)

type Server struct {
	Config *conf.Config
	Server *libnet.Server
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

func (s *Server) SDHeart() {
	work := etcd.NewWorker(conf.Conf.Etcd.Name, conf.Conf.Server.Addr, conf.Conf.Etcd.Root, conf.Conf.Etcd.Addrs)
	go work.HeartBeat()
}
