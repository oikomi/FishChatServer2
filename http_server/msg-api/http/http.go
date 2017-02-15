package http

import (
	"github.com/golang/glog"
	// "github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/common/net/xhttp"
	"github.com/oikomi/FishChatServer2/common/net/xhttp/router"
	// wctx "github.com/oikomi/FishChatServer2/common/net/xweb/context"
	"github.com/oikomi/FishChatServer2/http_server/msg-api/conf"
	"github.com/oikomi/FishChatServer2/http_server/msg-api/service"
	"net/http"
	// "strconv"
)

var (
	managerSvc *service.Service
)

func Init(c *conf.Config) {
	var err error
	managerSvc, err = service.New()
	if err != nil {
		glog.Error(err)
		return
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
	glog.Info("outerRouter")
	r.Group("/x/msg", func(cr *router.Router) {
		// cr.GuestGet("/offline", offlineMsgs)
	})
	return
}

// innerRouter init local router api path.
func localRouter(r *router.Router) {
}

// func offlineMsgs(c wctx.Context) {
// 	glog.Info("offlineMsgs")
// 	res := c.Result()
// 	uidStr := c.Request().Form.Get("uid")
// 	uid, err := strconv.ParseInt(uidStr, 10, 64)
// 	if err != nil {
// 		glog.Error(err)
// 		res["code"] = ecode.RequestErr
// 		return
// 	}
// 	res["data"], res["code"] = managerSvc.GetOfflineMsgs(uid)
// }
