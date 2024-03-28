package adguard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	"github.com/henrywhitaker3/adguard-exporter/internal/config"
)

type Client struct {
	conf config.Config
}

func NewClient(conf config.Config) *Client {
	return &Client{
		conf: conf,
	}
}

func (c *Client) do(ctx context.Context, method string, path string, out any) error {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.conf.Username, c.conf.Password)))
	url, err := url.Parse(fmt.Sprintf("%s%s", c.conf.Url, path))
	if err != nil {
		return err
	}

	req := &http.Request{
		Method: method,
		Header: http.Header{},
		URL:    url,
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d: %v", resp.StatusCode, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, out); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetStats(ctx context.Context) (*Stats, error) {
	out := &Stats{}
	err := c.do(ctx, http.MethodGet, "/control/stats", out)
	return out, err
}

func (c *Client) GetStatus(ctx context.Context) (*Status, error) {
	out := &Status{}
	err := c.do(ctx, http.MethodGet, "/control/status", out)
	return out, err
}

func (c *Client) GetDhcp(ctx context.Context) (*DhcpStatus, error) {
	out := &DhcpStatus{}
	err := c.do(ctx, http.MethodGet, "/control/dhcp/status", out)
	if err != nil {
		return nil, err
	}

	for i := range out.DynamicLeases {
		l := out.DynamicLeases[i]
		l.Type = "dynamic"
		out.DynamicLeases[i] = l
	}
	for i := range out.StaticLeases {
		l := out.StaticLeases[i]
		l.Type = "static"
		out.StaticLeases[i] = l
	}

	out.Leases = slices.Concat(out.DynamicLeases, out.StaticLeases)

	return out, nil
}

func (c *Client) Url() string {
	return c.conf.Url
}
