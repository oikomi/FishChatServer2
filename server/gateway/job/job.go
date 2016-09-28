package job

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/server/gateway/conf"
	"github.com/oikomi/FishChatServer2/service_discovery/etcd"
	"time"
)

var (
	AccessServerList map[string]*etcd.Member
)

func loadAccessServer(master *etcd.Master) {
	for {
		AccessServerList = master.Members()
		glog.Info(AccessServerList)
		time.Sleep(time.Second * 5)
	}
}

func DoServerDiscovery() {
	master, err := etcd.NewMaster(conf.Conf.Etcd.Root, conf.Conf.Etcd.Addrs)
	if err != nil {
		glog.Error("Error: cannot connec to etcd:", err)
		panic(err)
	}
	go loadAccessServer(master)
	master.WatchWorkers()
}
