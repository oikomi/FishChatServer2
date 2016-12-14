package etcd

import (
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"runtime"
	"time"
)

type Worker struct {
	Name     string
	IP       string
	rootPath string
	etcCli   *clientv3.Client
}

// workerInfo is the service register information to etcd
type WorkerInfo struct {
	Name string
	IP   string
	CPU  int
}

func NewWorker(name, ip, rootPath string, endpoints []string) *Worker {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		glog.Error("Error: cannot connec to etcd:", err)
	}
	w := &Worker{
		Name:     name,
		IP:       ip,
		rootPath: rootPath,
		etcCli:   etcdClient,
	}
	return w
}

func (w *Worker) HeartBeat() {
	for {
		info := &WorkerInfo{
			Name: w.Name,
			IP:   w.IP,
			CPU:  runtime.NumCPU(),
		}
		key := w.rootPath + w.Name
		value, err := json.Marshal(info)
		if err != nil {
			glog.Error(err)
		}
		resp, err := w.etcCli.Grant(context.TODO(), 10)
		if err != nil {
			glog.Error(err)
		}
		_, err = w.etcCli.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
		if err != nil {
			glog.Error("Error: cannot put to etcd:", err)
		}
		time.Sleep(time.Second * 5)
	}
}
