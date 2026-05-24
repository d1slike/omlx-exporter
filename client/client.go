package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/txix-open/isp-kit/http/httpcli"
	"github.com/txix-open/isp-kit/http/httpclix"
	"github.com/txix-open/isp-kit/log"
)

const (
	sessionCookieName = "omlx_admin_session"
	loginEndpoint     = "/admin/api/login"
)

type Client struct {
	httpCli      *httpcli.Client
	apiKey       string
	sessionToken string
}

func New(baseURL, apiKey string, logger log.Logger) *Client {
	cli := httpcli.New(httpcli.WithMiddlewares(httpclix.Log(logger)))
	cli.GlobalRequestConfig().BaseUrl = baseURL
	return &Client{
		httpCli: cli,
		apiKey:  apiKey,
	}
}

func (c *Client) RefreshSession(ctx context.Context) error {
	resp, err := c.httpCli.Post(loginEndpoint).
		RequestBody([]byte(`{"api_key":"`+c.apiKey+`","remember":true}`)).
		Header("Content-Type", "application/json").
		Do(ctx)
	if err != nil {
		return fmt.Errorf("login request: %w", err)
	}
	for _, cookie := range resp.Raw.Cookies() {
		if cookie.Name == sessionCookieName {
			c.sessionToken = cookie.Value
			break
		}
	}
	if c.sessionToken == "" {
		return fmt.Errorf("no session cookie received from login")
	}
	return nil
}

func (c *Client) FetchStats(ctx context.Context) (*StatsResponse, error) {
	if c.sessionToken == "" {
		return nil, fmt.Errorf("no session token, call RefreshSession first")
	}
	var result StatsResponse
	err := c.httpCli.Get("/admin/api/stats").
		QueryParams(map[string]any{"scope": "session"}).
		Cookie(&http.Cookie{Name: sessionCookieName, Value: c.sessionToken}).
		StatusCodeToError().
		JsonResponseBody(&result).
		DoWithoutResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch stats: %w", err)
	}
	return &result, nil
}

func (c *Client) NeedsSessionRefresh(err error) bool {
	var httpErr httpcli.ErrorResponse
	return errors.As(err, &httpErr) && httpErr.StatusCode == 401
}
