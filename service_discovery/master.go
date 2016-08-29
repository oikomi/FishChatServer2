package service_discovery

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

type Master struct {
	members map[string]*Member
	etcCli  *clientv3.Client
}

// Member is a client machine
type Member struct {
	InGroup bool
	IP      string
	Name    string
	CPU     int
}

func NewMaster(endpoints []string) *Master {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal("Error: cannot connec to etcd:", err)
	}

	master := &Master{
		members: make(map[string]*Member),
		etcCli:  etcdClient,
	}
	go master.WatchWorkers()
	return master
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
	rch := m.etcCli.Watch(context.Background(), "workers")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			//fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if ev.Type.String() == "expire" {
				member, ok := m.members[string(ev.Kv.Key)]
				if ok {
					member.InGroup = false
				}
			} else if ev.Type.String() == "set" || ev.Type.String() == "update" {
				info := &WorkerInfo{}
				err := json.Unmarshal(ev.Kv.Value, info)
				if err != nil {
					log.Print(err)
				}
				if _, ok := m.members[info.Name]; ok {
					m.UpdateWorker(info)
				} else {
					m.AddWorker(info)
				}
			} else if ev.Type.String() == "delete" {
				delete(m.members, string(ev.Kv.Key))
			}

		}

	}

	// watcher := m.etcCli.Watch("workers/", &client.WatcherOptions{
	// 	Recursive: true,
	// })
	// for {
	// 	res, err := watcher.Next(context.Background())
	// 	if err != nil {
	// 		log.Println("Error watch workers:", err)
	// 		break
	// 	}
	// 	if res.Action == "expire" {
	// 		member, ok := m.members[res.Node.Key]
	// 		if ok {
	// 			member.InGroup = false
	// 		}
	// 	} else if res.Action == "set" || res.Action == "update" {
	// 		info := &WorkerInfo{}
	// 		err := json.Unmarshal([]byte(res.Node.Value), info)
	// 		if err != nil {
	// 			log.Print(err)
	// 		}
	// 		if _, ok := m.members[info.Name]; ok {
	// 			m.UpdateWorker(info)
	// 		} else {
	// 			m.AddWorker(info)
	// 		}
	// 	} else if res.Action == "delete" {
	// 		delete(m.members, res.Node.Key)
	// 	}
	// }

}
