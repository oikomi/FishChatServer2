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

var rgClient *etcd.Client
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
	rgClient, err = etcd.New(conf)
	if err != nil {
		glog.Error(err)
		return
	}
	serviceID := fmt.Sprintf("%s-%s-%d", name, host, port)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceID)
	addrKey := fmt.Sprintf("/%s/%s/%s/addr", Prefix, name, serviceID)
	go func() {
		// invoke self-register with ticker
		ticker := time.NewTicker(interval)
		// should get first, if not exist, set it
		for {
			<-ticker.C
			glog.Info(serviceKey)
			getResp, err := rgClient.Get(context.Background(), serviceKey)
			glog.Info(getResp.Kvs)
			if err != nil {
				glog.Error(err)
				if _, err = rgClient.Put(context.Background(), addrKey, host+":"+fmt.Sprintf("%d", port)); err != nil {
					glog.Error(err)
				}
				resp, err := rgClient.Grant(context.Background(), ttl)
				if err != nil {
					glog.Error(err)
				}
				if _, err = rgClient.Put(context.Background(), serviceKey, "", etcd.WithLease(resp.ID)); err != nil {
					glog.Error(err)
				}
			} else {
				glog.Info("err is nil")
				resp, err := rgClient.Grant(context.Background(), ttl)
				if err != nil {
					glog.Error(err)
				}
				pres, err := rgClient.Put(context.Background(), serviceKey, "", etcd.WithLease(resp.ID))
				if err != nil {
					glog.Error(err)
				}
				glog.Info(pres)
			}
		}
	}()
	// initial register
	if _, err = rgClient.Put(context.Background(), addrKey, host+":"+fmt.Sprintf("%d", port)); err != nil {
		glog.Error(err)
		return
	}
	resp, err := rgClient.Grant(context.Background(), ttl)
	if err != nil {
		glog.Error(err)
	}
	if _, err = rgClient.Put(context.Background(), serviceKey, "", etcd.WithLease(resp.ID)); err != nil {
		glog.Error(err)
		return
	}
	glog.Info("end")
	return
}

// Unregister delete service from etcd
func Unregister() error {
	_, err := rgClient.Delete(context.Background(), serviceKey)
	if err != nil {
		log.Println("wonaming: deregister service error: ", err.Error())
	} else {
		log.Println("wonaming: deregistered service from etcd server.")
	}
	return err
}
