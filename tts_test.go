package openai

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Simplou/goxios"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing purposes.
type MockHTTPClient struct{}

func (c *MockHTTPClient) Post(url string, opts *goxios.RequestOpts) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("fake audio data")),
	}
	return resp, nil
}
func (c *MockHTTPClient) Get(url string, opts *goxios.RequestOpts) (*http.Response, error) {
	return &http.Response{}, nil
}

// MockFailingHTTPClient is a mock implementation of HTTPClient for testing purposes.
// This implementation always returns an error.
type MockFailingHTTPClient struct{}

func (c *MockFailingHTTPClient) Post(url string, opts *goxios.RequestOpts) (*http.Response, error) {
	json := goxios.JSON{}
	b, err := json.Marshal()
	if err != nil {
		return nil, err
	}
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(ioReader(b)),
	}
	return resp, errors.New("mock HTTP client always fails")
}

func (c *MockFailingHTTPClient) Get(url string, opts *goxios.RequestOpts) (*http.Response, error) {
	return &http.Response{}, nil
}

type OpenAIClientMock struct {
	MockClient
}

func TestTextToSpeech(t *testing.T) {
	client := &OpenAIClientMock{
		MockClient: MockClient{
			baseUrl: "https://fake.api.openai.com/v1",
		},
	}

	testCases := []struct {
		name            string
		httpClient      HTTPClient
		expectedErr     bool
		expectedMessage string
	}{
		{
			name:        "Successful request",
			httpClient:  new(MockHTTPClient),
			expectedErr: false,
		},
		{
			name:        "Failed request",
			httpClient:  new(MockFailingHTTPClient),
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := &SpeechRequestBody{
				Model: "fake-model",
				Input: "Hello, world!",
				Voice: "fake-voice",
			}

			resp, err := TextToSpeech(client, tc.httpClient, body)
			if tc.expectedErr && err == nil {
				t.Errorf("Expected an error, but got nil")
			}

			if !tc.expectedErr && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			if resp != nil {
				defer resp.Close()
				actualData, err := io.ReadAll(resp)
				if err != nil {
					t.Errorf("Error reading response body: %v", err)
				}
				expectedData := []byte("fake audio data")
				if !bytes.Equal(actualData, expectedData) {
					t.Errorf("Unexpected response body. Expected: %s, Got: %s", expectedData, actualData)
				}
			}
		})
	}
}
