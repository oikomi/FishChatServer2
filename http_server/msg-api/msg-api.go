package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/http_server/msg-api/conf"
	"github.com/oikomi/FishChatServer2/http_server/msg-api/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		glog.Errorf("conf.Init() error(%v)", err)
		panic(err)
	}
	glog.Infof("msg-api [version: %s] start", conf.Conf.Ver)
	http.Init(conf.Conf)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		glog.Info("msg-api get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			glog.Infof("msg-api [version: %s] exit", conf.Conf.Ver)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
