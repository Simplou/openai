package openai

import (
	"encoding/json"

	"github.com/Simplou/goxios"
)

type (
	ModerationRequest[Input string | []string] struct {
		Input Input  `json:"input"`
		Model string `json:"model,omitempty"`
	}

	ModerationResponse struct {
		Id      string `json:"id"`
		Model   string `json:"model"`
		Results []struct {
			Flagged    bool                     `json:"flagged"`
			Categories goxios.GenericJSON[bool] `json:"categories"`
		} `json:"results"`
		CategoryScores goxios.GenericJSON[float64] `json:"category_scores"`
	}
)

func Moderator[Input string | []string](api OpenAIClient, httpClient HTTPClient, body *ModerationRequest[Input]) (*ModerationResponse, *OpenAIErr) {
	api.AddHeader(contentTypeJSON)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, errCannotMarshalJSON(err)
	}
	options := goxios.RequestOpts{
		Body:    ioReader(b),
		Headers: api.Headers(),
	}
	res, err := httpClient.Post(api.BaseURL()+"/moderations", &options)
	if err != nil {
		return nil, errCannotSendRequest(err)
	}
	response := new(ModerationResponse)
	if err := goxios.DecodeJSON(res.Body, response); err != nil {
		return nil, errCannotDecodeJSON(err)
	}
	if err := res.Body.Close(); err != nil {
		return nil, errCloseBody(err)
	}
	return response, nil
}
