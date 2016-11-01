package xhttp

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/conf"
	"github.com/oikomi/FishChatServer2/common/net/netutil"
	"net"
	"net/http"
	"runtime"
	"time"
)

// Serve listen and serve http handlers with limit count.
func Serve(mux *http.ServeMux, c *conf.HTTPServer) (err error) {
	for _, addr := range c.Addrs {
		l, err := net.Listen("tcp", addr)
		if err != nil {
			glog.Errorf("net.Listen(\"tcp\", \"%s\") error(%v)", addr, err)
			return err
		}
		if c.MaxListen > 0 {
			l = netutil.LimitListener(l, c.MaxListen)
		}
		glog.Infof("start http listen addr: %s", addr)
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				server := &http.Server{Handler: mux, ReadTimeout: time.Duration(c.ReadTimeout), WriteTimeout: time.Duration(c.WriteTimeout)}
				if err := server.Serve(l); err != nil {
					glog.Errorf("server.Serve error(%v)", err)
				}
			}()
		}
	}
	return nil
}
