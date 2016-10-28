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

func loadAccessServerProc(master *etcd.Master) {
	for {
		// glog.Info("loadAccessServerProc")
		AccessServerList = master.Members()
		time.Sleep(time.Duration(conf.Conf.ServiceDiscovery.Interval))
	}
}

func ServerDiscoveryProc() {
	master, err := etcd.NewMaster(conf.Conf.Etcd)
	if err != nil {
		glog.Error("Error: cannot connect to etcd:", err)
		panic(err)
	}
	go loadAccessServerProc(master)
	master.WatchWorkers()
}
