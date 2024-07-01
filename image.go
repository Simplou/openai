package openai

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/Simplou/goxios"
)

type imageStyle string

func (is imageStyle) String() string {
	return string(is)
}

// NaturalImageStyle causes the model to produce more natural, less hyper-real looking images.
func NaturalImageStyle() imageStyle {
	return imageStyle("natural")
}

// VividImageStyle causes the model to lean towards generating hyper-real and dramatic images
func VividImageStyle() imageStyle {
	return imageStyle("vivid")
}

type (
	ImagesGenerationsRequestBody struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		N      int    `json:"n"`
		Size   string `json:"size"`
		Style  string `json:"style,omitempty"` //This param is only supported for dall-e-3.
	}

	ImagesGenerationsResponse struct {
		Created int64 `json:"created"`
		Data    []struct {
			Url string `json:"url"`
		} `json:"data"`
	}
)

func (igr *ImagesGenerationsResponse) Download(httpClient HTTPClient, filePaths []string) error {
	if len(filePaths) != len(igr.Data) {
		return errors.New("number of file paths does not match number of images")
	}
	for i, image := range igr.Data {
		res, err := httpClient.Get(image.Url, &goxios.RequestOpts{})
		if err != nil {
			return err
		}
		defer res.Body.Close()

		file, err := os.Create(filePaths[i])
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, res.Body)
		if err != nil {
			return err
		}
	}
	return nil
}

func ImagesGenerations(api OpenAIClient, httpClient HTTPClient, body *ImagesGenerationsRequestBody) (*ImagesGenerationsResponse, *OpenAIErr) {
	api.AddHeader(contentTypeJSON)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, errCannotMarshalJSON(err)
	}
	res, err := httpClient.Post(api.BaseURL()+"/images/generations", &goxios.RequestOpts{
		Body:    ioReader(b),
		Headers: api.Headers(),
	})
	if err != nil {
		return nil, errCannotSendRequest(err)
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, openaiHttpError(res)
	}
	images := new(ImagesGenerationsResponse)
	if err := goxios.DecodeJSON(res.Body, images); err != nil {
		return nil, errCannotDecodeJSON(err)
	}
	return images, nil
}
