package zookeeper

import (
	"encoding/json"
	"github.com/golang/glog"
	commconf "github.com/oikomi/FishChatServer2/common/conf"
	izk "github.com/samuel/go-zookeeper/zk"
	"path"
	"strings"
	"time"
)

const (
	_evAddServer = 0
	_evDelServer = 1
)

type event struct {
	ev int
	c  *commconf.Server
}

type zk struct {
	cli   *izk.Conn
	conf  *commconf.Zookeeper
	ev    chan event
	nodes map[*commconf.Server]string
}

func NewZKByServer(c *commconf.Zookeeper, server int) (z *zk) {
	var err error
	z = &zk{
		conf:  c,
		ev:    make(chan event, server*2),
		nodes: make(map[*commconf.Server]string, server),
	} // add & del server, double it
	if err = z.connect(); err != nil {
		go z.reconnect()
	}
	go z.serverproc()
	return
}

func NewZKByClient(c *commconf.Zookeeper) (z *zk) {
	var err error
	if c.Root != "/" {
		c.Root = strings.TrimRight(c.Root, "/")
	}
	z = &zk{
		conf: c,
	} // add & del server, double it
	if err = z.connect(); err != nil {
		go z.reconnect()
	}
	return
}

func (z *zk) connect() (err error) {
	var ev <-chan izk.Event
	if z.cli, ev, err = izk.Connect(z.conf.Addrs, time.Duration(z.conf.Timeout)); err == nil {
		go z.eventproc(ev)
	} else {
		glog.Error(err)
	}
	return
}

func (z *zk) reconnect() {
	var err error
	for err = z.connect(); err != nil; err = z.connect() {
		time.Sleep(time.Second)
	}
}

func (z *zk) eventproc(s <-chan izk.Event) {
	var (
		ok bool
		e  izk.Event
	)
	for {
		if e, ok = <-s; !ok {
			return
		}
		glog.Info("zookeeper get a event: %s", e.State.String())
	}
}

func (z *zk) AddServer(c *commconf.Server) {
	z.ev <- event{ev: _evAddServer, c: c}
}

func (z *zk) DelServer(c *commconf.Server) {
	z.ev <- event{ev: _evDelServer, c: c}
}

func (z *zk) serverproc() {
	var (
		ok   bool
		ev   event
		bs   []byte
		node string
		err  error
	)
	for {
		if ev, ok = <-z.ev; !ok {
			return
		}
		for {
			if z.cli == nil {
				time.Sleep(time.Second)
				continue
			}
			if ev.ev == _evAddServer {
				// add node
				if bs, err = json.Marshal(ev.c); err != nil {
					glog.Error(err)
					break
				}
				if node, err = z.cli.Create(z.conf.Root, bs, izk.FlagEphemeral|izk.FlagSequence, izk.WorldACL(izk.PermAll)); err != nil {
					glog.Error(err)
					time.Sleep(time.Second)
					continue
				}
				z.nodes[ev.c] = path.Join(z.conf.Root, node)
			} else if ev.ev == _evDelServer {
				// delete node
				if node, ok = z.nodes[ev.c]; ok {
					if err = z.cli.Delete(node, -1); err != nil {
						if err != izk.ErrNoNode {
							time.Sleep(time.Second)
							continue
						}
					}
					delete(z.nodes, ev.c)
				}
			}
			break
		}
	}
}

func (z *zk) servers() (svrs []*commconf.Server, ev <-chan izk.Event, err error) {
	var (
		addr  string
		node  string
		addrs []string
		bs    []byte
		svr   *commconf.Server
	)
	if z.cli == nil {
		return nil, nil, nil
	}
	if addrs, _, ev, err = z.cli.ChildrenW(z.conf.Root); err != nil {
		glog.Error(err)
		return
	} else if len(addrs) == 0 {
		glog.Warning("server(%s) have not node in zk", z.conf.Root)
		return
	}
	svrs = make([]*commconf.Server, 0, len(addrs))
	for _, addr = range addrs {
		node = path.Join(z.conf.Root, addr)
		if bs, _, err = z.cli.Get(node); err != nil {
			glog.Error(err)
			return
		}
		svr = new(commconf.Server)
		if err = json.Unmarshal(bs, svr); err != nil {
			glog.Error(err)
			return
		}
		svrs = append(svrs, svr)
	}
	return
}
