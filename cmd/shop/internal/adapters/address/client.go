package address

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"xsolla/cmd/shop/internal/app"
)

const (
	timeout = 15 * time.Second
)

type (
	Client struct {
		client *http.Client
		cfg    Config
	}

	Config struct {
		BasePath string
		APIKey string
	}
)

func New(cfg Config) *Client {
	client := &http.Client{
		Timeout: timeout,
	}

	return &Client{
		client: client,
		cfg:    cfg,
	}
}


func (c *Client) CheckAddress(ctx context.Context, a app.Address) error {
	// todo change endpoint path
	uri, err := url.Parse(fmt.Sprintf("%s%s", c.cfg.BasePath, fmt.Sprintf(c.cfg.BasePath, "")))
	if err != nil {
		return fmt.Errorf("url.Parse: %w", err)
	}

	addParams(uri, a.City, a.Street, a.House, a.Entrance, a.Flat)

	resp, err := c.request(ctx, uri.String(), http.MethodGet, nil)
	if err != nil {
		return fmt.Errorf("c.request: %w", err)
	}

	results := &responseCheckAddress{}
	err = json.Unmarshal(resp, results)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	// todo add logic to get info from response and validate it

	return nil
}

func (c *Client) request(ctx context.Context, url, method string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	// example for getting access
	req.Header.Add("", c.cfg.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.client.Do: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	// todo validate status code from response
	return nil, nil
}

func addParams(u *url.URL, _, _, _, _, _ string) {
	param := url.Values{}
	// todo add in url
	u.RawPath = param.Encode()
}
