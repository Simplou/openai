package openai

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"

	"github.com/Simplou/goxios"
)

const DefaultTranscriptionModel = "whisper-1"

type (
	TranscriptionsRequestBody struct {
		Model, Filename, AudioFilePath string
	}

	TranscriptionResponse struct {
		Text string `json:"text"`
	}
)

func Transcription(api OpenAIClient, httpClient HTTPClient, body *TranscriptionsRequestBody) (*TranscriptionResponse, *OpenAIErr) {
	file, err := os.Open(body.AudioFilePath)
	if err != nil {
		return nil, errCannotOpenFile(err)
	}
	defer file.Close()

	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)

	part, err := writer.CreateFormFile("file", body.Filename)
	if err != nil {
		return nil, errCannotCreateFormFile(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, errCannotCopyFileContent(err)
	}

	writer.WriteField("model", body.Model)

	writer.Close()

	api.AddHeader(goxios.Header{Key: "Content-Type", Value: writer.FormDataContentType()})
	requestOptions := goxios.RequestOpts{
		Headers: api.Headers(),
		Body:    b,
	}
	res, err := httpClient.Post(api.BaseURL()+"/audio/transcriptions", &requestOptions)
	if err != nil {
		return nil, errCannotSendRequest(err)
	}
	defer res.Body.Close()

	result := new(TranscriptionResponse)

	if err := goxios.DecodeJSON(res.Body, result); err != nil {
		return nil, errCannotDecodeJSON(err)
	}
	if res.StatusCode != 200 {
		return nil, openaiHttpError(res)
	}
	return result, nil
}
