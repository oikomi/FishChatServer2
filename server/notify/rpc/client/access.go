package client

import (
	"errors"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/protocol/rpc"
	"github.com/oikomi/FishChatServer2/server/notify/conf_discovery"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

type AccessServerRPCCli struct {
	conns map[string]*grpc.ClientConn
}

func NewAccessServerRPCCli() (accessServerRPCCli *AccessServerRPCCli, err error) {
	glog.Info("NewAccessServerRPCCli")
	var accessServerList []string
	conns := make(map[string]*grpc.ClientConn)
	for {
		if len(conf_discovery.AccessServerList) <= 0 {
			time.Sleep(time.Second * 5)
		} else {
			glog.Info(conf_discovery.AccessServerList)
			for _, v := range conf_discovery.AccessServerList {
				accessServerList = append(accessServerList, v.IP)
			}
			for _, accessServer := range accessServerList {
				conn, err := grpc.Dial(accessServer, grpc.WithInsecure())
				if err != nil {
					glog.Error(err)
				}
				conns[accessServer] = conn
			}
			accessServerRPCCli = &AccessServerRPCCli{
				conns: conns,
			}
			break
		}
	}
	go accessServerRPCCli.connProc()
	return
}

func (accessServerRPCCli *AccessServerRPCCli) connProc() {
	glog.Info(conf_discovery.AccessServerList)
	var accessServerList []string
	conns := make(map[string]*grpc.ClientConn)
	for {
		for _, v := range conf_discovery.AccessServerList {
			glog.Info(v.IP)
			accessServerList = append(accessServerList, v.IP)
		}
		for _, accessServer := range accessServerList {
			conn, err := grpc.Dial(accessServer, grpc.WithInsecure())
			if err != nil {
				glog.Error(err)
			}
			conns[accessServer] = conn
		}
		accessServerRPCCli.conns = conns
		time.Sleep(time.Second * 10)
	}
}

// FIXME can not use rr
func (accessServerRPCCli *AccessServerRPCCli) SendNotify(ctx context.Context, sendNotifyReq *rpc.ASSendNotifyReq) (res *rpc.ASSendNotifyRes, err error) {
	glog.Info(accessServerRPCCli.conns)
	glog.Info(sendNotifyReq.AccessServerAddr)
	if conn, ok := accessServerRPCCli.conns[sendNotifyReq.AccessServerAddr]; ok {
		a := rpc.NewAccessServerRPCClient(conn)
		if res, err = a.SendNotify(ctx, sendNotifyReq); err != nil {
			glog.Error(err)
		}
	} else {
		err = errors.New("no access server")
	}
	return
}
