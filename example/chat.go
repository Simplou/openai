package main

import (
	"context"
	"log"
	"os"

	"github.com/Simplou/goxios"
	openai "github.com/Simplou/openai"
)

func chat() {
	var (
		ctx        = context.Background()
		apiKey     = os.Getenv("OPENAI_KEY")
		client     = openai.New(ctx, apiKey)
		httpClient = goxios.New(ctx)
	)
 	body := &openai.CompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.Message{
			{Role: "user", Content: "Hello"},
		},
	}
	
	res, err := openai.ChatCompletion(client, httpClient, body)
	if err != nil{
		panic(err)
	}
	log.Println(res.Choices[0].Message.Content)
}
