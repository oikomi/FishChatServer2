package xweb

import (
	"github.com/oikomi/FishChatServer2/common/net/trace"
	"github.com/oikomi/FishChatServer2/common/net/xweb/context"
	ctx "golang.org/x/net/context"
	"net/http"
)

const (
	_family = "go_http_server"
)

// web http pattern router
type Router struct {
	mux     *http.ServeMux
	pattern string
}

// NewRouter new a router.
func NewRouter(mux *http.ServeMux) *Router {
	return &Router{mux: mux}
}

func (r *Router) join(pattern string) string {
	return r.pattern + pattern
}

func (r *Router) Group(pattern string) *Router {
	return &Router{mux: r.mux, pattern: r.join(pattern)}
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *Router) Handle(method, pattern string, handlers ...Handler) {
	r.mux.HandleFunc(r.join(pattern), func(w http.ResponseWriter, r *http.Request) {
		handler(method, w, r, handlers)
	})
	return
}

func (r *Router) HandlerFunc(method, pattern string, handlers ...HandlerFunc) {
	r.mux.HandleFunc(r.join(pattern), func(w http.ResponseWriter, r *http.Request) {
		handleFunc(method, w, r, handlers)
	})
	return
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (r *Router) Get(pattern string, handlers ...Handler) {
	r.mux.HandleFunc(r.join(pattern), func(w http.ResponseWriter, r *http.Request) {
		handler("GET", w, r, handlers)
	})
	return
}

func (r *Router) Post(pattern string, handlers ...Handler) {
	r.mux.HandleFunc(r.join(pattern), func(w http.ResponseWriter, r *http.Request) {
		handler("POST", w, r, handlers)
	})
	return
}

// GetFunc is a shortcut for router.HandleFunc("GET", path, handle)
func (r *Router) GetFunc(pattern string, handlers ...HandlerFunc) {
	r.mux.HandleFunc(r.join(pattern), func(w http.ResponseWriter, r *http.Request) {
		handleFunc("GET", w, r, handlers)
	})
	return
}

// PostFunc is a shortcut for router.HandleFunc("GET", path, handle)
func (r *Router) PostFunc(pattern string, handlers ...HandlerFunc) {
	r.mux.HandleFunc(r.join(pattern), func(w http.ResponseWriter, r *http.Request) {
		handleFunc("POST", w, r, handlers)
	})
	return
}

func handler(method string, w http.ResponseWriter, r *http.Request, handlers []Handler) {
	t := trace.WithHTTP(r)
	t.ServerReceive(_family, r.URL.Path, "")
	defer t.ServerSend()
	if r.Method != method {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	c := context.NewContext(trace.NewContext(ctx.Background(), t), r, w)
	for _, h := range handlers {
		h.ServeHTTP(c)
		if err := c.Err(); err != nil {
			break
		}
	}
}

func handleFunc(method string, w http.ResponseWriter, r *http.Request, handlers []HandlerFunc) {
	t := trace.WithHTTP(r)
	t.ServerReceive(_family, r.URL.Path, "")
	defer t.ServerSend()
	if r.Method != method {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	c := context.NewContext(trace.NewContext(ctx.Background(), t), r, w)
	for _, h := range handlers {
		h(c)
		if err := c.Err(); err != nil {
			break
		}
	}
}
