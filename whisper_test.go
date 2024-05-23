package openai

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/Simplou/goxios"
)

// MockWhisperHTTPClient is a mock implementation of HTTPClient for testing purposes.
type MockWhisperHTTPClient struct{}

func (c *MockWhisperHTTPClient) Get(url string, opts *goxios.RequestOpts) (*http.Response, error) {
	return &http.Response{}, nil
}

func (c *MockWhisperHTTPClient) Post(url string, opts *goxios.RequestOpts) (*http.Response, error) {
	json := goxios.JSON{
		"text": "Hello.",
	}
	b, err := json.Marshal()
	if err != nil {
		return nil, err
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(ioReader(b)),
	}
	return resp, nil
}

// MockFailingWhisperHTTPClient is a mock implementation of HTTPClient for testing purposes.
// This implementation always returns an error.
type MockFailingWhisperHTTPClient struct{}

func (c *MockFailingWhisperHTTPClient) Get(url string, opts *goxios.RequestOpts) (*http.Response, error) {
	return &http.Response{}, nil
}

func (c *MockFailingWhisperHTTPClient) Post(url string, opts *goxios.RequestOpts) (*http.Response, error) {
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

func TestTranscription(t *testing.T) {
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
			httpClient:  new(MockWhisperHTTPClient),
			expectedErr: false,
		},
		{
			name:        "Failed request",
			httpClient:  new(MockFailingWhisperHTTPClient),
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := &TranscriptionsRequestBody{
				Model:    DefaultTranscriptionModel,
				Filename: "hello.mp3",
			}
			body.AudioFilePath = "./temp/" + body.Filename
			resp, err := Transcription(client, tc.httpClient, body)
			if tc.expectedErr && err == nil {
				t.Errorf("Expected an error, but got nil")
			}

			if !tc.expectedErr && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			if resp != nil && resp.Text != "Hello." {
				t.Fail()
			}
		})
	}
}
