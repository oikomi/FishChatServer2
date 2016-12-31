package etcd

import (
	"errors"
	"fmt"
	etcd "github.com/coreos/etcd/clientv3"
	// "github.com/golang/glog"
	"github.com/golang/glog"
	"google.golang.org/grpc/naming"
	"strings"
	"time"
)

// EtcdResolver is the implementaion of grpc.naming.Resolver
type EtcdResolver struct {
	ServiceName string // service name to resolve
}

// NewResolver return EtcdResolver with service name
func NewResolver(serviceName string) *EtcdResolver {
	return &EtcdResolver{ServiceName: serviceName}
}

// Resolve to resolve the service from etcd, target is the dial address of etcd
// target example: "http://127.0.0.1:2379;http://127.0.0.1:12379;http://127.0.0.1:22379"
func (er *EtcdResolver) Resolve(target string) (naming.Watcher, error) {
	glog.Info("Resolve")
	if er.ServiceName == "" {
		return nil, errors.New("service_discovery: no service name provided")
	}
	// generate etcd client, return if error
	endpoints := strings.Split(target, ",")
	conf := etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	}
	client, err := etcd.New(conf)
	if err != nil {
		return nil, fmt.Errorf("service_discovery: creat etcd error: %s", err.Error())
	}
	// Return EtcdWatcher
	watcher := &EtcdWatcher{
		er: er,
		ec: client,
	}
	return watcher, nil
}
