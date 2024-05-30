package openai

import (
	"io"
	"net/http"

	"github.com/Simplou/goxios"
)

var (
	internalError = func(err error, t string) *OpenAIErr {
		return NewOpenAIErr(err, 500, t)
	}
	errCannotOpenFile = func(err error) *OpenAIErr {
		return internalError(err, "cannot_open_file")
	}
	errCannotCreateFormFile = func(err error) *OpenAIErr {
		return internalError(err, "cannot_create_form_file")
	}
	errCannotCopyFileContent = func(err error) *OpenAIErr {
		return internalError(err, "cannot_copy_file_content")
	}
	errCannotSendRequest = func(err error) *OpenAIErr {
		return internalError(err, "cannot_send_request")
	}
	errCannotDecodeJSON = func(err error) *OpenAIErr {
		return internalError(err, "cannot_decode_json")
	}
	errCannotMarshalJSON = func(err error) *OpenAIErr {
		return internalError(err, "cannot_marshal_json")
	}
	errCloseBody = func(err error) *OpenAIErr {
		return internalError(err, "close_body_error")
	}
)

type OpenAIErr struct {
	Err    JSONErr `json:"error"`
	status int
}

func (o *OpenAIErr) Error() string {
	return o.Err.Message
}

func (o *OpenAIErr) Status() int {
	return o.status
}

func NewOpenAIErr(err error, statusCode int, t string) *OpenAIErr {
	if err != nil {
		return &OpenAIErr{
			Err: JSONErr{
				Message: err.Error(),
				Type:    t,
			},
			status: statusCode,
		}
	}
	return nil
}

func closeBody(body io.ReadCloser, err error) *OpenAIErr {
	if err := body.Close(); err != nil {
		return errCloseBody(err)
	}
	if err, ok := err.(*OpenAIErr); ok {
		return err
	}
	if err != nil {
		return internalError(err, "internal_error")
	}
	return nil
}

func openaiHttpError(res *http.Response) *OpenAIErr {
	err := new(OpenAIErr)
	if err := goxios.DecodeJSON(res.Body, err); err != nil {
		return closeBody(res.Body, err)
	}
	err.status = res.StatusCode
	return closeBody(res.Body, err)
}

type JSONErr struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code,omitempty"`
}
