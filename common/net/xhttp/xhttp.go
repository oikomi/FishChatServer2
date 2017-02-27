package xhttp

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/conf"
	"github.com/oikomi/FishChatServer2/common/itime"
	"github.com/oikomi/FishChatServer2/common/net/trace"
	"github.com/oikomi/FishChatServer2/common/xtime"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	_get          = "GET"
	_post         = "POST"
	_allCheckFile = "/data/www/check.html"
	_moduleName   = "http_client"
)

var (
	_noKickUserAgent = "im"
)

func init() {
	n, err := os.Hostname()
	if err == nil {
		_noKickUserAgent = _noKickUserAgent + n
	}
}

// Client is http client.
type Client struct {
	conf      *conf.HTTPClient
	timer     *itime.Timer
	client    *http.Client
	dialer    *net.Dialer
	transport *http.Transport
}

// NewClient new a http client.
func NewClient(c *conf.HTTPClient) *Client {
	client := new(Client)
	client.conf = c
	client.timer = itime.NewTimer(c.Timer)
	client.dialer = &net.Dialer{
		Timeout:   time.Duration(c.Dial),
		KeepAlive: time.Duration(c.KeepAlive),
	}
	client.transport = &http.Transport{
		Dial:            client.dialer.Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.client = &http.Client{
		Transport: client.transport,
	}
	return client
}

// SetKeepAlive set http client keepalive.
func (client *Client) SetKeepAlive(d time.Duration) {
	client.dialer.KeepAlive = d
	client.conf.KeepAlive = xtime.Duration(d)
}

// SetTimeout set http client timeout.
func (client *Client) SetTimeout(d time.Duration) {
	client.conf.Timeout = xtime.Duration(d)
}

// SetDialTimeout set http client dial timeout.
func (client *Client) SetDialTimeout(d time.Duration) {
	client.dialer.Timeout = d
	client.conf.Dial = xtime.Duration(d)
}

// Get issues a GET to the specified URL.
func (client *Client) Get(c context.Context, uri, realIP string, params url.Values, res interface{}) (err error) {
	return client.Do(c, newRequest("GET", uri, realIP, params), res)
}

// Post issues a Post to the specified URL.
func (client *Client) Post(c context.Context, uri, realIP string, params url.Values, res interface{}) (err error) {
	return client.Do(c, newRequest("POST", uri, realIP, params), res)
}

// Do sends an HTTP request and returns an HTTP response.
func (client *Client) Do(c context.Context, req *http.Request, res interface{}) (err error) {
	if t, ok := trace.FromContext(c); ok {
		t = t.Fork()
		t.SetHTTP(req)
		t.ClientStart(_moduleName, req.URL.Path, "")
		defer t.ClientReceive()
	}
	req.Header.Set("User-Agent", _noKickUserAgent)
	td := client.timer.Start(time.Duration(client.conf.Timeout), func() {
		client.transport.CancelRequest(req)
	})
	resp, err := client.client.Do(req)
	td.Stop()
	if err != nil {
		glog.Errorf("httpClient.Do(%s) error(%v)", realURL(req), err)
		return
	}
	defer resp.Body.Close()
	if res == nil {
		return
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("ioutil.ReadAll(%s) uri(%s) error(%v)", bs, realURL(req), err)
		return
	}
	if err = json.Unmarshal(bs, res); err != nil {
		glog.Errorf("json.Unmarshal(%s) uri(%s) error(%v)", bs, realURL(req), err)
	}
	return
}

// Sign calc appkey and appsecret sign.
func Sign(params url.Values) (query string, err error) {
	if len(params) == 0 {
		return
	}
	if params.Get("appkey") == "" {
		err = fmt.Errorf("utils http get must have parameter appkey")
		return
	}
	if params.Get("appsecret") == "" {
		err = fmt.Errorf("utils http get must have parameter appsecret")
		return
	}
	if params.Get("sign") != "" {
		err = fmt.Errorf("utils http get must have not parameter sign")
		return
	}
	// sign
	secret := params.Get("appsecret")
	params.Del("appsecret")
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(tmp + secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	query = params.Encode()
	return
}

// newRequest new http request with method, uri, ip and values.
func newRequest(method, uri, realIP string, params url.Values) (req *http.Request) {
	enc, err := Sign(params)
	if err != nil {
		glog.Errorf("http check params or sign error(%v)", err)
		return
	}
	ru := uri
	if enc != "" {
		ru = uri + "?" + enc
	}
	if method == _get {
		req, err = http.NewRequest(_get, ru, nil)
	} else {
		req, err = http.NewRequest(_post, uri, strings.NewReader(enc))
	}
	if err != nil {
		glog.Errorf("http.NewRequest(%s, %s) error(%v)", method, ru, err)
		return
	}
	if method == _post {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("User-Agent", _noKickUserAgent)
	return
}

// realUrl return url with http://host/params.
func realURL(req *http.Request) string {
	if req.Method == "GET" {
		return req.URL.String()
	} else if req.Method == "POST" {
		ru := req.URL.Path
		if req.Body != nil {
			rd, ok := req.Body.(io.Reader)
			if ok {
				buf := bytes.NewBuffer([]byte{})
				buf.ReadFrom(rd)
				ru = ru + "?" + buf.String()
			}
		}
		return ru
	}
	return req.URL.Path
}

// InetAtoN conver ip addr to uint32.
func InetAtoN(s string) (sum uint32) {
	ip := net.ParseIP(s)
	if ip == nil {
		return
	}
	ip = ip.To4()
	if ip == nil {
		return
	}
	sum += uint32(ip[0]) << 24
	sum += uint32(ip[1]) << 16
	sum += uint32(ip[2]) << 8
	sum += uint32(ip[3])
	return sum
}

// InetNtoA conver uint32 to ip addr.
func InetNtoA(sum uint32) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte((sum >> 24) & 0xFF)
	ip[1] = byte((sum >> 16) & 0xFF)
	ip[2] = byte((sum >> 8) & 0xFF)
	ip[3] = byte(sum & 0xFF)
	return ip.String()
}
