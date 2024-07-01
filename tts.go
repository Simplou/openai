package openai

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Simplou/goxios"
)

const (
	alloy   = "alloy"
	echo    = "echo"
	fable   = "fable"
	onyx    = "onyx"
	nova    = "nova"
	shimmer = "shimmer"
)

var SpeechVoices = &openaiSpeechVoices{
	Alloy:   alloy,
	Echo:    echo,
	Fable:   fable,
	Onyx:    onyx,
	Nova:    nova,
	Shimmer: shimmer,
}

type (
	// openaiSpeechVoices holds the available voices for OpenAI speech synthesis.
	openaiSpeechVoices struct {
		Alloy   string
		Echo    string
		Fable   string
		Onyx    string
		Nova    string
		Shimmer string
	}

	// SpeechRequestBody represents the request body for the speech API.
	SpeechRequestBody struct {
		Model string `json:"model"` // The model for speech synthesis.
		Input string `json:"input"` // The input text for synthesis.
		Voice string `json:"voice"` // The voice to be used for synthesis.
	}
)

func TextToSpeech(api OpenAIClient, httpClient HTTPClient, body *SpeechRequestBody) (io.ReadCloser, *OpenAIErr) {
	api.AddHeader(contentTypeJSON)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, NewOpenAIErr(err, 500, "marshal_json_error")
	}
	options := goxios.RequestOpts{
		Headers: api.Headers(),
		Body:    ioReader(b),
	}
	res, err := httpClient.Post(api.BaseURL()+"/audio/speech", &options)
	if res.StatusCode != http.StatusOK {
		return nil, openaiHttpError(res)
	}
	if err != nil {
		return nil, closeBody(res.Body, err)
	}
	return res.Body, nil
}
