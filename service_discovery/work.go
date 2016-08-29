package service_discovery

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"log"
	"runtime"
	"time"
)

type Worker struct {
	Name   string
	IP     string
	etcCli *clientv3.Client
}

// workerInfo is the service register information to etcd
type WorkerInfo struct {
	Name string
	IP   string
	CPU  int
}

func NewWorker(name, IP string, endpoints []string) *Worker {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	w := &Worker{
		Name:   name,
		IP:     IP,
		etcCli: etcdClient,
	}
	go w.HeartBeat()
	return w
}

func (w *Worker) HeartBeat() {
	var err error
	for {
		info := &WorkerInfo{
			Name: w.Name,
			IP:   w.IP,
			CPU:  runtime.NumCPU(),
		}
		key := "workers/" + w.Name
		value, _ := json.Marshal(info)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err = w.etcCli.Put(ctx, key, string(value))
		cancel()
		// _, err := api.Set(context.Background(), key, string(value), &client.SetOptions{
		// 	TTL: time.Second * 10,
		// })
		if err != nil {
			log.Println("Error update workerInfo:", err)
		}
		time.Sleep(time.Second * 3)
	}
}
