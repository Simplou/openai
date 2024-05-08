package main

import (
	"fmt"
	"os"

	"github.com/Simplou/goxios"
	"github.com/Simplou/openai"
)

func ImageGenerated(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
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
		prompt := fmt.Sprintf("generate %s mascot %s from tool %s by reading a book", animal, mascot.Name, mascot.Tool)
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
