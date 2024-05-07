package openai

import (
	"bytes"
	"encoding/json"

	"github.com/Simplou/goxios"
)

// CompletionRequest represents the structure of the request sent to the OpenAI API.
type CompletionRequest struct {
	Model      string    `json:"model"`
	Messages   []Message `json:"messages"`
	ToolChoice string    `json:"tool_choice,omitempty"`
	Tools      []Tool    `json:"tools,omitempty"`
}

// Message represents a message in the conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool represents a tool that can be used during the conversation.
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function,omitempty"`
}

// Function represents a function call that can be used as a tool.
type Function struct {
	Name string `json:"name"`
}

// CompletionResponse represents the structure of the response received from the OpenAI API.
type CompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a response choice in the conversation.
type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs,omitempty"`
	FinishReason string      `json:"finish_reason"`
}

// Usage represents the token usage in the request and response.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func ChatCompletion(client OpenAIClient, httpClient HTTPClient, body *CompletionRequest) (*CompletionResponse, error) {
	client.AddHeader(goxios.Header{Key: "Content-Type", Value: "application/json"})
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	options := &goxios.RequestOpts{
		Headers: Headers(),
		Body:    bytes.NewBuffer(b),
	}
	res, err := httpClient.Post(client.BaseURL()+"/chat/completions", options)
	if err != nil {
		return nil, err
	}
	response := new(CompletionResponse)
	if err := goxios.DecodeJSON(res.Body, response); err != nil {
		return nil, err
	}

	if err := res.Body.Close(); err != nil {
		return nil, err
	}
	return response, nil
}
