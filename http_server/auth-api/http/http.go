package http

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/common/net/xhttp"
	"github.com/oikomi/FishChatServer2/common/net/xhttp/router"
	wctx "github.com/oikomi/FishChatServer2/common/net/xweb/context"
	"github.com/oikomi/FishChatServer2/http_server/auth-api/conf"
	"github.com/oikomi/FishChatServer2/http_server/auth-api/service"
	"net/http"
	"strconv"
)

var (
	authSvc *service.Service
)

func Init(c *conf.Config) {
	var err error
	authSvc, err = service.New()
	if err != nil {
		glog.Error(err)
	}
	// init external router
	extM := http.NewServeMux()
	extR := router.New(extM)
	outerRouter(extR)
	// init Outer serve
	if err = xhttp.Serve(extM, c.MultiHTTP.Outer); err != nil {
		glog.Errorf("xhttp.Serve error(%v)", err)
		panic(err)
	}
	// init local router
	localM := http.NewServeMux()
	localR := router.New(localM)
	localRouter(localR)
	// init local server
	if err = xhttp.Serve(localM, c.MultiHTTP.Local); err != nil {
		glog.Errorf("xhttp.Serve error(%v)", err)
		panic(err)
	}
}

// outerRouter init outer router api path.
func outerRouter(r *router.Router) {
	r.Group("/x/auth", func(cr *router.Router) {
		cr.GuestPost("/login", login)
	})
	return
}

// innerRouter init local router api path.
func localRouter(r *router.Router) {
	r.Get("/history/monitor/ping", ping)
}

func login(c wctx.Context) {
	res := c.Result()
	uidStr := c.Request().Form.Get("uid")
	pwStr := c.Request().Form.Get("pw")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		glog.Error(err)
		res["code"] = ecode.RequestErr
		return
	}
	res["data"], res["code"] = authSvc.Login(uid, pwStr)
}

// ping check server ok.
func ping(c wctx.Context) {
	// res := c.Result()
	// if err := hisSvc.Ping(c); err != nil {
	// 	res["code"] = err
	// 	log.Error("history service ping error(%v)", err)
	// 	http.Error(c.Response(), "", http.StatusServiceUnavailable)
	// 	c.Done()
	// }
}

// version check server ver.
// func version(c wctx.Context) {
// 	res := c.Result()
// 	res["data"] = map[string]interface{}{
// 		"version": conf.Conf..Version,
// 	}
// }
