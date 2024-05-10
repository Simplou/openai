package openai

import (
	"context"
	"net/http"

	"github.com/Simplou/goxios"
)

type Client struct {
	ctx    context.Context
	apiKey string
}

type OpenAIClient interface {
	Context() context.Context
	ApiKey() string
	BaseURL() string
	AddHeader(goxios.Header)
}

type HTTPClient interface {
	Post(string, *goxios.RequestOpts) (*http.Response, error)
	Get(string, *goxios.RequestOpts) (*http.Response, error)
}

func (c *Client) BaseURL() string {
	return "https://api.openai.com/v1"
}

func (c *Client) Context() context.Context {
	return c.ctx
}

func (c *Client) ApiKey() string {
	return c.apiKey
}

func New(ctx context.Context, apiKey string) *Client {
	openaiClient := &Client{ctx, apiKey}
	openaiClient.setAuthorizationHeader()
	return openaiClient
}
