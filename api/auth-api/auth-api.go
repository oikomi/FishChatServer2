package main

import (
	"flag"
	"github.com/oikomi/FishChatServer2/api/auth-api/conf"
	"github.com/oikomi/FishChatServer2/api/auth-api/http"
	// "github.com/oikomi/FishChatServer2/common/net/trace"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		glog.Errorf("conf.Init() error(%v)", err)
		panic(err)
	}
	glog.Info("auth-api [version: %s] start", conf.Conf.Ver)
	http.Init(conf.Conf)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		s := <-c
		glog.Info("auth-api get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			glog.Info("auth-api [version: %s] exit", conf.Conf.Ver)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
