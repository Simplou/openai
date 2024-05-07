package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

// openaiSpeechVoices holds the available voices for OpenAI speech synthesis.
type openaiSpeechVoices struct {
	Alloy   string
	Echo    string
	Fable   string
	Onyx    string
	Nova    string
	Shimmer string
}

// SpeechRequestBody represents the request body for the speech API.
type SpeechRequestBody struct {
	Model string `json:"model"` // The model for speech synthesis.
	Input string `json:"input"` // The input text for synthesis.
	Voice string `json:"voice"` // The voice to be used for synthesis.
}

func closeBody(body io.ReadCloser, err error) error {
	if err := body.Close(); err != nil {
		return err
	}
	return err
}

func TextToSpeech(api OpenAIClient, httpClient HTTPClient, body *SpeechRequestBody) (io.ReadCloser, error) {
	api.AddHeader(contentTypeJSON)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	options := goxios.RequestOpts{
		Headers: Headers(),
		Body:    bytes.NewBuffer(b),
	}
	res, err := httpClient.Post(api.BaseURL()+"/audio/speech", &options)
	if res.StatusCode != http.StatusOK {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, closeBody(res.Body, err)
		}
		errMessage := fmt.Sprintf("%s\n%s", res.Status, string(b))
		return nil, closeBody(res.Body, errors.New(errMessage))
	}
	if err != nil {
		return nil, closeBody(res.Body, err)
	}
	return res.Body, nil
}
