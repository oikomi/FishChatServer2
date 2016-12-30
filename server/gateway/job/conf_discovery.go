package job

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/conf_discovery/etcd"
	"github.com/oikomi/FishChatServer2/server/gateway/conf"
	"time"
)

var (
	AccessServerList map[string]*etcd.Member
)

func loadAccessServerProc(master *etcd.Master) {
	for {
		// glog.Info("loadAccessServerProc")
		AccessServerList = master.Members()
		time.Sleep(time.Second * 5)
	}
}

func ConfDiscoveryProc() {
	master, err := etcd.NewMaster(conf.Conf.Etcd)
	if err != nil {
		glog.Error("Error: cannot connect to etcd:", err)
		panic(err)
	}
	go loadAccessServerProc(master)
	master.WatchWorkers()
}
