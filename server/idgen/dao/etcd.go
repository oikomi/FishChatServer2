package dao

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
	"time"
)

type Etcd struct {
	EtcCli   *clientv3.Client
	rootPath string
}

func NewEtcd(c *commconf.Etcd) (etcd *Etcd, err error) {
	var etcdClient *clientv3.Client
	cfg := clientv3.Config{
		Endpoints:   c.Addrs,
		DialTimeout: time.Duration(c.Timeout),
	}
	if etcdClient, err = clientv3.New(cfg); err != nil {
		glog.Error("Error: cannot connec to etcd:", err)
		return
	}
	etcd = &Etcd{
		EtcCli:   etcdClient,
		rootPath: c.Root,
	}
	return
}
