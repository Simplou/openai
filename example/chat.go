package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Simplou/goxios"
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
	type function func(string)
	functionRegistry := goxios.GenericJSON[function]{}
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
		var argumentsMap goxios.GenericJSON[string]
		if err := json.Unmarshal([]byte(toolCalls[0].Function.Args), &argumentsMap); err != nil {
			panic(err)
		}
		functionRegistry[toolCalls[0].Function.Name](argumentsMap["email"])
	}
}

func chatByEmbedding(largeText string, query string) {
	createEmbedding := func(chunks []string) (*openai.EmbeddingResponse[[]float64], error) {
		emb, err := openai.CreateEmbedding[[]string, []float64](client, httpClient, &openai.EmbeddingRequest[[]string]{
			Model: "text-embedding-ada-002",
			Input: chunks,
		})
		if err != nil {
			return &openai.EmbeddingResponse[[]float64]{}, nil
		}
		return emb, nil
	}
	queryChunks := openai.ChunkText(openai.ChunkTextOpts{Text: query})
	queryEmb, err := createEmbedding(queryChunks)
	if err != nil {
		log.Println(err)
	}
	q := queryEmb.Data[0].Embedding
	qMatrix := make([][]float64, 0)
	qMatrix = append(qMatrix, q)
	chunks := openai.ChunkText(openai.ChunkTextOpts{Text: largeText})
	emb, err := createEmbedding(chunks)
	if err != nil {
		log.Println(err)
	}
	e := emb.Data[0].Embedding
	eMatrix := make([][]float64, 0)
	eMatrix = append(eMatrix, e)
	indexes, err := openai.FindMostRelevantEmbeddings(qMatrix, eMatrix)
	if err != nil {
		log.Println(err)
	}
	var relevantChunks []string
	for _, i := range indexes {
		relevantChunks = append(relevantChunks, chunks[i])
	}
	summary, err := openai.ChunksSummary(client, httpClient, relevantChunks, query)
	if err != nil {
		log.Println(err)
	}
	log.Println(summary)
}

func chatModerator(customerMessage string) {
	moderation, err := openai.Moderator(client, httpClient, &openai.ModerationRequest[string]{
		Input: customerMessage,
	})
	if err != nil {
		log.Println(err)
	}
	categories := make([]string, 0)
	for _, v := range moderation.Results {
		if v.Flagged {
			for category, value := range v.Categories {
				if value {
					categories = append(categories, category)
				}
			}
		}
	}
	if len(categories) == 0 {
		res, err := openai.ChatCompletion[openai.DefaultMessages](
			client,
			httpClient,
			&openai.CompletionRequest[openai.DefaultMessages]{
				Model: "gpt-3.5-turbo",
				Messages: openai.DefaultMessages{
					{Role: "user", Content: customerMessage},
				},
			},
		)
		if err != nil {
			log.Println(err)
		}
		log.Println(res.Choices[0].Message.Content)
	} else {
		s := strings.Join(categories, ", ")
		moderatorMessage := fmt.Sprintf("Your statement contains several disrespectful things: (%s)", s)
		log.Println(moderatorMessage)
	}

}
