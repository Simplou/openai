package main

import (
	"context"
	"os"

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
	chat()
	functionCall()
	tts()
	whisper()
	image()
}
