package simpleanalytics

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	userID  string
	apiKey  string
	baseURL string
	client  *http.Client
}

const defaultURL = "https://simpleanalytics.com"

var defaultHTTPClient = http.DefaultClient

type Option func(*Client)

func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

func NewClient(userId, apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL: defaultURL,
		client:  defaultHTTPClient,
		userID:  userId,
		apiKey:  apiKey,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type HTTPError struct {
	Code    int
	Message string
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("status %d (%v)", e.Code, http.StatusText(e.Code))
}

func (c *Client) get(ctx context.Context, path string, query url.Values) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s?%s", c.baseURL+path, query.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating request: %w", err)
	}
	req.Header.Set("User-Id", c.userID)
	req.Header.Set("Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/x-ndjson")
	r, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		return nil, HTTPError{Code: r.StatusCode, Message: string(body)}
	}
	return r.Body, nil
}
