package trace

import (
	"encoding/hex"
	"fmt"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	_ratio  = float32(1 / 1024.0)
	_ev     chan []byte
	_writer io.Writer
)

const (
	_clientStart   = 0
	_serverSend    = 1
	_serverReceive = 2
	_clientReceive = 3
	_userDefine    = 4

	_httpHeaderID       = "x1-bilispy-id"
	_httpHeaderSpanID   = "x1-bilispy-spanid"
	_httpHeaderParentID = "x1-bilispy-parentid"
	_httpHeaderSampled  = "x1-bilispy-sampled"
)

// Init init the trace, must set a trace writer and id generator.
func Init(writer io.Writer) {
	_ev = make(chan []byte, 10240)
	_writer = writer
	go writeproc()
}

// SetRatio set the trace ratio.
func SetRatio(ratio float32) {
	_ratio = ratio
}

// Stop stop the trace.
func Stop() {
	if _ev != nil {
		_ev <- nil
	}
}

func writeproc() {
	for {
		if d := <-_ev; d != nil {
			if _, err := _writer.Write(d); err != nil {
				glog.Error("_writer.Write() error(%v)", err)
			}
		} else {
			glog.Error("trace writeproc goroutine exit")
			return
		}
	}
}

func id() string {
	i := [16]byte(uuid.NewV1())
	return hex.EncodeToString(i[:])
}

// Trace is server and client called trace info.
type Trace struct {
	ID       string `json:"id"`
	SpanID   string `json:"span_id"`
	ParentID string `json:"parent_id"`
	Sampled  bool   `json:"sampled"`
}

// NewTrace a root trace.
func NewTrace() *Trace {
	t := new(Trace)
	t.ID = id()
	t.SpanID = t.ID
	t.ParentID = ""
	var sampled bool
	if _ratio <= 0 {
		sampled = false
	} else if _ratio >= 1 {
		sampled = true
	} else {
		sampled = (rand.Float32() <= _ratio)
	}
	t.Sampled = sampled
	return t
}

// InheritTrace fork a child trace from current trace.
func InheritTrace(id, spanID, parentID string, sampled bool) *Trace {
	t := new(Trace)
	t.ID = id
	t.SpanID = spanID
	t.ParentID = parentID
	t.Sampled = sampled
	return t
}

// WithHTTP init trace from http request.
func WithHTTP(req *http.Request) *Trace {
	var (
		sampled              bool
		id, spanID, parentID string
	)
	id = req.Form.Get(_httpHeaderID)
	spanID = req.Form.Get(_httpHeaderSpanID)
	parentID = req.Form.Get(_httpHeaderParentID)
	if str := req.Form.Get(_httpHeaderSampled); str == "true" {
		sampled = true
	} else {
		sampled = false
	}
	if id != "" && spanID != "" {
		return InheritTrace(id, spanID, parentID, sampled)
	}
	return NewTrace()
}

func record(module, name, env string, ev int, t *Trace) {
	if _ev == nil {
		return
	}
	if t.Sampled {
		select {
		case _ev <- []byte(fmt.Sprintf("%d\001%s\001%s\001%s\001%d\001%s.%s\001%s\001%d\n", time.Now().UnixNano(), t.ID, t.SpanID, t.ParentID, ev, module, name, env, 0)):
		default:
			glog.Error("trace chan full, discard the trace: %v", t)
		}
	}
}

// Fork fork a trace with different id.
func (t *Trace) Fork() *Trace {
	t1 := new(Trace)
	t1.ID = t.ID
	t1.SpanID = id()
	t1.ParentID = t.SpanID
	t1.Sampled = t.Sampled
	return t1
}

// ClientStart record the trace with a event ClientStart
// and called it when the call start.
func (t *Trace) ClientStart(module, name, env string) {
	record(module, name, env, _clientStart, t)
}

// ClientReceive record the trace with a event ClientReceive
// and called it when the call end.
func (t *Trace) ClientReceive() {
	record("-", "-", "", _clientReceive, t)
}

// ServerReceive record the trace with a event ServerReceive
// and called it when the call end.
func (t *Trace) ServerReceive(module, name, env string) {
	record(module, name, env, _serverReceive, t)
}

// ServerSend record the trace with a event ServerSend
// and called it when the call start.
func (t *Trace) ServerSend() {
	record("-", "-", "", _serverSend, t)
}

// Log record the trace with a event UserDefine and called it when you want
// more info.
func (t *Trace) Log(module, name, env string) {
	record(module, name, env, _userDefine, t)
}

// SetHTTP set trace id into http request.
func (t *Trace) SetHTTP(req *http.Request) {
	req.Header.Set(_httpHeaderID, t.ID)
	req.Header.Set(_httpHeaderSpanID, t.SpanID)
	req.Header.Set(_httpHeaderParentID, t.ParentID)
	req.Header.Set(_httpHeaderSampled, strconv.FormatBool(t.Sampled))
}

type contextKeyT string

var _contextKey = contextKeyT("go-common/net/trace.Trace")

// NewContext return a copy of the parent context
// and associates it with a trace.
func NewContext(c context.Context, t *Trace) context.Context {
	return context.WithValue(c, _contextKey, t)
}

// FromContext returns the trace bound to the context, if any.
func FromContext(c context.Context) (*Trace, bool) {
	t, ok := c.Value(_contextKey).(*Trace)
	return t, ok
}
