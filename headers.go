package openai

import (
	"github.com/Simplou/goxios"
)

var (
	headers = []goxios.Header{}
	contentTypeJSON = goxios.Header{Key: "Content-Type", Value: "application/json"}
)

func (c *Client) setAuthorizationHeader() {
	headers = append(headers, goxios.Header{Key: "Authorization", Value: "Bearer " + c.apiKey})
}

func (c *Client) AddHeader(h goxios.Header) {
	headers = append(headers, h)
}

func Headers() []goxios.Header {
	return headers
}
