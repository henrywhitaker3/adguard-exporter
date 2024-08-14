package adguard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/henrywhitaker3/adguard-exporter/internal/config"
	"github.com/mitchellh/mapstructure"
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

// func (c *Client) GetQueryLog(ctx context.Context) (map[string]map[string]int, []QueryTime, QueryPerClient, error) {
func (c *Client) GetQueryLog(ctx context.Context) (map[string]map[string]int, []QueryTime, []logEntry, error) {
	log := &queryLog{}
	err := c.do(ctx, http.MethodGet, "/control/querylog?limit=1000&response_status=all", log)
	if err != nil {
		return nil, nil, nil, err
	}

	types, err := c.getQueryTypes(log)
	if err != nil {
		return nil, nil, nil, err
	}
	times, err := c.getQueryTimes(log)
	if err != nil {
		return nil, nil, nil, err
	}

	return types, times, log.Log, nil
}

func (c *Client) getQueryTypes(log *queryLog) (map[string]map[string]int, error) {
	out := map[string]map[string]int{}
	for _, d := range log.Log {
		if d.Answer != nil && len(d.Answer) > 0 {
			if _, ok := out[d.Client]; !ok {
				out[d.Client] = map[string]int{}
			}
			for i := range d.Answer {
				switch v := d.Answer[i].Value.(type) {
				case string:
					out[d.Client][d.Answer[i].Type]++
				case map[string]any:
					dns65 := &type65{}
					mapstructure.Decode(v, dns65)
					out[d.Client]["TYPE"+strconv.Itoa(dns65.Hdr.Rrtype)]++
				}
			}
		}
	}
	return out, nil
}

func (c *Client) getQueryTimes(l *queryLog) ([]QueryTime, error) {
	out := []QueryTime{}
	for _, q := range l.Log {
		if q.Upstream == "" {
			q.Upstream = "self"
		}
		ms, err := strconv.ParseFloat(q.Elapsed, 32)
		if err != nil {
			log.Printf("ERROR - could not parse query elapsed time %v as float\n", q.Elapsed)
			continue
		}
		out = append(out, QueryTime{
			Elapsed:  time.Millisecond * time.Duration(ms),
			Client:   q.Client,
			Upstream: q.Upstream,
		})
	}
	return out, nil
}

func (c *Client) Url() string {
	return c.conf.Url
}
