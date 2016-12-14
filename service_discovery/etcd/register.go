package etcd

import (
	"fmt"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"
)

var client *etcd.Client
var serviceKey string

// Register is the helper function to self-register service into Etcd/Consul server
// should call Unregister when pocess stop
// name - service name
// host - service host
// port - service port
// target - etcd dial address, for example: "http://127.0.0.1:2379;http://127.0.0.1:12379"
// interval - interval of self-register to etcd
// ttl - ttl of the register information
func Register(name string, host string, port int, target string, interval time.Duration, ttl int64) (err error) {
	// get endpoints for register dial address
	endpoints := strings.Split(target, ",")
	glog.Info(endpoints)
	conf := etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}
	client, err = etcd.New(conf)
	if err != nil {
		glog.Error(err)
		return
	}
	// keyapi := etcd.NewKeysAPI(client)
	serviceID := fmt.Sprintf("%s-%s-%d", name, host, port)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceID)
	hostKey := fmt.Sprintf("/%s/%s/%s/host", Prefix, name, serviceID)
	portKey := fmt.Sprintf("/%s/%s/%s/port", Prefix, name, serviceID)
	go func() {
		// invoke self-register with ticker
		ticker := time.NewTicker(interval)
		// should get first, if not exist, set it
		for {
			<-ticker.C
			// _, err := client.Get(context.Background(), serviceKey, &etcd.GetOptions{Recursive: true})
			glog.Info(serviceKey)
			_, err = client.Get(context.Background(), serviceKey)
			if err != nil {
				glog.Error(err)
				if _, err = client.Put(context.Background(), hostKey, host); err != nil {
					glog.Error(err)
				}
				if _, err = client.Put(context.Background(), portKey, fmt.Sprintf("%d", port)); err != nil {
					glog.Error(err)
				}
				// setopt := &etcd.SetOptions{TTL: time.Duration(ttl) * time.Second, PrevExist: etcd.PrevExist, Dir: true}
				resp, err := client.Grant(context.Background(), ttl)
				if err != nil {
					glog.Error(err)
				}
				if _, err = client.Put(context.Background(), serviceKey, "1", etcd.WithLease(resp.ID)); err != nil {
					glog.Error(err)
				}
			} else {
				glog.Info("err is nil")
				// refresh set to true for not notifying the watcher
				// setopt := &etcd.SetOptions{TTL: time.Duration(ttl) * time.Second, PrevExist: etcd.PrevExist, Dir: true, Refresh: true}
				resp, err := client.Grant(context.Background(), 100)
				if err != nil {
					glog.Error(err)
				}
				pres, err := client.Put(context.Background(), serviceKey, "1", etcd.WithLease(resp.ID))
				if err != nil {
					glog.Error(err)
				}
				glog.Info(pres)
			}
		}
	}()
	// initial register
	if _, err = client.Put(context.Background(), hostKey, host); err != nil {
		glog.Error(err)
		return
	}
	if _, err = client.Put(context.Background(), portKey, fmt.Sprintf("%d", port)); err != nil {
		glog.Error(err)
		return
	}
	// setopt := &etcd.SetOptions{TTL: time.Duration(ttl) * time.Second, PrevExist: etcd.PrevExist, Dir: true}
	resp, err := client.Grant(context.Background(), ttl)
	if err != nil {
		glog.Error(err)
	}
	if _, err = client.Put(context.Background(), serviceKey, "1", etcd.WithLease(resp.ID)); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("end")
	return
}

// Unregister delete service from etcd
func Unregister() error {
	// keyapi := etcd.NewKeysAPI(client)
	// _, err := client.Delete(context.Background(), serviceKey, &etcd.DeleteOptions{Recursive: true})
	_, err := client.Delete(context.Background(), serviceKey)
	if err != nil {
		log.Println("wonaming: deregister service error: ", err.Error())
	} else {
		log.Println("wonaming: deregistered service from etcd server.")
	}
	return err
}
