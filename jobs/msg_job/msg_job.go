package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/conf"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/conf_discovery"
	"github.com/oikomi/FishChatServer2/jobs/msg_job/service"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

var (
	s *service.Service
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		glog.Error(err)
		panic(err)
	}
	go conf_discovery.ConfDiscoveryProc()
	s = service.New(conf.Conf)
	signalHandler()
}

func signalHandler() {
	var (
		err error
		ch  = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			glog.Infof("get a signal %s, stop the consume process", si.String())
			if err = s.Close(); err != nil {
				glog.Error("close consumer error :", err)
			}
			s.Wait()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
