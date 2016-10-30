package xweb

import (
	"github.com/oikomi/FishChatServer2/common/net/xweb/context"
)

type Handler interface {
	ServeHTTP(context.Context)
}

type HandlerFunc func(context.Context)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(c context.Context) {
	f(c)
}
