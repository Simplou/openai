# OpenAI Go

[![GoReport](https://img.shields.io/badge/%F0%9F%93%9D%20goreport-A%2B-75C46B?style=flat-square)](https://goreportcard.com/report/github.com/Simplou/openai)

This package provides an interface for interacting with services offered by OpenAI, such as Text-to-Speech (TTS), Transcription, and Chat using advanced language models.

## Installation

To install the OpenAI package, use the go get command:

```bash
go get github.com/Simplou/openai/v1
```

## Usage

### Text-to-Speech (TTS)

```go
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Simplou/openai/v1"
)

const fileName = "hello.mp3"

var (
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

	openai "github.com/Simplou/openai/v1"
)

const fileName = "hello.mp3"

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

	openai "github.com/Simplou/openai/v1"
)

func main() {
	body := &openai.CompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.Message{
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

## Contribution

If you want to contribute improvements to this package, feel free to open an issue or send a pull request.

## License

This package is licensed under the MIT License.  See the [LICENSE](LICENSE) file for details.
