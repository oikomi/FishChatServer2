package router

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/ecode"
	// "github.com/oikomi/FishChatServer2/common/net/trace"
	"github.com/oikomi/FishChatServer2/common/net/xweb"
	wctx "github.com/oikomi/FishChatServer2/common/net/xweb/context"
	"net/http"
	// "net/url"
	"strconv"
	// "time"
)

// var (
// 	allowOrigin = map[string]string{
// 		"www.bilibili.com":     "http://www.bilibili.com",
// 		"space.bilibili.com":   "http://space.bilibili.com",
// 		"member.bilibili.com":  "http://member.bilibili.com",
// 		"search.bilibili.com":  "http://search.bilibili.com",
// 		"bangumi.bilibili.com": "http://bangumi.bilibili.com",
// 	}
// )

type Identify interface {
	Access(c wctx.Context, ip string) (mid int64, err error)
	Verify(c wctx.Context) (err error)
}

// Router is http router.
// Statsd for stat, need set value by called.
type Router struct {
	m *http.ServeMux
	r *xweb.Router
	// Identify Identify
	// degrades map[string]degrade.Ratio
}

func New(m *http.ServeMux) *Router {
	r := &Router{
		m: m,
		r: xweb.NewRouter(m),
		// degrades: map[string]degrade.Ratio{},
	}
	return r
}

// Group add group paths.
// eg: r.Group("/x/reply", func(cr *Router){
// 		cr.Get("add", handlerFunc)
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

// func (r *Router) isDegrade(p string) bool {
// 	ratio, ok := r.degrades[p]
// 	if !ok {
// 		return false
// 	}
// 	return !ratio.Pass()
// }

func (r *Router) preHandler(c wctx.Context) {
	req := c.Request()
	// res := c.Response()
	// if r.isDegrade(req.URL.Path) {
	// 	c.Cancel()
	// 	res := c.Result()
	// 	res["code"] = ecode.Degrade
	// 	r.writerHandler(c)
	// 	return
	// }
	req.ParseForm()
	// if req.Form.Get("jsonp") == "jsonp" || req.Form.Get("cross_domain") == "true" {
	// 	u, err := url.Parse(req.Referer())
	// 	if err == nil {
	// 		// if origin, ok := allowOrigin[u.Host]; ok {
	// 		// 	res.Header().Set("Access-Control-Allow-Origin", origin)
	// 		// 	res.Header().Set("Access-Control-Allow-Credentials", "true")
	// 		// }
	// 	}
	// }
}

func (r *Router) guestHandler(c wctx.Context) {
	var (
	// req = c.Request()
	// mid, _ = r.Identify.Access(c, c.RemoteIP())
	)
	// if mid == 0 {
	// 	midStr := req.Form.Get("mid")
	// 	if midStr != "" {
	// 		var err error
	// 		mid, err = strconv.ParseInt(midStr, 10, 64)
	// 		if err != nil {
	// 			glog.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
	// 		}
	// 	}
	// }
	// c.Set("mid", mid)
}

func (r *Router) userHandler(c wctx.Context) {
	// var (
	// 	mid, err = r.Identify.Access(c, c.RemoteIP())
	// )
	// if (err != nil && err != ecode.OK) || mid == 0 {
	// 	if err == nil || err == ecode.OK {
	// 		err = ecode.RequestErr
	// 	}
	c.Result()["code"] = ecode.OK
	r.writerHandler(c)
	c.Cancel()
	return
	// }
	// c.Set("mid", mid)
}

func (r *Router) verifyHandler(c wctx.Context) {
	// err := r.Identify.Verify(c)
	// if err != nil && err != ecode.OK {
	c.Result()["code"] = ecode.OK
	r.writerHandler(c)
	c.Cancel()
	return
	// }
	// midStr := c.Request().Form.Get("mid")
	// if midStr != "" {
	// 	mid, err := strconv.ParseInt(midStr, 10, 64)
	// 	if err != nil {
	// 		glog.Error("strconv.ParseInt(%s) error(%v)", midStr, err)
	// 	}
	// 	c.Set("mid", mid)
	// }
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
		// p := req.URL.Path[1:]
		// tmsub := time.Now().Sub(c.Now())
		// if r.Statsd != nil {
		// 	r.Statsd.Incr(fmt.Sprintf("%s.%d", p, ret))
		// 	r.Statsd.Timing(p, int64(tmsub/time.Millisecond))
		// }
		// if t, ok := trace.FromContext(c); ok && t.Sampled {
		// 	tid = t.ID
		// }
		// glog.InfoTrace(tid, req.URL.Path, fmt.Sprintf("method:%s,ip:%s,params:%s,ret:%d", req.Method, ip, params.Encode(), ret), tmsub.Seconds())
	}()
	if ec := res["code"]; ec != nil {
		if code, ok := ec.(error); !ok {
			glog.Error("http(%s) response have not ecode return", req.URL.Path)
			return
		} else {
			ret = ecode.From(code)
		}
	}
	res["code"] = ret
	res["message"] = ret.String()
	if bs, err = json.Marshal(res); err != nil {
		glog.Error("json.Marshal(%v) error(%v)", res, err)
		return
	}
	resp.Header().Set("Content-Type", "application/json;charset=utf-8")
	if params.Get("jsonp") == "jsonp" {
		cb := params.Get("callback")
		if cb != "" {
			bs = []byte(fmt.Sprintf("%s(%s)", cb, bs))
		}
		script := params.Get("script")
		if script == "script" {
			resp.Header().Set("Content-Type", "text/html;charset=utf-8")
			bs = []byte(fmt.Sprintf(
				`<script type="text/javascript">
					document.domain = 'bilibili.com';
					window.parent.%s;
				</script>`, bs))
		}
	}
	if _, err = resp.Write(bs); err != nil {
		glog.Error("c.Response.Write(%s, %s) failed (%v)", req.URL.Path, params.Encode(), err)
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
	// dg := map[string]degrade.Ratio{}
	// ratio := degrade.Ratio{}
	// ratio.SetSeed(percent)
	// dg[pt] = ratio
	// for k, v := range r.degrades {
	// 	if k != pt {
	// 		dg[k] = v
	// 	}
	// }
	// r.degrades = dg
	res["code"] = ecode.OK
	// res["degrades"] = r.degrades
}
