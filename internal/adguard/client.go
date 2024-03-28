package adguard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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

func (c *Client) get(ctx context.Context, path string) (*http.Response, error) {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.conf.Username, c.conf.Password)))
	url, err := url.Parse(fmt.Sprintf("%s%s", c.conf.Url, path))
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		Header: http.Header{},
		URL:    url,
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", auth))
	req = req.WithContext(ctx)

	return http.DefaultClient.Do(req)
}

func (c *Client) GetStats(ctx context.Context) (*Stats, error) {
	resp, err := c.get(ctx, "/control/stats")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d: %v", resp.StatusCode, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out := &Stats{}
	if err := json.Unmarshal(body, out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) Url() string {
	return c.conf.Url
}
