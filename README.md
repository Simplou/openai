# OpenAI Go

[![GoReport](https://img.shields.io/badge/%F0%9F%93%9D%20goreport-A%2B-75C46B?style=flat-square)](https://goreportcard.com/report/github.com/Simplou/openai)

This package provides an interface for interacting with services offered by OpenAI, such as Text-to-Speech (TTS), Transcription, and Chat (Completion, Function Calling) using advanced language models.

## Installation

To install the OpenAI package, use the go get command:

```bash
go get github.com/Simplou/openai
```

## Usage

### Text-to-Speech (TTS)

```go
package main

import (
	"fmt"
	"io"
	"os"
	"context"

	"github.com/Simplou/goxios"
	"github.com/Simplou/openai"
)

const fileName = "hello.mp3"

var (
	ctx        = context.Background()
	apiKey     = os.Getenv("OPENAI_KEY")
	client     = openai.New(ctx, apiKey)
	httpClient = goxios.New(ctx)
	audioFilePath = fmt.Sprintf("./temp/%s", fileName)
)

func main() {
	audio, err := openai.TextToSpeech(client, httpClient, &openai.SpeechRequestBody{
		Model: "tts-1",
		Input: "Hello",
		Voice: openai.SpeechVoices.Onyx,
	})
	if err != nil {
		panic(err)
	}
	defer audio.Close()
	b, err := io.ReadAll(audio)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(audioFilePath, b, os.ModePerm); err != nil {
		panic(err)
	}
}
```

### Audio Transcription

```go
package main

import (
	"log"
	"os"
	"context"

	"github.com/Simplou/goxios"
    "github.com/Simplou/openai"
)

const fileName = "hello.mp3"

var (
	ctx        = context.Background()
	apiKey     = os.Getenv("OPENAI_KEY")
	client     = openai.New(ctx, apiKey)
	httpClient = goxios.New(ctx)
)

func main() {
	transcription, err := openai.Transcription(client, httpClient, &openai.TranscriptionsRequestBody{
		Model: openai.DefaultTranscriptionModel,
		Filename: fileName,
		AudioFilePath: audioFilePath,
	})
	if err != nil{
		panic(err)
	}
	println(transcription.Text)
}
```

### Chat Completion

```go
package main

import (
	"log"
	"os"
	"context"

	"github.com/Simplou/goxios"
	"github.com/Simplou/openai"
)

var (
	ctx        = context.Background()
	apiKey     = os.Getenv("OPENAI_KEY")
	client     = openai.New(ctx, apiKey)
	httpClient = goxios.New(ctx)
)

func main() {
	body := &openai.CompletionRequest[openai.DefaultMessages]{
		Model: "gpt-3.5-turbo",
		Messages: openai.DefaultMessages{
			{Role: "user", Content: "Hello"},
		},
	}

	res, err := openai.ChatCompletion(client, httpClient, body)
	if err != nil {
		panic(err)
	}
	log.Println(res.Choices[0].Message.Content)
}
```

### Function Calling

```go
package main

import (
	"encoding/json"
	"os"
	"context"

	"github.com/Simplou/goxios"
	"github.com/Simplou/openai"
)

var (
	ctx        = context.Background()
	apiKey     = os.Getenv("OPENAI_KEY")
	client     = openai.New(ctx, apiKey)
	httpClient = goxios.New(ctx)
)

func main() {
	type function func(string)
	functionRegistry := goxios.GenericJSON[function]{}
	sendEmailFnName  := "sendEmail"
	functionRegistry[sendEmailFnName] = func(email string) {
		println("email ", email)
	}
	body := &openai.CompletionRequest[openai.DefaultMessages]{
		Model: "gpt-3.5-turbo",
		Messages: openai.DefaultMessages{
			{Role: "user", Content: "send email to 93672097+gabrielluizsf@users.noreply.github.com"},
		},
		Tools: []openai.Tool{
			{
				Type: "function",
				Function: openai.Function{
					Name:        sendEmailFnName,
					Description: "send email",
					Parameters: openai.FunctionParameters{
						Type: "object",
						FunctionProperties: openai.FunctionProperties{
							"email": {
								Type:        "string",
								Description: "email provided by user",
							},
						},
					},
				},
			},
		},
		ToolChoice: "auto",
	}
	res, err := openai.ChatCompletion(client, httpClient, body)
	if err != nil {
		panic(err)
	}
	toolCalls := res.Choices[0].Message.ToolCalls
	if len(toolCalls) > 0 {
		var argumentsMap goxios.GenericJSON[string]
		if err := json.Unmarshal([]byte(toolCalls[0].Function.Args), &argumentsMap); err != nil {
			panic(err)
		}
		functionRegistry[toolCalls[0].Function.Name](argumentsMap["email"])
	}
}

```

## Contribution

If you want to contribute improvements to this package, feel free to open an issue or send a pull request.

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
