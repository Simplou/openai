package openai

import (
	"encoding/json"
	"net/http"

	"github.com/Simplou/goxios"
)

type DefaultMessages []Message[string]
type MediaMessages []Message[[]MediaMessage]

// CompletionRequest represents the structure of the request sent to the OpenAI API.
type CompletionRequest[T any] struct {
	Model      string `json:"model"`
	Messages   T      `json:"messages"`
	ToolChoice string `json:"tool_choice,omitempty"`
	Tools      []Tool `json:"tools,omitempty"`
}

type MediaMessage struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageUrl *imageUrl `json:"image_url,omitempty"`
}

type imageUrl struct {
	Url string `json:"url"`
}

func ImageUrl(url string) *imageUrl {
	return &imageUrl{url}
}

// Message represents a message in the conversation.
type Message[T string | []MediaMessage] struct {
	Role      string     `json:"role"`
	Content   T          `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name string `json:"name"`
		Args string `json:"arguments"`
	} `json:"function"`
}

// Tool represents a tool that can be used during the conversation.
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function,omitempty"`
}

// Function represents a function call that can be used as a tool.
type Function struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  FunctionParameters `json:"parameters"`
}

type FunctionProperties goxios.GenericJSON[FunctionPropertie]

// FunctionParameters represents the parameters of a function
type FunctionParameters struct {
	Type               string `json:"type"`
	FunctionProperties `json:"properties"`
}

// FunctionPropertie represents a property of a function.
type FunctionPropertie struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
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
	Index        int             `json:"index"`
	Message      Message[string] `json:"message"`
	Logprobs     interface{}     `json:"logprobs,omitempty"`
	FinishReason string          `json:"finish_reason"`
}

// Usage represents the token usage in the request and response.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens"`
}

func ChatCompletion[Messages any](api OpenAIClient, httpClient HTTPClient, body *CompletionRequest[Messages]) (*CompletionResponse, *OpenAIErr) {
	api.AddHeader(contentTypeJSON)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, errCannotMarshalJSON(err)
	}
	options := &goxios.RequestOpts{
		Headers: api.Headers(),
		Body:    ioReader(b),
	}
	res, err := httpClient.Post(api.BaseURL()+"/chat/completions", options)
	if err != nil {
		return nil, errCannotSendRequest(err)
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, openaiHttpError(res)
	}
	response := new(CompletionResponse)
	if err := goxios.DecodeJSON(res.Body, response); err != nil {
		return nil, errCannotDecodeJSON(err)
	}

	if err := res.Body.Close(); err != nil {
		return nil, errCloseBody(err)
	}
	return response, nil
}
