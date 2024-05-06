package openaigo

import (
	"github.com/Simplou/goxios"
)

var headers = []goxios.Header{}

func (c *Client) setAuthorizationHeader(){
	headers = append(headers, goxios.Header{Key: "Authorization", Value: "Bearer "+c.apiKey})
}

func (c *Client) AddHeader(h goxios.Header) {
	headers = append(headers, h)
}

func Headers() []goxios.Header {
	return headers
}
