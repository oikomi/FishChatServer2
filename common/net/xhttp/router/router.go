package router

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	"github.com/oikomi/FishChatServer2/common/net/xweb"
	wctx "github.com/oikomi/FishChatServer2/common/net/xweb/context"
	"net/http"
	"strconv"
)

type Identify interface {
	Access(c wctx.Context, ip string) (mid int64, err error)
	Verify(c wctx.Context) (err error)
}

// Router is http router.
// Statsd for stat, need set value by called.
type Router struct {
	m *http.ServeMux
	r *xweb.Router
}

func New(m *http.ServeMux) *Router {
	r := &Router{
		m: m,
		r: xweb.NewRouter(m),
	}
	return r
}

// Group add group paths.
// eg: r.Group("/x/auth", func(cr *Router){
// 		cr.Get("login", handlerFunc)
// })
func (r *Router) Group(p string, f func(r *Router)) {
	or := r.r
	cr := or.Group(p)
	r.r = cr
	f(r)
	r.r = or
}

// Get support http method Get, no business access and verify.
func (r *Router) Get(p string, hf xweb.HandlerFunc) {
	r.r.GetFunc(p, r.preHandler, hf, r.writerHandler)
}

// Post support http method Post, no business access and verify.
func (r *Router) Post(p string, hf xweb.HandlerFunc) {
	r.r.PostFunc(p, r.preHandler, hf, r.writerHandler)
}

// GuestGet support http method Get, identify access will be called before hf, but cann't affect hf.
func (r *Router) GuestGet(p string, hf xweb.HandlerFunc) {
	r.r.GetFunc(p, r.preHandler, r.guestHandler, hf, r.writerHandler)
}

// GuestPost support http method Post, identify access will be called before hf, but cann't affect hf.
func (r *Router) GuestPost(p string, hf xweb.HandlerFunc) {
	r.r.PostFunc(p, r.preHandler, r.guestHandler, hf, r.writerHandler)
}

// UserGet support http method Get, identify access will be called before hf, and hf not be called when access failed.
func (r *Router) UserGet(p string, hf xweb.HandlerFunc) {
	r.r.GetFunc(p, r.preHandler, r.userHandler, hf, r.writerHandler)
}

// UserPost support http method Post, identify access will be called before hf, and hf not be called when access failed.
func (r *Router) UserPost(p string, hf xweb.HandlerFunc) {
	r.r.PostFunc(p, r.preHandler, r.userHandler, hf, r.writerHandler)
}

// VerifyGet support http method Get, identify verify will be called before hf, and hf not be called when verify failed.
func (r *Router) VerifyGet(p string, hf xweb.HandlerFunc) {
	r.r.GetFunc(p, r.preHandler, r.verifyHandler, hf, r.writerHandler)
}

// VerifyPost support http method Post, identify verify will be called before hf, and hf not be called when verify failed.
func (r *Router) VerifyPost(p string, hf xweb.HandlerFunc) {
	r.r.PostFunc(p, r.preHandler, r.verifyHandler, hf, r.writerHandler)
}

// Degrade add api that can degrade other api.
func (r *Router) Degrade(p string) {
	r.r.GetFunc(p, r.preHandler, r.verifyHandler, r.degradeHandler, r.writerHandler)
}

func (r *Router) preHandler(c wctx.Context) {
	req := c.Request()
	req.ParseForm()
}

func (r *Router) guestHandler(c wctx.Context) {
}

func (r *Router) userHandler(c wctx.Context) {
	c.Result()["code"] = ecode.OK
	r.writerHandler(c)
	c.Cancel()
	return
}

func (r *Router) verifyHandler(c wctx.Context) {
	c.Result()["code"] = ecode.OK
	r.writerHandler(c)
	c.Cancel()
	return
}

// writerHandler is a json writer for http server.
func (r *Router) writerHandler(c wctx.Context) {
	var (
		bs  []byte
		err error
		// tid    = "-"
		ret    = ecode.OK
		req    = c.Request()
		resp   = c.Response()
		res    = c.Result()
		params = req.Form
		// ip     = c.RemoteIP()
	)
	defer func() {
	}()
	if ec := res["code"]; ec != nil {
		if code, ok := ec.(error); !ok {
			glog.Errorf("http(%s) response have not ecode return", req.URL.Path)
			return
		} else {
			ret = ecode.From(code)
		}
	}
	res["code"] = ret
	res["message"] = ret.String()
	if bs, err = json.Marshal(res); err != nil {
		glog.Errorf("json.Marshal(%v) error(%v)", res, err)
		return
	}
	resp.Header().Set("Content-Type", "application/json;charset=utf-8")
	if _, err = resp.Write(bs); err != nil {
		glog.Errorf("c.Response.Write(%s, %s) failed (%v)", req.URL.Path, params.Encode(), err)
	}
}

func (r *Router) degradeHandler(c wctx.Context) {
	res := c.Result()
	params := c.Request().Form
	pt := params.Get("path")
	per := params.Get("percent")
	// check params
	if pt == "" && per == "" {
		res["code"] = ecode.OK
		// res["degrades"] = r.degrades
		return
	}
	percent, err := strconv.Atoi(per)
	if err != nil || percent < 0 || percent > 100 {
		res["code"] = ecode.RequestErr
		// res["degrades"] = r.degrades
		res["message"] = "percent is not number or in [0,100]"
		return
	}
	res["code"] = ecode.OK
}
