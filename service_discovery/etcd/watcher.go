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
var Prefix = "wonaming"

// EtcdWatcher is the implementaion of grpc.naming.Watcher
type EtcdWatcher struct {
	// er: EtcdResolver
	er *EtcdResolver
	// ec: Etcd Client
	ec *etcd.Client
	// addrs is the service addrs cache
	addrs []string
}

// Close do nothing
func (ew *EtcdWatcher) Close() {
}

// Next to return the updates
func (ew *EtcdWatcher) Next() ([]*naming.Update, error) {
	// key is the etcd key/value dir to watch
	key := fmt.Sprintf("/%s/%s", Prefix, ew.er.ServiceName)
	// keyapi := etcd.NewKeysAPI(*ew.ec)
	// ew.addrs is nil means it is intially called
	if ew.addrs == nil {
		// query addresses from etcd
		// resp, _ := keyapi.Get(context.Background(), key, &etcd.GetOptions{Recursive: true})
		resp, _ := ew.ec.Get(context.Background(), key)
		addrs, empty := extractAddrs(resp)
		dropEmptyDir(ew.ec, empty)
		// addrs is not empty, return the updates
		// addrs is empty, should to watch new data
		if len(addrs) != 0 {
			ew.addrs = addrs
			return GenUpdates([]string{}, addrs), nil
		}
	}
	// generate etcd Watcher
	// w := keyapi.Watcher(key, &etcd.WatcherOptions{Recursive: true})
	rch := ew.ec.Watch(context.Background(), key)
	for _ = range rch {
		// _, err := w.Next(context.Background())
		// if err == nil {
		// query addresses from etcd
		resp, err := ew.ec.Get(context.Background(), key)
		if err != nil {
			continue
		}
		addrs, empty := extractAddrs(resp)
		dropEmptyDir(ew.ec, empty)

		updates := GenUpdates(ew.addrs, addrs)
		// update ew.addrs
		ew.addrs = addrs
		// if addrs updated, return it
		if len(updates) != 0 {
			return updates, nil
		}
		// }
	}
	// should not goto here for ever
	return []*naming.Update{}, nil
}

// helper function to extract addrs rom etcd response
func extractAddrs(resp *etcd.GetResponse) (addrs, empty []string) {
	addrs = []string{}
	empty = []string{}
	if resp == nil || resp.Kvs == nil || len(resp.Kvs) == 0 {
		return addrs, empty
	}
	for _, v := range resp.Kvs {
		// prevValue, err = strconv.Atoi(string(value.Value))
		// if err != nil {
		// 	glog.Error(err)
		// 	return
		// }
		// prevIndex = value.ModRevision
		// glog.Info(prevValue)
		// glog.Info(prevIndex)

		// node should contain host & port both
		host := ""
		port := ""
		what := v.Key[len(v.Key)-4 : len(v.Key)]
		if string(what) == "host" {
			host = string(v.Value)
		}
		if string(what) == "port" {
			port = string(v.Value)
		}
		// if one of host&port has no value, the addr is set partly, should not return
		if host != "" && port != "" {
			addrs = append(addrs, fmt.Sprintf("%s:%s", host, port))
		}
		if host == "" && port == "" {
			empty = append(empty, string(v.Key))
		}
	}
	// for _, node := range resp.Node.Nodes {
	// 	// node should contain host & port both
	// 	host := ""
	// 	port := ""
	// 	for _, v := range node.Nodes {
	// 		// get the last 4 chars
	// 		what := v.Key[len(v.Key)-4 : len(v.Key)]
	// 		if what == "host" {
	// 			host = v.Value
	// 		}
	// 		if what == "port" {
	// 			port = v.Value
	// 		}
	// 	}
	// 	// if one of host&port has no value, the addr is set partly, should not return
	// 	if host != "" && port != "" {
	// 		addrs = append(addrs, fmt.Sprintf("%s:%s", host, port))
	// 	}
	// 	if host == "" && port == "" {
	// 		empty = append(empty, node.Key)
	// 	}
	// }
	return addrs, empty
}

func dropEmptyDir(keyapi *etcd.Client, empty []string) {
	if keyapi == nil || len(empty) == 0 {
		return
	}
	for _, key := range empty {
		_, err := keyapi.Delete(context.Background(), key)
		if err != nil {
			glog.Info("wonaming: delete empty service dir error: ", err.Error())
		}
	}
}
