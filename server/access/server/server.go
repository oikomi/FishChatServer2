package server

import (
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/conf_discovery/etcd"
	"github.com/oikomi/FishChatServer2/libnet"
	"github.com/oikomi/FishChatServer2/protocol/external"
	"github.com/oikomi/FishChatServer2/server/access/client"
	"github.com/oikomi/FishChatServer2/server/access/conf"
	"github.com/oikomi/FishChatServer2/server/access/rpc"
)

type Server struct {
	Server    *libnet.Server
	RPCServer *rpc.RPCServer
}

func New() (s *Server) {
	s = &Server{}
	go rpc.RPCServerInit()
	return
}

func (s *Server) sessionLoop(client *client.Client) {
	for {
		reqData, err := client.Session.Receive()
		if err != nil {
			glog.Error(err)
		}
		if reqData != nil {
			baseCMD := &external.Base{}
			if err = proto.Unmarshal(reqData, baseCMD); err != nil {
				if err = client.Session.Send(&external.Error{
					Cmd:     external.ErrServerCMD,
					ErrCode: ecode.ServerErr.Uint32(),
					ErrStr:  ecode.ServerErr.String(),
				}); err != nil {
					glog.Error(err)
				}
				continue
			}
			if err = client.Parse(baseCMD.Cmd, reqData); err != nil {
				glog.Error(err)
				continue
			}
		}
	}
}

func (s *Server) Loop(rpcClient *rpc.RPCClient) {
	for {
		session, err := s.Server.Accept()
		if err != nil {
			glog.Error(err)
		}
		glog.Info("a new client ", session.ID())
		c := client.New(session, rpcClient)
		go s.sessionLoop(c)
	}
}

func (s *Server) SDHeart() {
	glog.Info("SDHeart")
	work1 := etcd.NewWorker(conf.Conf.ConfDiscovery.Gateway.Name, conf.Conf.Server.Addr, conf.Conf.ConfDiscovery.Gateway.Root,
		conf.Conf.ConfDiscovery.Gateway.Addrs)
	go work1.HeartBeat()
	work2 := etcd.NewWorker(conf.Conf.ConfDiscovery.MsgJob.Name, conf.Conf.RPCServer.Addr, conf.Conf.ConfDiscovery.MsgJob.Root,
		conf.Conf.ConfDiscovery.MsgJob.Addrs)
	go work2.HeartBeat()
	glog.Info(conf.Conf.ConfDiscovery.Notify)
	work3 := etcd.NewWorker(conf.Conf.ConfDiscovery.Notify.Name, conf.Conf.RPCServer.Addr, conf.Conf.ConfDiscovery.Notify.Root,
		conf.Conf.ConfDiscovery.Notify.Addrs)
	go work3.HeartBeat()
	work4 := etcd.NewWorker(conf.Conf.ConfDiscovery.Logic.Name, conf.Conf.RPCServer.Addr, conf.Conf.ConfDiscovery.Logic.Root,
		conf.Conf.ConfDiscovery.Logic.Addrs)
	go work4.HeartBeat()
}
