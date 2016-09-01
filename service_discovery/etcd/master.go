package etcd

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	"time"
)

type Master struct {
	members  map[string]*Member
	etcCli   *clientv3.Client
	rootPath string
}

// Member is a client machine
type Member struct {
	InGroup bool
	IP      string
	Name    string
	CPU     int
}

func NewMaster(rootPath string, endpoints []string) *Master {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}
	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		glog.Error("Error: cannot connec to etcd:", err)
	}
	master := &Master{
		members:  make(map[string]*Member),
		etcCli:   etcdClient,
		rootPath: rootPath,
	}
	go master.WatchWorkers()
	return master
}

func (m *Master) Members() (ms map[string]*Member) {
	ms = m.members
	return
}

func (m *Master) AddWorker(info *WorkerInfo) {
	member := &Member{
		InGroup: true,
		IP:      info.IP,
		Name:    info.Name,
		CPU:     info.CPU,
	}
	m.members[member.Name] = member
}

func (m *Master) UpdateWorker(info *WorkerInfo) {
	member := m.members[info.Name]
	member.InGroup = true
}

func (m *Master) WatchWorkers() {
	for {
		rch := m.etcCli.Watch(context.Background(), m.rootPath, clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				glog.Info(ev.Type.String())
				//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				if ev.Type.String() == "EXPIRE" {
					member, ok := m.members[string(ev.Kv.Key)]
					if ok {
						member.InGroup = false
					}
				} else if ev.Type.String() == "PUT" {
					info := &WorkerInfo{}
					err := json.Unmarshal(ev.Kv.Value, info)
					if err != nil {
						glog.Error(err)
					}
					if _, ok := m.members[info.Name]; ok {
						m.UpdateWorker(info)
					} else {
						m.AddWorker(info)
					}
				} else if ev.Type.String() == "DELETE" {
					delete(m.members, string(ev.Kv.Key))
				}

			}
		}
		glog.Info("end WatchWorkers")
	}
}
