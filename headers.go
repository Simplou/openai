package openai

import (
	"github.com/Simplou/goxios"
)

var (
	contentTypeJSON = goxios.Header{Key: "Content-Type", Value: "application/json"}
)

func (c *Client) setAuthorizationHeader() {
	c.headers = append(c.headers, goxios.Header{Key: "Authorization", Value: "Bearer " + c.apiKey})
}

func (c *Client) AddHeader(h goxios.Header) {
	c.headers = append(c.headers, h)
}

func (c *Client) Headers() []goxios.Header {
	return c.headers
}
