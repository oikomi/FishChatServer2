package job

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/gateway/conf"
	"github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"time"
)

var (
	MsgServerList map[string]*etcd.Member
)

func loadMsgServer(master *etcd.Master) {
	for {
		MsgServerList = master.Members()
		glog.Info(MsgServerList)
		time.Sleep(time.Second * 5)
	}
}

func DoServerDiscovery() {
	master, err := etcd.NewMaster(conf.Conf.Etcd.Root, conf.Conf.Etcd.Addrs)
	if err != nil {
		glog.Error("Error: cannot connec to etcd:", err)
		panic(err)
	}
	go loadMsgServer(master)
	master.WatchWorkers()
}
