package etcd

import (
	"fmt"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	. "github.com/oikomi/FishChatServer2/service_discovery/lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

// prefix is the root Dir of services in etcd
var Prefix = "im"

// EtcdWatcher is the implementaion of grpc.naming.Watcher
type EtcdWatcher struct {
	// er: EtcdResolver
	er *EtcdResolver
	// ec: Etcd Client
	ec *etcd.Client
	// addrs is the service addrs cache
	addrs         []string
	isInitialized bool
}

// Close do nothing
func (ew *EtcdWatcher) Close() {
}

// Next to return the updates
func (ew *EtcdWatcher) Next() ([]*naming.Update, error) {
	// key is the etcd key/value dir to watch
	key := fmt.Sprintf("/%s/%s", Prefix, ew.er.ServiceName)
	// ew.addrs is nil means it is intially called
	if ew.addrs == nil {
		// query addresses from etcd
		resp, _ := ew.ec.Get(context.Background(), key, etcd.WithPrefix())
		addrs, empty := extractAddrs(resp)
		dropEmptyDir(ew.ec, empty)
		// addrs is not empty, return the updates
		// addrs is empty, should to watch new data
		if len(addrs) != 0 {
			ew.addrs = addrs
			return GenUpdates([]string{}, addrs), nil
		}
	}
	for {
		// generate etcd Watcher
		rch := ew.ec.Watch(context.Background(), key, etcd.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				// glog.Info(ev.Type, string(ev.Kv.Key), string(ev.Kv.Value))
				if ev.Type.String() == "EXPIRE" {
					return []*naming.Update{{Op: naming.Delete, Addr: string(ev.Kv.Value)}}, nil
				} else if ev.Type.String() == "PUT" {
					return []*naming.Update{{Op: naming.Add, Addr: string(ev.Kv.Value)}}, nil
				} else if ev.Type.String() == "DELETE" {
					return []*naming.Update{{Op: naming.Delete, Addr: string(ev.Kv.Value)}}, nil
				}
			}
		}
	}
}

// helper function to extract addrs rom etcd response
func extractAddrs(resp *etcd.GetResponse) (addrs, empty []string) {
	addrs = []string{}
	empty = []string{}
	if resp == nil || resp.Kvs == nil || len(resp.Kvs) == 0 {
		return addrs, empty
	}
	for _, v := range resp.Kvs {
		addr := ""
		what := v.Key[len(v.Key)-4 : len(v.Key)]
		if string(what) == "addr" {
			addr = string(v.Value)
		}
		if addr != "" {
			addrs = append(addrs, addr)
		}
	}
	glog.Info(addrs)
	return addrs, empty
}

func dropEmptyDir(keyapi *etcd.Client, empty []string) {
	if keyapi == nil || len(empty) == 0 {
		return
	}
	for _, key := range empty {
		_, err := keyapi.Delete(context.Background(), key)
		if err != nil {
			glog.Error("service_discovery: delete empty service dir error: ", err.Error())
		}
	}
}
