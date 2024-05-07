package openai

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/Simplou/goxios"
)

type MockClient struct {
	baseUrl string
}

func (c MockClient) Context() context.Context {
	return context.TODO()
}
func (c MockClient) ApiKey() string {
	return "mock_api_key"
}
func (c MockClient) BaseURL() string {
	return c.baseUrl
}

func (c MockClient) AddHeader(h goxios.Header) {}

type testHTTPClient struct {
	req *http.Request
}

func (c *testHTTPClient) Post(url string, opts *goxios.RequestOpts) (*http.Response, error) {
	statusCode := http.StatusOK
	body := goxios.JSON{
		"id":    "123",
		"model": "gpt-3.5-turbo",
		"choices": []Choice{
			{
				Index:        0,
				Message:      Message{Role: "assistant", Content: "Hi"},
				FinishReason: "",
			},
		},
	}
	resBytes, err := body.Marshal()
	if err != nil {
		return nil, err
	}
	resReader := bytes.NewBuffer(resBytes)
	res := &http.Response{
		Request:       c.req,
		Status:        http.StatusText(statusCode),
		StatusCode:    statusCode,
		Header:        http.Header{},
		Body:          io.NopCloser(resReader),
		ContentLength: int64(len(resBytes)),
	}

	return res, nil
}

func TestChatCompletionRequest(t *testing.T) {
	mockClient := MockClient{"http://localhost:399317"}
	httpClient := testHTTPClient{}
	completionRequest := &CompletionRequest{
		Model:    "gpt-3.5-turbo",
		Messages: []Message{{Role: "user", Content: "Hello!"}},
	}

	response, err := ChatCompletion(mockClient, &httpClient, completionRequest)
	if err != nil {
		t.Errorf("Erro ao chamar ChatCompletion: %v", err)
	}

	expectedID := "123"
	if response.ID != expectedID {
		t.Errorf("ID da resposta esperado: %s, ID recebido: %s", expectedID, response.ID)
	}
}
