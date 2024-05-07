package main

import (
	"log"

	"github.com/Simplou/openai/v1"
)

func chat() {
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
