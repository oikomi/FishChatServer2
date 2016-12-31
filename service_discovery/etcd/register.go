package etcd

import (
	"fmt"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/xtime"
	"golang.org/x/net/context"
	"strings"
	"time"
)

var rgClient *etcd.Client
var serviceKey string
var stopSignal = make(chan bool, 1)

// Register is the helper function to self-register service into Etcd/Consul server
// should call Unregister when pocess stop
// name - service name
// host - service host
// port - service port
// target - etcd dial address, for example: "http://127.0.0.1:2379;http://127.0.0.1:12379"
// interval - interval of self-register to etcd
// ttl - ttl of the register information
func Register(name string, rpcServerAddr string, target string, interval xtime.Duration, ttl xtime.Duration) (err error) {
	// get endpoints for register dial address
	endpoints := strings.Split(target, ",")
	conf := etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}
	rgClient, err = etcd.New(conf)
	if err != nil {
		glog.Error(err)
		return
	}
	serviceID := fmt.Sprintf("%s-%s", name, rpcServerAddr)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceID)
	addrKey := fmt.Sprintf("/%s/%s/%s/addr", Prefix, name, serviceID)
	go func() {
		// invoke self-register with ticker
		ticker := time.NewTicker(time.Duration(interval))
		// should get first, if not exist, set it
		for {
			_, err := rgClient.Get(context.Background(), serviceKey)
			if err != nil {
				resp, err := rgClient.Grant(context.Background(), int64(time.Duration(ttl)/time.Second))
				if err != nil {
					glog.Error(err)
				}
				if _, err = rgClient.Put(context.Background(), addrKey, rpcServerAddr, etcd.WithLease(resp.ID)); err != nil {
					glog.Error(err)
				}
				resp, err = rgClient.Grant(context.Background(), int64(time.Duration(ttl)))
				if err != nil {
					glog.Error(err)
				}
				if _, err = rgClient.Put(context.Background(), serviceKey, "", etcd.WithLease(resp.ID)); err != nil {
					glog.Error(err)
				}
			} else {
				resp, err := rgClient.Grant(context.Background(), int64(time.Duration(ttl)/time.Second))
				if err != nil {
					glog.Error(err)
				}
				if _, err = rgClient.Put(context.Background(), addrKey, rpcServerAddr, etcd.WithLease(resp.ID)); err != nil {
					glog.Error(err)
				}
			}
			select {
			case <-stopSignal:
				return
			case <-ticker.C:
			}
		}
	}()
	// initial register
	resp, err := rgClient.Grant(context.Background(), int64(time.Duration(ttl)/time.Second))
	if err != nil {
		glog.Error(err)
	}
	if _, err = rgClient.Put(context.Background(), addrKey, rpcServerAddr, etcd.WithLease(resp.ID)); err != nil {
		glog.Error(err)
		return
	}
	resp, err = rgClient.Grant(context.Background(), int64(time.Duration(ttl)/time.Second))
	if err != nil {
		glog.Error(err)
	}
	if _, err = rgClient.Put(context.Background(), serviceKey, "", etcd.WithLease(resp.ID)); err != nil {
		glog.Error(err)
		return
	}
	return
}

// Unregister delete service from etcd
func Unregister() (err error) {
	stopSignal <- true
	stopSignal = make(chan bool, 1) // just a hack to avoid multi UnRegister deadlock
	_, err = rgClient.Delete(context.Background(), serviceKey)
	if err != nil {
		glog.Error(err)
	}
	return
}
