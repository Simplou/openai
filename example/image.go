package main

import (
	"fmt"

	"github.com/Simplou/goxios"
	"github.com/Simplou/openai"
)

func ImageGenerated(filePath string) bool {
	return fileExists(filePath)
}

type mascot struct {
	Tool          string
	Name          string
	imageFilePath string
}

func NewMascot(tool, name, imageFilePath string) *mascot {
	return &mascot{
		Tool:          tool,
		Name:          name,
		imageFilePath: imageFilePath,
	}
}

const imagePath = "./temp/"

func image() {
	animals := goxios.GenericJSON[*mascot]{
		"gopher": NewMascot("golang", "Gopher", imagePath+"gopher.png"),
		"mouse":  NewMascot("c++", "Keith", imagePath+"keith.png"),
		"whale":  NewMascot("docker", "Moby Dock", imagePath+"moby-dock.png"),
	}
	for animal, mascot := range animals {
		prompt := fmt.Sprintf("generate %s mascot %s from tool %s creating a robot", animal, mascot.Name, mascot.Tool)
		if !ImageGenerated(mascot.imageFilePath) {
			body := &openai.ImagesGenerationsRequestBody{
				Model:  "dall-e-3",
				Prompt: prompt,
				N:      1,
				Size:   "1024x1024",
				Style:  openai.VividImageStyle().String(),
			}
			res, err := openai.ImagesGenerations(client, httpClient, body)
			if err != nil {
				panic(err)
			}
			if err := res.Download(httpClient, []string{mascot.imageFilePath}); err != nil {
				panic(err)
			}
		}
	}
}

func vision() {
	body := &openai.CompletionRequest[openai.MediaMessages]{
		Model: "gpt-4-turbo",
		Messages: openai.MediaMessages{
			{
				Role: "user",
				Content: []openai.MediaMessage{
					{Type: "text", Text: "Create a detailed prompt describing the distinct characteristics (such as color, eye features, and overall shape) of the Golang gopher depicted in the provided image. This prompt will be used to instruct OpenAI DALL-E-3 in generating a highly realistic image of the Golang gopher engaged in the process of creating a robot."},
					{Type: "image_url", ImageUrl: openai.ImageUrl("https://raw.githubusercontent.com/egonelbre/gophers/master/sketch/science/power-to-the-masses.png")},
				},
			},
		},
	}
	res, err := openai.ChatCompletion(client, httpClient, body)
	if err != nil {
		panic(err)
	}
	prompt := res.Choices[0].Message.Content
	println(prompt)
	filePath := imagePath + "gopher.png"
	if !ImageGenerated(filePath) {
		body := &openai.ImagesGenerationsRequestBody{
			Model:  "dall-e-3",
			Prompt: prompt,
			N:      1,
			Size:   "1024x1024",
			Style:  openai.VividImageStyle().String(),
		}
		res, err := openai.ImagesGenerations(client, httpClient, body)
		if err != nil {
			panic(err)
		}
		if err := res.Download(httpClient, []string{filePath}); err != nil {
			panic(err)
		}
	}
}
