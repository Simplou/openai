package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Simplou/openai"
)

const fileName = "hello.mp3"

var (
	audioFilePath = fmt.Sprintf("./temp/%s", fileName)
)

func AudioGenerated(filePath string) bool {
	return fileExists(filePath)
}

func tts() {
	fileExists := AudioGenerated(audioFilePath)
	if !fileExists {
		audio, openaiErr := openai.TextToSpeech(client, httpClient, &openai.SpeechRequestBody{
			Model: "tts-1",
			Input: "Hello",
			Voice: openai.SpeechVoices.Onyx,
		})
		if openaiErr != nil {
			b, err := json.Marshal(openaiErr)
			if err != nil {
				panic(err)
			}
			log.Fatal(string(b))
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
}

func whisper() {
	transcription, err := openai.Transcription(client, httpClient, &openai.TranscriptionsRequestBody{
		Model:         openai.DefaultTranscriptionModel,
		Filename:      fileName,
		AudioFilePath: audioFilePath,
	})
	if err != nil {
		panic(err)
	}
	println(transcription.Text)
}
