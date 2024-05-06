package main

import (
	"context"
	"log"
	"os"

	"github.com/Simplou/goxios"
	openaigo "github.com/Simplou/openai-go"
)

func chat() {
	var (
		ctx        = context.Background()
		apiKey     = os.Getenv("OPENAI_KEY")
		client     = openaigo.New(ctx, apiKey)
		httpClient = goxios.New(ctx)
	)
 	body := &openaigo.CompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openaigo.Message{
			{Role: "user", Content: "Hello"},
		},
	}
	
	res, err := openaigo.ChatCompletion(client, httpClient, body)
	if err != nil{
		panic(err)
	}
	log.Println(res.Choices[0].Message.Content)
}
