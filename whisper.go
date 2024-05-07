package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/Simplou/goxios"
)

type TranscriptionsRequestBody struct {
	Model, Filename, AudioFilePath string
}

const DefaultTranscriptionModel = "whisper-1"

type TranscriptionResponse struct {
	Text string `json:"text"`
}

func Transcription(api OpenAIClient, httpClient HTTPClient, body *TranscriptionsRequestBody) (*TranscriptionResponse, error) {
	file, err := os.Open(body.AudioFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)

	part, err := writer.CreateFormFile("file", body.Filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	writer.WriteField("model", body.Model)

	writer.Close()

	api.AddHeader(goxios.Header{Key: "Content-Type", Value: writer.FormDataContentType()})
	requestOptions := goxios.RequestOpts{
		Headers: Headers(),
		Body:    b,
	}
	res, err := httpClient.Post(api.BaseURL()+"/audio/transcriptions", &requestOptions)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := new(TranscriptionResponse)

	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, closeBody(res.Body, err)
		}
		errMessage := fmt.Sprintf("%s %s", res.Status, string(b))
		return nil, closeBody(res.Body, errors.New(errMessage))
	}
	return result, nil
}
