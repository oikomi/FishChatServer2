package http

import (
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/common/net/xhttp"
	"github.com/oikomi/FishChatServer2/common/net/xhttp/router"
	wctx "github.com/oikomi/FishChatServer2/common/net/xweb/context"
	"github.com/oikomi/FishChatServer2/http_server/group-api/conf"
	"github.com/oikomi/FishChatServer2/http_server/group-api/service"
	"net/http"
	"strconv"
)

var (
	groupSvc *service.Service
)

func Init(c *conf.Config) {
	var err error
	groupSvc, err = service.New()
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
	r.Group("/x/group", func(cr *router.Router) {
		cr.GuestPost("/create", createGroup)
		cr.GuestPost("/join", joinGroup)
		cr.GuestPost("/quit", quitGroup)
	})
	return
}

// innerRouter init local router api path.
func localRouter(r *router.Router) {
}

func createGroup(c wctx.Context) {
	res := c.Result()
	uidStr := c.Request().Form.Get("uid")
	groupNameStr := c.Request().Form.Get("name")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		glog.Error(err)
		res["code"] = ecode.RequestErr
		return
	}
	res["code"] = groupSvc.CreateGroup(uid, groupNameStr)
}

func joinGroup(c wctx.Context) {
	res := c.Result()
	uidStr := c.Request().Form.Get("uid")
	gidStr := c.Request().Form.Get("gid")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		glog.Error(err)
		res["code"] = ecode.RequestErr
		return
	}
	gid, err := strconv.ParseInt(gidStr, 10, 64)
	if err != nil {
		glog.Error(err)
		res["code"] = ecode.RequestErr
		return
	}
	res["code"] = groupSvc.JoinGroup(uid, gid)
}

func quitGroup(c wctx.Context) {
	res := c.Result()
	uidStr := c.Request().Form.Get("uid")
	gidStr := c.Request().Form.Get("gid")
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		glog.Error(err)
		res["code"] = ecode.RequestErr
		return
	}
	gid, err := strconv.ParseInt(gidStr, 10, 64)
	if err != nil {
		glog.Error(err)
		res["code"] = ecode.RequestErr
		return
	}
	res["code"] = groupSvc.QuitGroup(uid, gid)
}
