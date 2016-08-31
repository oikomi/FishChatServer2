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
	glog.Info(endpoints)
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
	// var err error
	for {
		info := &WorkerInfo{
			Name: w.Name,
			IP:   w.IP,
			CPU:  runtime.NumCPU(),
		}
		key := w.rootPath + "/" + w.Name
		value, _ := json.Marshal(info)
		// ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
		_, err := w.etcCli.Put(context.TODO(), key, string(value))
		// cancel()
		if err != nil {
			glog.Error("Error: cannot put to etcd:", err)
		}
		time.Sleep(time.Second * 10)

		// getResp, err := w.etcCli.Get(context.TODO(), key)
		// if err != nil {
		// 	glog.Error("Error: cannot put to etcd:", err)
		// }
		// glog.Info(getResp.Kvs)
	}
}
