package main

import (
	"encoding/json"
	"log"

	"github.com/Simplou/openai"
)

func chat() {
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

func functionCall() {
	functionRegistry := map[string]func(string){}
	sendEmailFnName := "sendEmail"
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
		var argumentsMap map[string]string
		if err := json.Unmarshal([]byte(toolCalls[0].Function.Args), &argumentsMap); err != nil {
			panic(err)
		}
		functionRegistry[toolCalls[0].Function.Name](argumentsMap["email"])
	}
}
