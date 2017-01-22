package xhbase

import (
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"golang.org/x/net/context"
)

const (
	_module = "hbase"
)

type Client struct {
	c gohbase.Client
	// testCell *gohbase.HBaseCell
}

func (c *Client) Scan(ctx context.Context, s *hrpc.Scan) ([]*hrpc.Result, error) {
	return c.c.Scan(s)
}

func (c *Client) Get(ctx context.Context, g *hrpc.Get) (*hrpc.Result, error) {
	return c.c.Get(g)
}

func (c *Client) Put(ctx context.Context, p *hrpc.Mutate) (*hrpc.Result, error) {
	return c.c.Put(p)
}

func (c *Client) Delete(ctx context.Context, d *hrpc.Mutate) (*hrpc.Result, error) {
	return c.c.Delete(d)
}

func (c *Client) Append(ctx context.Context, a *hrpc.Mutate) (*hrpc.Result, error) {
	return c.c.Append(a)
}

func (c *Client) Increment(ctx context.Context, i *hrpc.Mutate) (int64, error) {
	return c.c.Increment(i)
}

func (c *Client) Close() {
	c.c.Close()
}

func NewClient(zkquorum string, options ...gohbase.Option) *Client {
	return &Client{
		c: gohbase.NewClient(zkquorum),
	}
}
