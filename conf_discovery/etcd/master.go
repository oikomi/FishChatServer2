package etcd

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
	"time"
)

// Master is a server
type Master struct {
	members  map[string]*Member
	etcCli   *clientv3.Client
	rootPath string
}

// Member is a client
type Member struct {
	InGroup bool
	IP      string
	Name    string
	CPU     int
}

func NewMaster(c *commconf.Etcd) (master *Master, err error) {
	var etcdClient *clientv3.Client
	cfg := clientv3.Config{
		Endpoints:   c.Addrs,
		DialTimeout: time.Duration(c.Timeout),
	}
	if etcdClient, err = clientv3.New(cfg); err != nil {
		glog.Error("Error: cannot connec to etcd:", err)
		return
	}
	master = &Master{
		members:  make(map[string]*Member),
		etcCli:   etcdClient,
		rootPath: c.Root,
	}
	return
}

func (m *Master) Members() (ms map[string]*Member) {
	ms = m.members
	return
}

func (m *Master) addWorker(key string, info *WorkerInfo) {
	member := &Member{
		InGroup: true,
		IP:      info.IP,
		Name:    info.Name,
		CPU:     info.CPU,
	}
	m.members[key] = member
}

func (m *Master) updateWorker(key string, info *WorkerInfo) {
	member := m.members[key]
	member.InGroup = true
}

func (m *Master) WatchWorkers() {
	rch := m.etcCli.Watch(context.Background(), m.rootPath, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			// glog.Info(ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
			if ev.Type.String() == "EXPIRE" {
				member, ok := m.members[string(ev.Kv.Key)]
				if ok {
					member.InGroup = false
					delete(m.members, string(ev.Kv.Key))
				}
			} else if ev.Type.String() == "PUT" {
				info := &WorkerInfo{}
				err := json.Unmarshal(ev.Kv.Value, info)
				if err != nil {
					glog.Error(err)
				}
				if _, ok := m.members[string(ev.Kv.Key)]; ok {
					m.updateWorker(string(ev.Kv.Key), info)
				} else {
					m.addWorker(string(ev.Kv.Key), info)
				}
			} else if ev.Type.String() == "DELETE" {
				delete(m.members, string(ev.Kv.Key))
			}
		}
	}
}
